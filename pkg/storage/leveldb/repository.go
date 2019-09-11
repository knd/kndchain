package leveldb

import (
	"errors"
	"log"
	"os"
	"path"

	"github.com/syndtr/goleveldb/leveldb"

	"github.com/knd/kndchain/pkg/listing"
	"github.com/knd/kndchain/pkg/mining"
)

// LevelDB keeps blockchain in local key-value db
type LevelDB struct {
	PathToTransactionData string
	PathToBlockData       string
	PathToChainData       string
	transactionDB         *leveldb.DB
	blockDB               *leveldb.DB
	chainDB               *leveldb.DB
}

// NewRepository creates a repository to interact with LevelDB
func NewRepository(pathToDataDir string) *LevelDB {
	r := &LevelDB{
		PathToTransactionData: path.Join(pathToDataDir, "transactionDatadir"),
		PathToBlockData:       path.Join(pathToDataDir, "blockDatadir"),
		PathToChainData:       path.Join(pathToDataDir, "chainDatadir"),
	}

	if dirExisted, _ := exists(r.PathToTransactionData); !dirExisted {
		if err := os.Mkdir(r.PathToTransactionData, os.ModeDir); err != nil {
			log.Fatalf("Failed to create dir=%s", r.PathToTransactionData)
		}
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

	transactionDB, err := leveldb.OpenFile(r.PathToTransactionData, nil)
	if err != nil {
		log.Fatalf("Failed to open leveldb#openfile dir=%s", r.PathToTransactionData)
	}
	r.transactionDB = transactionDB

	blockDB, err := leveldb.OpenFile(r.PathToBlockData, nil)
	if err != nil {
		log.Fatalf("Failed to open leveldb#openfile dir=%s", r.PathToBlockData)
	}
	r.blockDB = blockDB

	chainDB, err := leveldb.OpenFile(r.PathToChainData, nil)
	if err != nil {
		log.Fatalf("Failed to open leveldb#openfile dir=%s", r.PathToChainData)
	}
	r.chainDB = chainDB

	return r
}

// ErrAddNilBlock is used when no mined block is given to add
var ErrAddNilBlock = errors.New("Mined block is not given to add")

// AddBlock adds mined block into blockchain
func (db *LevelDB) AddBlock(minedBlock *mining.Block) error {
	if minedBlock == nil {
		return ErrAddNilBlock
	}

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

// Close closes any db connections
func (db *LevelDB) Close() error {
	err := db.transactionDB.Close()
	if err != nil {
		return err
	}
	err = db.blockDB.Close()
	if err != nil {
		return err
	}
	return db.chainDB.Close()
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
