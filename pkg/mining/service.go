package mining

import (
	"errors"
	"log"
	"time"

	"github.com/knd/kndchain/pkg/hashing"
	"github.com/knd/kndchain/pkg/listing"
)

const (
	// DefaultGenesisLastHash is default last genesis block hash if not given from genesis config
	DefaultGenesisLastHash = "0x000"

	// DefaultGenesisHash is default genesis hash if not given from genesis config
	DefaultGenesisHash = "0x000"
)

// ErrMissingLastBlock is used when new block is not provided last block hash
var ErrMissingLastBlock = errors.New("Missing last block")

// Service provides block creating operations
type Service interface {
	MineNewBlock(lastBlock *Block, data []string) (*Block, error)
	AddBlock(minedBlock *Block) error
}

// Repository provides access to in-memory blockchain
type Repository interface {
	// AddBlock adds a minedBlock into blockchain
	AddBlock(minedBlock *Block) error
}

type service struct {
	blockchain Repository
	listing    listing.Service
}

// NewService creates a creating service with necessary dependencies
func NewService(r Repository, l listing.Service) Service {
	newS := &service{r, l}

	// Do not have genesis block
	if l.GetBlockCount() == 0 {
		var genesisBlock *Block
		var err error
		if genesisBlock, err = CreateGenesisBlock(GenesisConfig{}); err != nil {
			log.Fatal("Cannot create genesis block")
		}
		if err = newS.AddBlock(genesisBlock); err != nil {
			log.Fatal("cannot add genesis block to blockchain")
		}
	}

	return newS
}

// CreateGenesisBlock returns the genesis block created from config
func CreateGenesisBlock(genesisConfig GenesisConfig) (*Block, error) {
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

	timestamp := time.Now()
	hash := hashing.SHA256Hash(timestamp, *lastBlock.Hash, data)

	return mineBlock(timestamp, lastBlock.Hash, &hash, data), nil
}

// AddBlock adds a minedBlock into blockchain
func (s *service) AddBlock(minedBlock *Block) error {
	// TODO: validations of minedBlock
	return s.blockchain.AddBlock(minedBlock)
}

func mineBlock(timestamp time.Time, lastHash *string, hash *string, data []string) *Block {
	return &Block{
		Timestamp: timestamp,
		LastHash:  lastHash,
		Hash:      hash,
		Data:      data,
	}
}
