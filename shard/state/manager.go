package state

import (
	"fmt"
	"github.com/phoreproject/synapse/chainhash"
	"github.com/phoreproject/synapse/csmt"
	"github.com/phoreproject/synapse/primitives"
	"github.com/prysmaticlabs/go-ssz"
	"github.com/sirupsen/logrus"
	"sync"
)

type stateInfo struct {
	db     csmt.TreeDatabase
	parent *stateInfo
	slot   uint64
}

// ShardStateManager keeps track of state by reverting and applying blocks.
type ShardStateManager struct {
	stateLock *sync.RWMutex
	// each of these TreeDatabase's is a TreeMemoryCache except for the finalized block which is the regular state DB.
	// When finalizing a block, we commit the hash of the finalized block and recurse through each parent until we reach
	// the previously justified block (which should be a real DB, not a cache).
	finalizedDB *stateInfo
	stateMap    map[chainhash.Hash]*stateInfo
	tipDB       *stateInfo

	shardInfo ShardInfo
}

// NewShardStateManager constructs a new shard state manager which keeps track of shard state.
func NewShardStateManager(stateDB csmt.TreeDatabase, stateSlot uint64, tipBlockHash chainhash.Hash, shardInfo ShardInfo) *ShardStateManager {
	tip := &stateInfo{
		db:   stateDB,
		slot: stateSlot,
	}

	return &ShardStateManager{
		finalizedDB:  tip,
		tipDB:        tip,
		shardInfo: shardInfo,
		stateMap: map[chainhash.Hash]*stateInfo{
			tipBlockHash: tip,
		},
		stateLock: new(sync.RWMutex),
	}
}

// GetTip gets the state of the current tip.
func (sm *ShardStateManager) GetTip() csmt.TreeDatabase {
	sm.stateLock.Lock()
	defer sm.stateLock.Unlock()
	return sm.tipDB.db
}

// GetTipSlot gets the slot of the tip of the blockchain.
func (sm *ShardStateManager) GetTipSlot() uint64 {
	sm.stateLock.Lock()
	defer sm.stateLock.Unlock()

	return sm.tipDB.slot
}

func (sm *ShardStateManager) has(c chainhash.Hash) bool {
	_, found := sm.stateMap[c]
	return found
}

// Has checks if the ShardStateManager has a certain state for a block.
func (sm *ShardStateManager) Has(c chainhash.Hash) bool {
	sm.stateLock.RLock()
	defer sm.stateLock.RUnlock()
	return sm.has(c)
}

func executeBlockTransactions(a csmt.TreeTransactionAccess, b primitives.ShardBlock, si ShardInfo) error {
	for _, tx := range b.Body.Transactions {
		_, err := Transition(a, tx.TransactionData, si)
		if err != nil {
			return err
		}
	}
	return nil
}

// Add adds a block to the state map.
func (sm *ShardStateManager) Add(block *primitives.ShardBlock) (*chainhash.Hash, error) {
	sm.stateLock.Lock()
	defer sm.stateLock.Unlock()

	blockHash, _ := ssz.HashTreeRoot(block)
	if sm.has(blockHash) {
		// we already got this block
		return nil, nil
	}

	logrus.WithField("block hash", chainhash.Hash(blockHash)).Debug("add block action")

	previousTree, found := sm.stateMap[block.Header.PreviousBlockHash]
	if !found {
		return nil, fmt.Errorf("could not find parent block with hash: %s", block.Header.PreviousBlockHash)
	}

	newCache, err := csmt.NewTreeMemoryCache(previousTree.db)
	if err != nil {
		return nil, err
	}

	newCacheTree := csmt.NewTree(newCache)
	err = newCacheTree.Update(func(access csmt.TreeTransactionAccess) error {
		return executeBlockTransactions(access, *block, sm.shardInfo)
	})
	if err != nil {
		return nil, err
	}

	sm.stateMap[blockHash] = &stateInfo{
		db:     newCache,
		parent: previousTree,
		slot:   block.Header.Slot,
	}

	return newCache.Hash()
}

// Finalize removes unnecessary state from the state map.
func (sm *ShardStateManager) Finalize(finalizedHash chainhash.Hash, finalizedSlot uint64) error {
	sm.stateLock.Lock()
	defer sm.stateLock.Unlock()


	logrus.WithField("block hash", finalizedHash).Debug("finalize block action")

	// first, let's start at the tip and commit everything
	finalizeNode, found := sm.stateMap[finalizedHash]
	if !found {
		return fmt.Errorf("could not find block to finalize %s on slot %d %d", finalizedHash, finalizedSlot, sm.shardInfo.ShardID)
	}
	memCache, isCache := finalizeNode.db.(*csmt.TreeMemoryCache)
	for isCache {
		fmt.Println("flushing")
		if err := memCache.Flush(); err != nil {
			return err
		}
		prevBlock := finalizeNode.parent
		// if this is null, that means something went wrong because there is no underlying store.
		memCache, isCache = prevBlock.db.(*csmt.TreeMemoryCache)

		finalizeNode = prevBlock
	}

	// after we've flushed all of that, we should clean up any states with slots before that.
	for k, node := range sm.stateMap {
		if node.slot > finalizedSlot {
			continue
		}
		if k.IsEqual(&finalizedHash) {
			continue
		}
		if node.slot < finalizedSlot {
			delete(sm.stateMap, k)
		}
		if node.slot == finalizedSlot && !k.IsEqual(&finalizedHash) {
			delete(sm.stateMap, k)
		}
	}

	sm.finalizedDB = sm.stateMap[finalizedHash]

	return nil
}

// SetTip sets the tip of the state.
func (sm *ShardStateManager) SetTip(newTipHash chainhash.Hash) error {
	sm.stateLock.Lock()
	defer sm.stateLock.Unlock()

	logrus.WithField("block hash", newTipHash).Debug("new tip block action")

	newTip, found := sm.stateMap[newTipHash]
	if !found {
		return fmt.Errorf("couldn't find state root for tip: %s", newTipHash)
	}

	sm.tipDB = newTip

	return nil
}
