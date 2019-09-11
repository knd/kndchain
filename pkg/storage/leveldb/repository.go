package leveldb

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"time"

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
	lastBlockHash         string
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

// ErrPersistTransaction indicates when there is error persisting transaction
var ErrPersistTransaction = errors.New("Failed to persist transaction")

// ErrPersistBlock indicates when there is error persisting block
var ErrPersistBlock = errors.New("Failed to persist block")

// ErrPersistBlockchain indicates where there is error persisting blockchain
var ErrPersistBlockchain = errors.New("Failed to persist blockchain")

// AddBlock adds mined block into blockchain
func (db *LevelDB) AddBlock(minedBlock *mining.Block) error {
	if minedBlock == nil {
		return ErrAddNilBlock
	}

	rBlock := toRepoBlock(minedBlock)

	// add to transaction db
	for _, tx := range rBlock.Data {
		txBytes, err := json.Marshal(tx)
		if err != nil {
			panic(err)
		}
		if err := db.transactionDB.Put(
			[]byte(fmt.Sprintf("%v", tx.Input.Timestamp)),
			txBytes,
			nil); err != nil {
			return ErrPersistTransaction
		}
	}

	// add to block db
	blockBytes, err := json.Marshal(rBlock)
	if err != nil {
		panic(err)
	}
	if err := db.blockDB.Put(
		[]byte(*rBlock.Hash),
		[]byte(blockBytes),
		nil); err != nil {
		return ErrPersistBlock
	}

	// add to chain db
	if err := db.chainDB.Put(
		[]byte(fmt.Sprintf("%v", rBlock.Timestamp)),
		[]byte(*rBlock.Hash),
		nil); err != nil {
		return ErrPersistBlockchain
	}

	db.lastBlockHash = *rBlock.LastHash

	return nil
}

func toRepoBlock(miningBlock *mining.Block) *Block {
	var transactions []Transaction
	for _, miningBlockTransaction := range miningBlock.Data {
		transactions = append(transactions, Transaction{
			ID:     miningBlockTransaction.ID,
			Output: miningBlockTransaction.Output,
			Input: Input{
				Timestamp: miningBlockTransaction.Input.Timestamp,
				Amount:    miningBlockTransaction.Input.Amount,
				Address:   miningBlockTransaction.Input.Address,
				Signature: miningBlockTransaction.Input.Signature,
			},
		})
	}
	return &Block{
		Timestamp:  miningBlock.Timestamp.Unix(),
		LastHash:   miningBlock.LastHash,
		Hash:       miningBlock.Hash,
		Nonce:      miningBlock.Nonce,
		Difficulty: miningBlock.Difficulty,
		Data:       transactions,
	}
}

// GetBlockCount returns the latest block count in blockchain
func (db *LevelDB) GetBlockCount() uint32 {
	// TODO: Implement this
	return 0
}

// GetLastBlock returns the last block in blockchain
func (db *LevelDB) GetLastBlock() listing.Block {
	var rBlock Block
	if db.lastBlockHash != "" {
		blockBytes, err := db.blockDB.Get(
			[]byte(db.lastBlockHash),
			nil)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(blockBytes, &rBlock)
		if err != nil {
			panic(err)
		}
	}

	return toListingBlock(rBlock)
}

func toListingBlock(b Block) listing.Block {
	var transactions []listing.Transaction
	for _, tx := range b.Data {
		transactions = append(transactions, listing.Transaction{
			ID:     tx.ID,
			Output: tx.Output,
			Input: listing.Input{
				Timestamp: tx.Input.Timestamp,
				Amount:    tx.Input.Amount,
				Address:   tx.Input.Address,
				Signature: tx.Input.Signature,
			},
		})
	}

	return listing.Block{
		Timestamp:  time.Unix(b.Timestamp, 0),
		LastHash:   b.LastHash,
		Hash:       b.Hash,
		Nonce:      b.Nonce,
		Difficulty: b.Difficulty,
		Data:       transactions,
	}
}

// GetBlockchain returns a list of blocks from genesis block
func (db *LevelDB) GetBlockchain() *listing.Blockchain {
	// TODO: Implement this

	return nil
}

// ReplaceChain replace the current blockchain with the newchain
func (db *LevelDB) ReplaceChain(newChain *mining.Blockchain) error {
	// TODO: Implement this

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
