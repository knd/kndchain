package main

import (
	"log"

	"github.com/knd/kndchain/pkg/listing"
	"github.com/knd/kndchain/pkg/mining"
	"github.com/knd/kndchain/pkg/networking/pubsub"
	"github.com/knd/kndchain/pkg/wallet"
)

// Miner provides entry to mining actions
type Miner interface {
	Mine() (*mining.Block, error)
}

type miner struct {
	service              mining.Service
	lister               listing.Service
	transactionPool      wallet.TransactionPool
	wal                  wallet.Wallet
	comm                 pubsub.Service
	rewardTxInputAddress string
	rewardAmount         uint64
}

// NewMiner creates a miner with necessary dependencies
func NewMiner(s mining.Service, l listing.Service, p wallet.TransactionPool, w wallet.Wallet, c pubsub.Service, rewardTxInputAddress string, rewardAmount uint64) Miner {
	return &miner{s, l, p, w, c, rewardTxInputAddress, rewardAmount}
}

func (m *miner) Mine() (*mining.Block, error) {
	validTransactions := m.transactionPool.ValidTransactions()
	rewardTransaction, _ := wallet.CreateRewardTransaction(m.wal, m.rewardTxInputAddress, m.rewardAmount)

	validTransactions = append(validTransactions, rewardTransaction)

	lastBlock := m.lister.GetLastBlock()
	mb := &mining.Block{
		Timestamp:  lastBlock.Timestamp,
		LastHash:   lastBlock.LastHash,
		Hash:       lastBlock.Hash,
		Data:       fromListingtoMiningTransactions(lastBlock.Data),
		Nonce:      lastBlock.Nonce,
		Difficulty: lastBlock.Difficulty,
	}
	minedBlock, err := m.service.MineNewBlock(mb, fromPooltoMiningTransactions(validTransactions))
	if err != nil {
		log.Printf("Failed to create mined block: %s", err.Error())
		return nil, err
	}

	err = m.service.AddBlock(minedBlock)
	if err != nil {
		log.Printf("Failed to add block to chain: %s", err.Error())
		return nil, err
	}

	err = m.comm.BroadcastBlockchain(m.lister.GetBlockchain())
	if err != nil {
		log.Printf("Failed to broadcast blockchain: %s", err.Error())
		return minedBlock, err
	}

	err = m.transactionPool.Clear()
	if err != nil {
		log.Printf("Failed to clear transaction pool: %s", err.Error())
		return minedBlock, err
	}

	return minedBlock, nil
}

func fromListingtoMiningTransactions(data []listing.Transaction) []mining.Transaction {
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

func fromPooltoMiningTransactions(data []wallet.Transaction) []mining.Transaction {
	var mTxs []mining.Transaction
	for _, transaction := range data {
		mTxs = append(mTxs, mining.Transaction{
			ID:     transaction.GetID(),
			Output: transaction.GetOutput(),
			Input: mining.Input{
				Timestamp: transaction.GetInput().Timestamp,
				Amount:    transaction.GetInput().Amount,
				Address:   transaction.GetInput().Address,
				Signature: transaction.GetInput().Signature,
			},
		})
	}
	return mTxs
}
