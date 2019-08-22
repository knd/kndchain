package memory

import (
	"github.com/knd/kndchain/pkg/mining"
)

// MemStorage keeps blockchain in memory
type MemStorage struct {
	blockchain *Blockchain
}

// AddBlock adds mined block into blockchain
func (m *MemStorage) AddBlock(minedBlock *mining.Block) error {
	// TODO: validations of minedBlock

	newB := Block{
		Timestamp: minedBlock.Timestamp,
		LastHash:  minedBlock.LastHash,
		Hash:      minedBlock.Hash,
		Data:      minedBlock.Data,
	}

	m.blockchain.chain = append(m.blockchain.chain, newB)

	return nil
}
