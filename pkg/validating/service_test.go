package validating

import (
	"testing"
	"time"

	"github.com/knd/kndchain/pkg/hashing"
	"github.com/stretchr/testify/assert"
)

func TestService_IsInvalidChainWhenGenesisBlockIsInvalid(t *testing.T) {
	validatingService := NewService()
	lastHash := "0x123"
	hash := "0x456"
	blockchain := &Blockchain{
		chain: []Block{
			Block{
				Timestamp: time.Now(),
				LastHash:  &lastHash,
				Hash:      &hash,
				Data:      []string{"has initial tx which is not supposed to be"},
			},
		},
	}

	// perform test & verification
	assert.False(t, validatingService.IsValidChain(blockchain))
}

func TestService_IsInvalidChainWhenLastHashIsTampered(t *testing.T) {
	validatingService := NewService()
	genesisTimestamp := time.Now()
	lastHash := "0x123"
	hash := "0x456"
	tamperedLashHash := "tampered"
	blockchain := &Blockchain{
		chain: []Block{
			Block{
				Timestamp: genesisTimestamp,
				LastHash:  &lastHash,
				Hash:      &hash,
				Data:      []string{},
			},
			Block{
				Timestamp: time.Now(),
				LastHash:  &tamperedLashHash,
				Hash:      &hash,
				Data:      []string{},
			},
		},
	}

	// perform test & verification
	assert.False(t, validatingService.IsValidChain(blockchain))
}

func TestService_IsInvalidChainWhenTimestampIsNotInOrder(t *testing.T) {
	validatingService := NewService()

	genesisLastHash := "0x123"
	genesisHash := "0x456"
	genesisTimestamp := time.Now()
	timestamp1 := time.Now().Add(time.Duration(100))
	timestamp2 := time.Now().Add(time.Duration(200))

	genesisBlock := Block{
		Timestamp: genesisTimestamp,
		LastHash:  &genesisLastHash,
		Hash:      &genesisHash,
		Data:      []string{},
	}

	blockA := Block{
		Timestamp: timestamp2,
		LastHash:  &genesisHash,
		Data:      []string{"txA"},
	}
	blockAHash := hashing.SHA256Hash(timestamp2, genesisHash, []string{"txA"})
	blockA.Hash = &blockAHash

	blockB := Block{
		Timestamp: timestamp1,
		LastHash:  &blockAHash,
		Data:      []string{"txB"},
	}
	blockBHash := hashing.SHA256Hash(timestamp1, blockAHash, []string{"txB"})
	blockB.Hash = &blockBHash

	blockchain := &Blockchain{chain: []Block{genesisBlock, blockA, blockB}}

	// perform test & verification
	assert.False(t, validatingService.IsValidChain(blockchain))
}

func TestService_IsValidChainWhenChainContainsOnlyValidBlocks(t *testing.T) {
	validatingService := NewService()

	genesisLastHash := "0x123"
	genesisHash := "0x456"
	genesisTimestamp := time.Now()
	timestamp1 := time.Now().Add(time.Duration(100))
	timestamp2 := time.Now().Add(time.Duration(200))

	genesisBlock := Block{
		Timestamp: genesisTimestamp,
		LastHash:  &genesisLastHash,
		Hash:      &genesisHash,
		Data:      []string{},
	}

	blockA := Block{
		Timestamp: timestamp1,
		LastHash:  &genesisHash,
		Data:      []string{"txA"},
	}
	blockAHash := hashing.SHA256Hash(timestamp1, genesisHash, blockA.Data)
	blockA.Hash = &blockAHash

	blockB := Block{
		Timestamp: timestamp2,
		LastHash:  &blockAHash,
		Data:      []string{"txB"},
	}
	blockBHash := hashing.SHA256Hash(timestamp2, blockAHash, blockB.Data)
	blockB.Hash = &blockBHash

	blockchain := &Blockchain{chain: []Block{genesisBlock, blockA, blockB}}

	// perform test & verification
	assert.True(t, validatingService.IsValidChain(blockchain))
}
