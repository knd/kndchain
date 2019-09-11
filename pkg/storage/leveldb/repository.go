package leveldb

import (
	"log"
	"os"
	"path"

	"github.com/knd/kndchain/pkg/listing"
	"github.com/knd/kndchain/pkg/mining"
)

// LevelDB keeps blockchain in local key-value db
type LevelDB struct {
	PathToBlockData string
	PathToChainData string
}

// NewRepository creates a repository to interact with LevelDB
func NewRepository(pathToDataDir string) *LevelDB {
	r := &LevelDB{
		PathToBlockData: path.Join(pathToDataDir, "blockDatadir"),
		PathToChainData: path.Join(pathToDataDir, "chainDatadir"),
	}

	if dirExisted, _ := exists(r.PathToBlockData); !dirExisted {
		if err := os.Mkdir(r.PathToBlockData, os.ModeDir); err != nil {
			log.Fatalf("Failed to create dir=%s", r.PathToBlockData)
		}
	}

	if dirExisted, _ := exists(r.PathToChainData); !dirExisted {
		if err := os.Mkdir(r.PathToChainData, os.ModeDir); err != nil {
			log.Fatalf("Failed to create dir=%s", r.PathToChainData)
		}
	}

	return r
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

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
