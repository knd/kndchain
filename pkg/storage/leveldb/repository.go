package leveldb

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"sync"

	"github.com/knd/kndchain/pkg/listing"
	"github.com/knd/kndchain/pkg/mining"
	"github.com/syndtr/goleveldb/leveldb"
)

// LevelDB keeps blockchain in local key-value db
type LevelDB struct {
	PathToTransactionData string
	PathToBlockData       string
	PathToChainData       string
	transactionDB         *leveldb.DB
	blockDB               *leveldb.DB
	chainDB               *leveldb.DB
	currentBlockHash      string
	blockCount            uint32
	mutex                 *sync.Mutex
}

// NewRepository creates a repository to interact with LevelDB
func NewRepository(pathToDataDir string) *LevelDB {
	r := &LevelDB{
		PathToTransactionData: path.Join(pathToDataDir, "transactionDatadir"),
		PathToBlockData:       path.Join(pathToDataDir, "blockDatadir"),
		PathToChainData:       path.Join(pathToDataDir, "chainDatadir"),
		mutex:                 &sync.Mutex{},
	}

	if dirExisted, _ := exists(r.PathToTransactionData); !dirExisted {
		if err := os.Mkdir(r.PathToTransactionData, os.ModePerm); err != nil {
			log.Println(err)
			log.Fatalf("Failed to create dir=%s", r.PathToTransactionData)
		}
	}

	if dirExisted, _ := exists(r.PathToBlockData); !dirExisted {
		if err := os.Mkdir(r.PathToBlockData, os.ModePerm); err != nil {
			log.Println(err)
			log.Fatalf("Failed to create dir=%s", r.PathToBlockData)
		}
	}

	if dirExisted, _ := exists(r.PathToChainData); !dirExisted {
		if err := os.Mkdir(r.PathToChainData, os.ModePerm); err != nil {
			log.Println(err)
			log.Fatalf("Failed to create dir=%s", r.PathToChainData)
		}
	}

	transactionDB, err := leveldb.OpenFile(r.PathToTransactionData, nil)
	if err != nil {
		log.Fatalf("Failed to open leveldb#openfile dir=%s, %v", r.PathToTransactionData, err)
	}
	r.transactionDB = transactionDB

	blockDB, err := leveldb.OpenFile(r.PathToBlockData, nil)
	if err != nil {
		log.Fatalf("Failed to open leveldb#openfile dir=%s, %v", r.PathToBlockData, err)
	}
	r.blockDB = blockDB

	chainDB, err := leveldb.OpenFile(r.PathToChainData, nil)
	if err != nil {
		log.Fatalf("Failed to open leveldb#openfile dir=%s, %v", r.PathToChainData, err)
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

	db.mutex.Lock()
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
		[]byte(rBlock.Hash),
		[]byte(blockBytes),
		nil); err != nil {
		return ErrPersistBlock
	}

	// add to chain db
	if err := db.chainDB.Put(
		[]byte(fmt.Sprintf("%v", rBlock.Timestamp)),
		[]byte(rBlock.Hash),
		nil); err != nil {
		return ErrPersistBlockchain
	}

	db.currentBlockHash = rBlock.Hash
	db.blockCount++
	db.mutex.Unlock()

	log.Printf("Added block. Timestamp: %d, BlockHash=%s, Count=%d", rBlock.Timestamp, db.currentBlockHash, db.blockCount)

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
		Timestamp:  miningBlock.Timestamp,
		LastHash:   *miningBlock.LastHash,
		Hash:       *miningBlock.Hash,
		Nonce:      miningBlock.Nonce,
		Difficulty: miningBlock.Difficulty,
		Data:       transactions,
	}
}

// GetBlockCount returns the latest block count in blockchain
func (db *LevelDB) GetBlockCount() uint32 {
	if db.blockCount != 0 {
		return db.blockCount
	}

	var count uint32
	iter := db.chainDB.NewIterator(nil, nil)
	for iter.Next() {
		count++
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		panic(err)
	}
	return count
}

// GetLastBlock returns the last block in blockchain
func (db *LevelDB) GetLastBlock() listing.Block {
	if db.currentBlockHash != "" {
		db.mutex.Lock()
		lBlock := db.GetBlockByHash(db.currentBlockHash)
		if lBlock == nil {
			panic("No last block found by currentBlockHash")
		}
		db.mutex.Unlock()
		return *lBlock
	}

	db.mutex.Lock()
	iter := db.chainDB.NewIterator(nil, nil)
	var currentBlockHash string
	var lastBlockHashBytes []byte
	for iter.Next() {
		lastBlockHashBytes = iter.Value()
		currentBlockHash = string(lastBlockHashBytes[:])
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		panic(err)
	}
	db.currentBlockHash = currentBlockHash
	lBlock := db.GetBlockByHash(db.currentBlockHash)
	if lBlock == nil {
		panic(err)
	}
	db.mutex.Lock()

	return *lBlock
}

// GetBlockByHash returns block with given block hash
func (db *LevelDB) GetBlockByHash(hash string) *listing.Block {
	var rBlock Block
	blockBytes, err := db.blockDB.Get(
		[]byte(hash),
		nil)
	if err == leveldb.ErrNotFound {
		return nil
	}
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(blockBytes, &rBlock)
	if err != nil {
		panic(err)
	}

	lBlock := toListingBlock(rBlock)
	return &lBlock
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
		Timestamp:  b.Timestamp,
		LastHash:   &b.LastHash,
		Hash:       &b.Hash,
		Nonce:      b.Nonce,
		Difficulty: b.Difficulty,
		Data:       transactions,
	}
}

// GetBlockchain returns a list of blocks from genesis block
func (db *LevelDB) GetBlockchain() *listing.Blockchain {
	lBlockchain := &listing.Blockchain{}

	iter := db.chainDB.NewIterator(nil, nil)
	for iter.Next() {
		blockBytes, err := db.blockDB.Get(iter.Value(), nil)
		if err != nil {
			panic(err)
		}
		var rBlock Block
		err = json.Unmarshal(blockBytes, &rBlock)
		if err != nil {
			panic(err)
		}
		log.Printf("GetBlockchain(): LastHash=%s", rBlock.LastHash)
		log.Printf("GetBlockchain(): Hash=%s", rBlock.Hash)
		log.Println("---")
		lBlockchain.Chain = append(lBlockchain.Chain, toListingBlock(rBlock))
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		panic(err)
	}

	return lBlockchain
}

// ReplaceChain replace the current blockchain with the newchain
func (db *LevelDB) ReplaceChain(newChain *mining.Blockchain) error {
	db.mutex.Lock()
	err := db.DeleteAllData()
	if err != nil {
		return err
	}

	var lastMinedBlock mining.Block
	for _, minedBlock := range newChain.Chain {
		lastMinedBlock = minedBlock
		err = db.AddBlock(&minedBlock)
		if err != nil {
			panic(err)
		}
	}

	db.currentBlockHash = *lastMinedBlock.Hash
	db.blockCount = uint32(len(newChain.Chain))
	db.mutex.Unlock()

	return nil
}

// DeleteAllData delete everything
func (db *LevelDB) DeleteAllData() error {
	deleteDB(db.transactionDB)
	deleteDB(db.blockDB)
	deleteDB(db.chainDB)
	db.currentBlockHash = ""
	db.blockCount = 0
	log.Printf("Blockcount after delete=%d", db.GetBlockCount())
	return nil
}

func deleteDB(db *leveldb.DB) error {
	var txToDelete [][]byte
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		txToDelete = append(txToDelete, iter.Key())
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		panic(err)
	}
	for _, txTimestamp := range txToDelete {
		err = db.Delete(txTimestamp, nil)
		if err != nil {
			panic(err)
		}
	}
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
