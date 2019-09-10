package mining

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/knd/kndchain/pkg/hashing"
	"github.com/knd/kndchain/pkg/listing"
	"github.com/knd/kndchain/pkg/validating"
)

// ErrMissingLastBlock is used when new block is not provided last block hash
var ErrMissingLastBlock = errors.New("Missing last block")

// ErrInvalidChain is used when trying to replace an invalid chain
var ErrInvalidChain = errors.New("Invalid chain is given to replace original chain")

// ErrInvalidTransactions is used when trying to replace chain with invalid transactions
var ErrInvalidTransactions = errors.New("Invalid transactions")

// ErrShorterChain is used when trying to replace a shorter chain
var ErrShorterChain = errors.New("Current chain is the longest, Incoming chain is no longer, No replacement")

// Service provides block creating operations
type Service interface {
	MineNewBlock(lastBlock *Block, data []Transaction) (*Block, error)
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
	MineRate   int64
}

// NewService creates a creating service with necessary dependencies
func NewService(r Repository, l listing.Service, v validating.Service, mineRate int64) Service {
	return &service{r, l, v, mineRate}
}

// CreateGenesisBlock returns the genesis block created from config
func CreateGenesisBlock(genesisLastHash string, genesisHash string, genesisDifficulty uint32, genesisNonce uint32) (*Block, error) {
	data := []Transaction{}

	return yieldBlock(time.Now(), &genesisLastHash, &genesisHash, data, genesisNonce, genesisDifficulty), nil
}

func adjustBlockDifficulty(lastBlock Block, blockTimestamp time.Time, mineRate int64) uint32 {
	if blockTimestamp.Sub(lastBlock.Timestamp) < (time.Duration(mineRate) * time.Millisecond) {
		return lastBlock.Difficulty + 1
	} else if blockTimestamp.Sub(lastBlock.Timestamp) > (time.Duration(mineRate) * time.Millisecond) {
		if lastBlock.Difficulty <= 1 {
			return 1
		}
		return lastBlock.Difficulty - 1
	}

	return lastBlock.Difficulty
}

// MineNewBlock returns a new block
func (s *service) MineNewBlock(lastBlock *Block, data []Transaction) (*Block, error) {
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
		difficulty = adjustBlockDifficulty(*lastBlock, timestamp, s.MineRate)
		hash = hashing.SHA256Hash(timestamp.Unix(), *lastBlock.Hash, data, nonce, difficulty)
		if hexStringToBinary(hash)[:difficulty] == strings.Repeat("0", int(difficulty)) {
			break
		}
	}

	return yieldBlock(timestamp, lastBlock.Hash, &hash, data, nonce, difficulty), nil
}

// HexStringToBinary converts the hex string to binary string representation
func hexStringToBinary(s string) string {
	res := ""
	b, _ := hex.DecodeString(s)
	for _, c := range b {
		binary, _ := strconv.Atoi(fmt.Sprintf("%.b", c))
		res = fmt.Sprintf("%s%s", res, fmt.Sprintf("%08d", binary))
	}
	return res
}

// AddBlock adds a minedBlock into blockchain
func (s *service) AddBlock(minedBlock *Block) error {
	if minedBlock == nil {
		return errors.New("No block provided to add")
	}
	return s.blockchain.AddBlock(minedBlock)
}

func yieldBlock(timestamp time.Time, lastHash *string, hash *string, data []Transaction, nonce uint32, difficulty uint32) *Block {
	return &Block{
		Timestamp:  timestamp,
		LastHash:   lastHash,
		Hash:       hash,
		Data:       data,
		Nonce:      nonce,
		Difficulty: difficulty,
	}
}

// ReplaceChain replaces valid incoming chain with existing chain
func (s *service) ReplaceChain(newChain *Blockchain) error {
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
	if valid, err := s.validating.ContainsValidTransactions(vChain); !valid || err != nil {
		log.Printf("MiningService#ReplaceChain: Failed to replace chain %v", err)
		return ErrInvalidTransactions
	}

	return s.blockchain.ReplaceChain(newChain)
}
