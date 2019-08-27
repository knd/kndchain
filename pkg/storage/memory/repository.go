package memory

import (
	"errors"
	"log"

	"github.com/knd/kndchain/pkg/listing"
	"github.com/knd/kndchain/pkg/mining"
)

// ErrAddNilBlock is used when no mined block is given to add
var ErrAddNilBlock = errors.New("Mined block is not given to add")

// MemStorage keeps blockchain in memory
type MemStorage struct {
	blockchain *Blockchain
}

// NewRepository creates a blockchain repository
func NewRepository() *MemStorage {
	return &MemStorage{
		blockchain: &Blockchain{},
	}
}

// AddBlock adds mined block into blockchain
func (m *MemStorage) AddBlock(minedBlock *mining.Block) error {
	if minedBlock == nil {
		return ErrAddNilBlock
	}

	newB := Block{
		Timestamp:  minedBlock.Timestamp,
		LastHash:   minedBlock.LastHash,
		Hash:       minedBlock.Hash,
		Data:       minedBlock.Data,
		Nonce:      minedBlock.Nonce,
		Difficulty: minedBlock.Difficulty,
	}

	m.blockchain.chain = append(m.blockchain.chain, newB)

	return nil
}

// GetBlockCount returns the latest block count in blockchain
func (m *MemStorage) GetBlockCount() uint32 {
	return uint32(len(m.blockchain.chain))
}

// GetLastBlock returns the last block in blockchain
func (m *MemStorage) GetLastBlock() listing.Block {
	if m.GetBlockCount() == 0 {
		log.Fatal("Blockchain is empty")
	}

	lastBlock := m.blockchain.chain[m.GetBlockCount()-1]

	return listing.Block{
		Timestamp:  lastBlock.Timestamp,
		LastHash:   lastBlock.LastHash,
		Hash:       lastBlock.Hash,
		Data:       lastBlock.Data,
		Nonce:      lastBlock.Nonce,
		Difficulty: lastBlock.Difficulty,
	}
}

// GetBlockchain returns a list of blocks from genesis block
func (m *MemStorage) GetBlockchain() *listing.Blockchain {
	var res []listing.Block
	for _, block := range m.blockchain.chain {
		res = append(res, listing.Block{
			Timestamp:  block.Timestamp,
			LastHash:   block.LastHash,
			Hash:       block.Hash,
			Data:       block.Data,
			Nonce:      block.Nonce,
			Difficulty: block.Difficulty,
		})
	}
	return &listing.Blockchain{Chain: res}
}

// ReplaceChain replace the current blockchain with the newchain
func (m *MemStorage) ReplaceChain(newChain *mining.Blockchain) error {
	if newChain == nil || len(newChain.Chain) < 1 {
		log.Fatal("Blockchain is nil/ zero block")
	}

	newBc := []Block{}
	for _, newBlock := range newChain.Chain {
		newBc = append(newBc, Block{
			Timestamp:  newBlock.Timestamp,
			LastHash:   newBlock.LastHash,
			Hash:       newBlock.Hash,
			Data:       newBlock.Data,
			Nonce:      newBlock.Nonce,
			Difficulty: newBlock.Difficulty,
		})
	}

	m.blockchain.chain = newBc

	return nil
}
