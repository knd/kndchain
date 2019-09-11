package leveldb

import (
	"github.com/knd/kndchain/pkg/listing"
	"github.com/knd/kndchain/pkg/mining"
)

// LevelDB keeps blockchain in local key-value db
type LevelDB struct {
}

// NewRepository creates a repository to interact with LevelDB
func NewRepository() *LevelDB {
	return &LevelDB{}
}

// AddBlock adds mined block into blockchain
func (db *LevelDB) AddBlock(minedBlock *mining.Block) error {
	return nil
}

// GetBlockCount returns the latest block count in blockchain
func (db *LevelDB) GetBlockCount() uint32 {
	return 0
}

// GetLastBlock returns the last block in blockchain
func (db *LevelDB) GetLastBlock() listing.Block {
	return listing.Block{}
}

// GetBlockchain returns a list of blocks from genesis block
func (db *LevelDB) GetBlockchain() *listing.Blockchain {
	return nil
}

// ReplaceChain replace the current blockchain with the newchain
func (db *LevelDB) ReplaceChain(newChain *mining.Blockchain) error {
	return nil
}
