package mining

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
var ErrMissingLastBlock = errors.New("Missing last block")

// Service provides block creating operations
type Service interface {
	CreateGenesisBlock(genesisConfig GenesisConfig) (*Block, error)
	MineNewBlock(lastBlock *Block, data []string) (*Block, error)
}

type service struct{}

// NewService creates a creating service with necessary dependencies
func NewService() Service {
	return &service{}
}

// CreateGenesisBlock returns the genesis block created from config
func (s *service) CreateGenesisBlock(genesisConfig GenesisConfig) (*Block, error) {
	// validations
	var lastBlockHash string
	if genesisConfig.LastHash == nil {
		lastBlockHash = DefaultGenesisLastHash
	} else {
		lastBlockHash = *genesisConfig.LastHash
	}

	var blockHash string
	if genesisConfig.Hash == nil {
		blockHash = DefaultGenesisHash
	} else {
		blockHash = *genesisConfig.Hash
	}

	var data []string
	if genesisConfig.Data != nil {
		data = *genesisConfig.Data
	}

	return mineBlock(time.Now(), &lastBlockHash, &blockHash, data), nil
}

// MineNewBlock returns a new block
func (s *service) MineNewBlock(lastBlock *Block, data []string) (*Block, error) {
	// validations
	if lastBlock == nil {
		return nil, ErrMissingLastBlock
	}

	randomHash := "randomHash" // TODO: Implement hash
	return mineBlock(time.Now(), lastBlock.Hash, &randomHash, data), nil
}

func mineBlock(timestamp time.Time, lastHash *string, hash *string, data []string) *Block {
	return &Block{
		Timestamp: timestamp,
		LastHash:  lastHash,
		Hash:      hash,
		Data:      data,
	}
}
