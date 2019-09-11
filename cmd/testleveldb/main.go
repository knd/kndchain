package main

import (
	"log"
	"time"

	"github.com/knd/kndchain/pkg/mining"
	"github.com/knd/kndchain/pkg/storage/leveldb"
)

func scenarioA(repository *leveldb.LevelDB) {
	lastHashA := "0x000"
	hashA := "0x000"
	blockA := mining.Block{
		Timestamp:  time.Now(),
		LastHash:   &lastHashA,
		Hash:       &hashA,
		Data:       []mining.Transaction{},
		Nonce:      0,
		Difficulty: 15,
	}
	lastHashB := "blockAHash"
	hashB := "hashB"
	blockB := mining.Block{
		Timestamp: time.Now(),
		LastHash:  &lastHashB,
		Hash:      &hashB,
		Data: []mining.Transaction{
			mining.Transaction{
				ID:     "111",
				Output: map[string]uint64{"0x123": 10},
				Input: mining.Input{
					Timestamp: time.Now().Unix(),
					Amount:    100,
					Address:   "0x111",
					Signature: "abc",
				},
			},
		},
		Nonce:      1,
		Difficulty: 16,
	}

	repository.AddBlock(&blockA)
	repository.AddBlock(&blockB)

	lBlock := repository.GetLastBlock()
	log.Printf("Timestamp=%d", lBlock.Timestamp.Unix())
	log.Printf("LastHash=%s", *lBlock.LastHash)
	log.Printf("Hash=%s", *lBlock.Hash)
	log.Printf("Nonce=%d", lBlock.Nonce)
	log.Printf("Difficulty=%d", lBlock.Difficulty)
	log.Printf("len(data)=%d", len(lBlock.Data))
	if len(lBlock.Data) > 0 {
		for _, tx := range lBlock.Data {
			log.Printf("ID=%s", tx.ID)
			log.Printf("Output=%v", tx.Output)
			log.Printf("Input=%v", tx.Input)
		}
	}

	log.Printf("BlockCount=%d", repository.GetBlockCount())

	repository.DeleteAllData()
}

func scenarioB(repository *leveldb.LevelDB) {
	lastHashA := "0x000"
	hashA := "0x000"
	blockA := mining.Block{
		Timestamp:  time.Now(),
		LastHash:   &lastHashA,
		Hash:       &hashA,
		Data:       []mining.Transaction{},
		Nonce:      0,
		Difficulty: 15,
	}
	lastHashB := "blockAHash"
	hashB := "hashB"
	blockB := mining.Block{
		Timestamp: time.Now(),
		LastHash:  &lastHashB,
		Hash:      &hashB,
		Data: []mining.Transaction{
			mining.Transaction{
				ID:     "111",
				Output: map[string]uint64{"0x123": 10},
				Input: mining.Input{
					Timestamp: time.Now().Unix(),
					Amount:    100,
					Address:   "0x111",
					Signature: "abc",
				},
			},
		},
		Nonce:      1,
		Difficulty: 16,
	}
	repository.AddBlock(&blockA)
	repository.AddBlock(&blockB)

	lastHashC := "0x000C"
	hashC := "0x000C"
	blockC := mining.Block{
		Timestamp:  time.Now(),
		LastHash:   &lastHashC,
		Hash:       &hashC,
		Data:       []mining.Transaction{},
		Nonce:      1,
		Difficulty: 7,
	}
	lastHashD := "blockDHash"
	hashD := "hashD"
	blockD := mining.Block{
		Timestamp: time.Now(),
		LastHash:  &lastHashD,
		Hash:      &hashD,
		Data: []mining.Transaction{
			mining.Transaction{
				ID:     "DDD",
				Output: map[string]uint64{"0xDDD": 100},
				Input: mining.Input{
					Timestamp: time.Now().Unix(),
					Amount:    999,
					Address:   "DDD0x",
					Signature: "dddddd",
				},
			},
		},
		Nonce:      0,
		Difficulty: 1,
	}

	mBlockchain := &mining.Blockchain{
		Chain: []mining.Block{blockC, blockD},
	}
	repository.ReplaceChain(mBlockchain)

	lBlock := repository.GetLastBlock()
	log.Printf("Timestamp=%d", lBlock.Timestamp.Unix())
	log.Printf("LastHash=%s", *lBlock.LastHash)
	log.Printf("Hash=%s", *lBlock.Hash)
	log.Printf("Nonce=%d", lBlock.Nonce)
	log.Printf("Difficulty=%d", lBlock.Difficulty)
	log.Printf("len(data)=%d", len(lBlock.Data))
	if len(lBlock.Data) > 0 {
		for _, tx := range lBlock.Data {
			log.Printf("ID=%s", tx.ID)
			log.Printf("Output=%v", tx.Output)
			log.Printf("Input=%v", tx.Input)
		}
	}

	log.Printf("BlockCount=%d", repository.GetBlockCount())

	repository.DeleteAllData()
}

func main() {
	repository := leveldb.NewRepository("/Users/knd/kndchainDatadir")
	defer repository.Close()
	// scenarioA(repository)
	scenarioB(repository)
}
