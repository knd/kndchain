package syncing

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/knd/kndchain/pkg/mining"
	"github.com/knd/kndchain/pkg/wallet"
)

// Service provides access to blockchain/transactions syncing operations
type Service interface {
	SyncBlockchain(nodeURL string) error
	SyncTransactionPool(nodeURL string) error
}

type service struct {
	m mining.Service
	p wallet.TransactionPool
}

// NewService creates a syncing service with necessary dependencies
func NewService(m mining.Service, p wallet.TransactionPool) Service {
	return &service{m, p}
}

// SyncBlockchain obtains the full blockchain from nodeEndpoint url
func (s *service) SyncBlockchain(nodeURL string) error {
	req, _ := http.NewRequest("GET", nodeURL, nil)
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()

	var bc Blockchain
	json.NewDecoder(resp.Body).Decode(&bc)

	var blocks []mining.Block
	for _, b := range bc.Chain {
		blocks = append(blocks, mining.Block{
			Timestamp:  b.Timestamp,
			LastHash:   b.LastHash,
			Hash:       b.Hash,
			Data:       b.Data,
			Nonce:      b.Nonce,
			Difficulty: b.Difficulty,
		})
	}
	mbc := &mining.Blockchain{Chain: blocks}

	return s.m.ReplaceChain(mbc)
}

// SyncTransactionPool obtains the full transaction pool from nodeEndpoint url
func (s *service) SyncTransactionPool(nodeURL string) error {
	req, _ := http.NewRequest("GET", nodeURL, nil)
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()

	var incomingPool TransactionPool
	json.NewDecoder(resp.Body).Decode(&incomingPool)

	pool := make(map[string]wallet.Transaction)
	for _, t := range incomingPool {
		tx := &wallet.Tx{
			ID: t.ID,
			Input: wallet.Input{
				Timestamp: t.Input.Timestamp,
				Amount:    t.Input.Amount,
				Address:   t.Input.Address,
				Signature: t.Input.Signature,
			},
			Output: map[string]uint64(t.Output),
		}
		pool[tx.ID] = tx
	}

	return s.p.SetPool(pool)
}
