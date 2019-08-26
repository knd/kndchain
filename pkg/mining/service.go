package mining

import (
	"errors"
	"strings"
	"time"

	"github.com/knd/kndchain/pkg/hashing"
	"github.com/knd/kndchain/pkg/listing"
	"github.com/knd/kndchain/pkg/validating"
)

const (
	// MineRate (1000 milliseconds) adjusts the difficulty of mining operation
	MineRate int = 1000
)

const (
	// DefaultGenesisLastHash is default last genesis block hash if not given from genesis config
	DefaultGenesisLastHash = "0x000"

	// DefaultGenesisHash is default genesis hash if not given from genesis config
	DefaultGenesisHash = "0x000"

	// DefaultGenesisDifficulty is default difficulty in genesis block
	DefaultGenesisDifficulty uint32 = 3

	// DefaultGenesisNonce is default nonce in genesis block
	DefaultGenesisNonce uint32 = 0
)

// ErrMissingLastBlock is used when new block is not provided last block hash
var ErrMissingLastBlock = errors.New("Missing last block")

// ErrInvalidChain is used when trying to replace an invalid chain
var ErrInvalidChain = errors.New("Invalid chain is given to replace original chain")

// ErrShorterChain is used when trying to replace a shorter chain
var ErrShorterChain = errors.New("Shorter chain is given to replace original chain")

// Service provides block creating operations
type Service interface {
	MineNewBlock(lastBlock *Block, data []string) (*Block, error)
	AddBlock(minedBlock *Block) error
	ReplaceChain(newChain *Blockchain) error
}

// Repository provides access to in-memory blockchain
type Repository interface {
	// AddBlock adds a minedBlock into blockchain
	AddBlock(minedBlock *Block) error
	ReplaceChain(newChain *Blockchain) error
}

type service struct {
	blockchain Repository
	listing    listing.Service
	validating validating.Service
}

// NewService creates a creating service with necessary dependencies
func NewService(r Repository, l listing.Service, v validating.Service, c *GenesisConfig) Service {
	newS := &service{r, l, v}

	// Do not have genesis block
	// if l.GetBlockCount() == 0 {
	// 	var genesisBlock *Block
	// 	var err error
	// 	if genesisBlock, err = CreateGenesisBlock(c); err != nil {
	// 		log.Fatal("Cannot create genesis block")
	// 	}
	// 	if err = newS.AddBlock(genesisBlock); err != nil {
	// 		log.Fatal("cannot add genesis block to blockchain")
	// 	}
	// }

	return newS
}

// CreateGenesisBlock returns the genesis block created from config
func CreateGenesisBlock(genesisConfig *GenesisConfig) (*Block, error) {
	// validations
	var lastBlockHash string
	if genesisConfig == nil || genesisConfig.LastHash == nil {
		lastBlockHash = DefaultGenesisLastHash
	} else {
		lastBlockHash = *genesisConfig.LastHash
	}

	var blockHash string
	if genesisConfig == nil || genesisConfig.Hash == nil {
		blockHash = DefaultGenesisHash
	} else {
		blockHash = *genesisConfig.Hash
	}

	var blockDifficulty uint32
	if genesisConfig == nil {
		blockDifficulty = DefaultGenesisDifficulty
	} else {
		blockDifficulty = genesisConfig.Difficulty
	}

	var blockNonce uint32
	if genesisConfig == nil {
		blockNonce = DefaultGenesisNonce
	} else {
		blockNonce = genesisConfig.Nonce
	}

	var data []string
	if genesisConfig != nil && genesisConfig.Data != nil {
		data = *genesisConfig.Data
	}

	return yieldBlock(time.Now(), &lastBlockHash, &blockHash, data, blockNonce, blockDifficulty), nil
}

// AdjustBlockDifficulty adjusts the difficulty in the current mining block
func AdjustBlockDifficulty(lastBlock Block, blockTimestamp time.Time) uint32 {
	if blockTimestamp.Sub(lastBlock.Timestamp) < (time.Duration(MineRate) * time.Millisecond) {
		return lastBlock.Difficulty + 1
	} else if blockTimestamp.Sub(lastBlock.Timestamp) > (time.Duration(MineRate) * time.Millisecond) {
		if lastBlock.Difficulty <= 1 {
			return 1
		}
		return lastBlock.Difficulty - 1
	}

	return lastBlock.Difficulty
}

// MineNewBlock returns a new block
func (s *service) MineNewBlock(lastBlock *Block, data []string) (*Block, error) {
	// validations
	if lastBlock == nil {
		return nil, ErrMissingLastBlock
	}

	difficulty := lastBlock.Difficulty
	var nonce uint32
	var timestamp time.Time
	var hash string
	for {
		nonce++
		timestamp = time.Now()
		difficulty = AdjustBlockDifficulty(*lastBlock, timestamp)
		hash = hashing.SHA256Hash(timestamp, *lastBlock.Hash, data, nonce, difficulty)
		if hash[:difficulty] == strings.Repeat("0", int(difficulty)) {
			break
		}
	}

	return yieldBlock(timestamp, lastBlock.Hash, &hash, data, nonce, difficulty), nil
}

// AddBlock adds a minedBlock into blockchain
func (s *service) AddBlock(minedBlock *Block) error {
	// TODO: validations of minedBlock
	return s.blockchain.AddBlock(minedBlock)
}

func yieldBlock(timestamp time.Time, lastHash *string, hash *string, data []string, nonce uint32, difficulty uint32) *Block {
	return &Block{
		Timestamp:  timestamp,
		LastHash:   lastHash,
		Hash:       hash,
		Data:       data,
		Nonce:      nonce,
		Difficulty: difficulty,
	}
}

func (s *service) ReplaceChain(newChain *Blockchain) error {
	// validations
	if newChain == nil {
		return ErrInvalidChain
	}
	if uint32(len(newChain.Chain)) <= s.listing.GetBlockCount() {
		return ErrShorterChain
	}
	vChain := toValidatingChain(newChain)
	if !s.validating.IsValidChain(vChain) {
		return ErrInvalidChain
	}

	return s.blockchain.ReplaceChain(newChain)
}

func toValidatingChain(newChain *Blockchain) *validating.Blockchain {
	vBlockchain := &validating.Blockchain{}
	for _, block := range newChain.Chain {
		vBlock := &validating.Block{
			Timestamp: block.Timestamp,
			LastHash:  block.LastHash,
			Hash:      block.Hash,
			Data:      block.Data,
		}
		vBlockchain.Chain = append(vBlockchain.Chain, *vBlock)
	}
	return vBlockchain
}
