package main

import (
	"fmt"
	"time"

	"github.com/knd/kndchain/pkg/calculating"
	"github.com/knd/kndchain/pkg/listing"
	"github.com/knd/kndchain/pkg/mining"
	"github.com/knd/kndchain/pkg/storage/memory"
	"github.com/knd/kndchain/pkg/validating"
)

// Type defines available storage types
type Type int

const (
	// JSON stores data in JSON files saved on disk
	JSON Type = iota

	// Memory stores data in memory
	Memory
)

func toMiningTransactions(data []listing.Transaction) []mining.Transaction {
	var mTxs []mining.Transaction
	for _, transaction := range data {
		mTxs = append(mTxs, mining.Transaction{
			ID:     transaction.ID,
			Output: transaction.Output,
			Input: mining.Input{
				Timestamp: transaction.Input.Timestamp,
				Amount:    transaction.Input.Amount,
				Address:   transaction.Input.Address,
				Signature: transaction.Input.Signature,
			},
		})
	}
	return mTxs
}

func main() {
	// set up storage
	storageType := Memory

	var miner mining.Service
	var lister listing.Service
	var validator validating.Service

	switch storageType {
	case Memory:
		storage := memory.NewRepository()

		lister = listing.NewService(storage)
		validator = validating.NewService(lister, calculating.NewService())
		miner = mining.NewService(storage, lister, validator, nil)
	}

	fmt.Println("Staring now")

	genesisBlock, _ := mining.CreateGenesisBlock(nil)
	miner.AddBlock(genesisBlock)

	var durations []float64
	for i := 0; i < 10000; i++ {
		lastBlock := lister.GetLastBlock()

		mb := &mining.Block{
			Timestamp:  lastBlock.Timestamp,
			LastHash:   lastBlock.LastHash,
			Hash:       lastBlock.Hash,
			Data:       toMiningTransactions(lastBlock.Data),
			Nonce:      lastBlock.Nonce,
			Difficulty: lastBlock.Difficulty,
		}

		newB, _ := miner.MineNewBlock(mb, []mining.Transaction{mining.Transaction{ID: "dummy-tx"}})

		miner.AddBlock(newB)

		durationDiff := newB.Timestamp.Sub(lastBlock.Timestamp) // in nano
		durationDiffInMillis := float64(durationDiff) / float64(time.Millisecond)

		durations = append(durations, durationDiffInMillis)

		var sumDuration float64
		for _, duration := range durations {
			sumDuration = sumDuration + duration
		}
		averageDuration := float64(sumDuration) / float64(len(durations))

		fmt.Printf(
			"Time to mine block: %.2f ms. Difficulty: %d. Average time: %.2f ms", durationDiffInMillis,
			newB.Difficulty,
			averageDuration,
		)
		fmt.Println()
	}
}
