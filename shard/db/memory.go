package db

import (
	"fmt"
	"github.com/phoreproject/synapse/chainhash"
	"github.com/phoreproject/synapse/primitives"
	"github.com/prysmaticlabs/go-ssz"
)

// MemoryBlockDB is a block database stored in memory.
type MemoryBlockDB struct {
	blocks map[chainhash.Hash]primitives.ShardBlock
}

// NewMemoryBlockDB creates a new block database stored in memory.
func NewMemoryBlockDB() *MemoryBlockDB {
	return &MemoryBlockDB{
		blocks: make(map[chainhash.Hash]primitives.ShardBlock),
	}
}

// GetBlockForHash gets a block from the database.
func (m *MemoryBlockDB) GetBlockForHash(h chainhash.Hash) (*primitives.ShardBlock, error) {
	if b, found := m.blocks[h]; found {
		return &b, nil
	} else {
		return nil, fmt.Errorf("couldn't find block in database with hash %s", h)
	}
}

// SetBlock sets a block in the block database.
func (m *MemoryBlockDB) SetBlock(b *primitives.ShardBlock) error {
	blockHash, err := ssz.HashTreeRoot(b)
	if err != nil {
		return err
	}

	m.blocks[blockHash] = *b
	return nil
}

// Close closes the database, which does nothing for an in-memory database.
func (m *MemoryBlockDB) Close() error {
	return nil
}

var _ ShardBlockDatabase = &MemoryBlockDB{}