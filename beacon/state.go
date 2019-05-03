package beacon

import (
	"errors"
	"fmt"
	"time"

	"github.com/phoreproject/prysm/shared/ssz"

	"github.com/phoreproject/synapse/primitives"
	logger "github.com/sirupsen/logrus"
)

// StoreBlock adds a block header to the current chain. The block should already
// have been validated by this point.
func (b *Blockchain) StoreBlock(block *primitives.Block) error {
	err := b.db.SetBlock(*block)
	if err != nil {
		return err
	}

	return nil
}

// AddBlockToStateMap calculates the state after applying block and adds it
// to the state map.
func (b *Blockchain) AddBlockToStateMap(block *primitives.Block) (*primitives.State, error) {
	return b.stateManager.AddBlockToStateMap(block)
}

// ProcessBlock is called when a block is received from a peer.
func (b *Blockchain) ProcessBlock(block *primitives.Block, checkTime bool) error {
	genesisTime := b.stateManager.GetGenesisTime()

	validationStart := time.Now()

	// VALIDATE BLOCK HERE
	if checkTime && (block.BlockHeader.SlotNumber*uint64(b.config.SlotDuration)+genesisTime > uint64(time.Now().Unix()) || block.BlockHeader.SlotNumber == 0) {
		return errors.New("block slot too soon")
	}

	seen := b.View.Index.Has(block.BlockHeader.ParentRoot)
	if !seen {
		return errors.New("do not have parent block")
	}

	blockHash, err := ssz.TreeHash(block)
	if err != nil {
		return err
	}

	seen = b.View.Index.Has(blockHash)

	if seen {
		// we've already processed this block
		return nil
	}

	validationTime := time.Since(validationStart)

	blockHashStr := fmt.Sprintf("%x", blockHash)

	logger.WithField("hash", blockHashStr).Info("processing new block")

	stateCalculationStart := time.Now()

	initialState, found := b.stateManager.GetStateForHash(block.BlockHeader.ParentRoot)
	if !found {
		return errors.New("could not find state for parent block")
	}

	// Check state root.
	if true {
		h, err := b.GetHashBySlot(block.BlockHeader.SlotNumber - 1)
		if err != nil {
			return err
		}
		lastBlock, err := b.GetBlockByHash(h)
		if err != nil {
			return err
		}
		h, err = ssz.TreeHash(lastBlock)
		if err != nil {
			return err
		}
		/*
			expectedState, found := b.stateManager.GetStateForHash(h)
			if !found {
				return errors.New("could not find state for block")
			}

			expectedStateRoot, err := ssz.TreeHash(expectedState)
			if err != nil {
				return err
			}
		*/
		if block.BlockHeader.StateRoot != h {
			return errors.New("StateRoot doesn't match")
		}
	}

	initialJustifiedSlot := initialState.JustifiedSlot
	initialFinalizedSlot := initialState.FinalizedSlot

	newState, err := b.AddBlockToStateMap(block)
	if err != nil {
		return err
	}

	stateCalculationTime := time.Since(stateCalculationStart)

	blockStorageStart := time.Now()

	err = b.StoreBlock(block)
	if err != nil {
		return err
	}

	blockStorageTime := time.Since(blockStorageStart)

	logger.Debug("applied with new state")

	node, err := b.View.Index.AddBlockNodeToIndex(block, blockHash)
	if err != nil {
		return err
	}

	// TODO: these two database operations should be in a single transaction

	databaseTipUpdateStart := time.Now()

	// set the block node in the database
	err = b.db.SetBlockNode(blockNodeToDisk(*node))
	if err != nil {
		return err
	}

	// update the parent node in the database
	err = b.db.SetBlockNode(blockNodeToDisk(*node.Parent))
	if err != nil {
		return err
	}

	databaseTipUpdateTime := time.Since(databaseTipUpdateStart)

	attestationUpdateStart := time.Now()

	for _, a := range block.BlockBody.Attestations {
		participants, err := newState.GetAttestationParticipants(a.Data, a.ParticipationBitfield, b.config)
		if err != nil {
			return err
		}

		err = b.db.SetLatestAttestationsIfNeeded(participants, a)
		if err != nil {
			return err
		}
	}

	attestationUpdateEnd := time.Since(attestationUpdateStart)

	logger.Debug("updating chain head")

	updateChainHeadStart := time.Now()

	err = b.UpdateChainHead()
	if err != nil {
		return err
	}

	updateChainHeadTime := time.Since(updateChainHeadStart)

	connectBlockSignalStart := time.Now()

	b.ConnectBlockNotifier.Signal(*block)

	connectBlockSignalTime := time.Since(connectBlockSignalStart)

	finalizedStateUpdateStart := time.Now()

	finalizedNode := node.GetAncestorAtSlot(newState.FinalizedSlot)
	if finalizedNode == nil {
		return errors.New("could not find finalized node in block index")
	}
	finalizedState, found := b.stateManager.GetStateForHash(finalizedNode.Hash)
	if !found {
		return errors.New("could not find finalized block Hash in state map")
	}

	if initialFinalizedSlot != newState.FinalizedSlot {
		err := b.db.SetFinalizedState(*finalizedState)
		if err != nil {
			return err
		}
	}

	finalizedNodeAndState := blockNodeAndState{finalizedNode, *finalizedState}
	b.View.finalizedHead = finalizedNodeAndState

	err = b.db.SetFinalizedHead(finalizedNode.Hash)
	if err != nil {
		return err
	}

	justifiedNode := node.GetAncestorAtSlot(newState.JustifiedSlot)
	if justifiedNode == nil {
		return errors.New("could not find justified node in block index")
	}

	justifiedState, found := b.stateManager.GetStateForHash(justifiedNode.Hash)
	if !found {
		return errors.New("could not find justified block Hash in state map")
	}
	justifiedNodeAndState := blockNodeAndState{justifiedNode, *justifiedState}
	b.View.justifiedHead = justifiedNodeAndState

	if initialJustifiedSlot != newState.JustifiedSlot {
		err := b.db.SetJustifiedState(*justifiedState)
		if err != nil {
			return err
		}
	}

	err = b.db.SetJustifiedHead(justifiedNode.Hash)
	if err != nil {
		return err
	}

	finalizedStateUpdateTime := time.Since(finalizedStateUpdateStart)

	stateCleanupStart := time.Now()

	err = b.stateManager.DeleteStateBeforeFinalizedSlot(finalizedNode.Slot)
	if err != nil {
		return err
	}

	stateCleanupTime := time.Since(stateCleanupStart)

	logger.WithFields(logger.Fields{
		"validation":         validationTime,
		"stateCalculation":   stateCalculationTime,
		"storage":            blockStorageTime,
		"databaseTipUpdate":  databaseTipUpdateTime,
		"attestationUpdate":  attestationUpdateEnd,
		"updateChainHead":    updateChainHeadTime,
		"connectBlockSignal": connectBlockSignalTime,
		"finalizedUpdate":    finalizedStateUpdateTime,
		"stateCleanup":       stateCleanupTime,
		"totalTime":          time.Since(validationStart),
	}).Debug("performance report for processing")

	return nil
}

// GetState gets a copy of the current state of the blockchain.
func (b *Blockchain) GetState() primitives.State {
	return b.stateManager.GetHeadState()
}
