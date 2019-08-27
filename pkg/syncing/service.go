package syncing

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/knd/kndchain/pkg/mining"
)

// Service provides access to blockchain syncing operations
type Service interface {
	SyncBlockchain(nodeURL string) error
}

type service struct {
	m mining.Service
}

// NewService creates a syncing service with necessary dependencies
func NewService(m mining.Service) Service {
	return &service{m}
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
