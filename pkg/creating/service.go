package creating

import (
	"errors"
	"time"
)

const (
	// DefaultGenesisLastHash is default last genesis block hash if not given from genesis config
	DefaultGenesisLastHash = "0x000"

	// DefaultGenesisHash is default genesis hash if not given from genesis config
	DefaultGenesisHash = "0x000"
)

// ErrMissingLastBlockHash is used when new block is not provided last block hash
var ErrMissingLastBlockHash = errors.New("Missing last block hash")

// Service provides block creating operations
type Service interface {
	CreateGenesisBlock(genesisConfig GenesisConfig) (*Block, error)
	CreateBlock(lastHash *string, data []string) (*Block, error)
}

type service struct{}

// NewService creates a creating service with necessary dependencies
func NewService() Service {
	return &service{}
}

// CreateGenesisBlock returns the genesis block created from config
func (s *service) CreateGenesisBlock(genesisConfig GenesisConfig) (*Block, error) {
	// validations
	lastBlockHash := *genesisConfig.LastHash
	if genesisConfig.LastHash == nil {
		lastBlockHash = DefaultGenesisLastHash
	}
	blockHash := *genesisConfig.Hash
	if genesisConfig.Hash == nil {
		blockHash = DefaultGenesisHash
	}

	return createBlock(time.Now(), &lastBlockHash, &blockHash, *genesisConfig.Data), nil
}

// CreateBlock returns the new block with given lastHash, data
func (s *service) CreateBlock(lastHash *string, data []string) (*Block, error) {
	// validations
	if lastHash == nil {
		return nil, ErrMissingLastBlockHash
	}

	return createBlock(time.Now(), lastHash, nil, data), nil
}

func createBlock(timestamp time.Time, lastHash *string, hash *string, data []string) *Block {
	return &Block{
		Timestamp: timestamp,
		LastHash:  lastHash,
		Hash:      hash,
		Data:      data,
	}
}
