package mining

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/knd/kndchain/pkg/hashing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateGenesisBlock(t *testing.T) {
	assert := assert.New(t)

	t.Run("creates default genesis block", func(t *testing.T) {
		jsonData := "{}"

		var genesisConfig GenesisConfig
		if err := json.Unmarshal([]byte(jsonData), &genesisConfig); err != nil {
			t.FailNow()
		}

		// perform test
		genesisBlock, err := CreateGenesisBlock(&genesisConfig)

		// test verification
		assert.Nil(err)
		assert.NotEmpty(genesisBlock.Timestamp)
		assert.Equal(*genesisBlock.LastHash, DefaultGenesisLastHash)
		assert.Equal(*genesisBlock.Hash, DefaultGenesisHash)
		assert.Empty(genesisBlock.Data)
	})

	t.Run("creates genesis block with given config", func(t *testing.T) {
		jsonData := `{ "lastHash": "0x123", "hash": "0x456", "data": ["tx1", "tx2"] }`

		var genesisConfig GenesisConfig
		if err := json.Unmarshal([]byte(jsonData), &genesisConfig); err != nil {
			t.FailNow()
		}

		// perform test
		genesisBlock, err := CreateGenesisBlock(&genesisConfig)

		// test verification
		assert.Nil(err)
		assert.NotEmpty(genesisBlock.Timestamp)
		assert.Equal("0x123", *genesisBlock.LastHash)
		assert.Equal("0x456", *genesisBlock.Hash)
		assert.Equal([]string{"tx1", "tx2"}, genesisBlock.Data)
	})
}

func TestService(t *testing.T) {
	assert := assert.New(t)
	var miningService Service
	var mockedRepository *MockedRepository
	var mockedListing *MockedListing
	var mockedValidating *MockedValidating

	beforeEach := func() {
		mockedRepository = new(MockedRepository)
		mockedListing = new(MockedListing)
		mockedValidating = new(MockedValidating)
		miningService = NewService(mockedRepository, mockedListing, mockedValidating, nil)
	}

	t.Run("mines new block", func(t *testing.T) {
		beforeEach()
		mockedListing.On("GetBlockCount").Return(1)
		lastHash := "0x123"
		hash := "0x456"
		lastBlock := Block{
			Timestamp: time.Now(),
			LastHash:  &lastHash,
			Hash:      &hash,
			Data:      []string{"tx1"},
		}
		data := []string{"tx2"}

		// perform test
		newBlock, err := miningService.MineNewBlock(&lastBlock, data)

		// test verification
		assert.Nil(err)
		assert.NotEmpty(newBlock.Timestamp)
		assert.Equal("0x456", *newBlock.LastHash)
		assert.Equal(*newBlock.Hash, hashing.SHA256Hash(data, newBlock.Timestamp, *lastBlock.Hash))
		assert.Equal(data, newBlock.Data)
	})

	t.Run("adds block to chain", func(t *testing.T) {
		beforeEach()
		mockedListing.On("GetBlockCount").Return(1)
		LastHash := "0x123"
		Hash := "0x456"
		minedBlock := &Block{
			Timestamp: time.Now(),
			LastHash:  &LastHash,
			Hash:      &Hash,
			Data:      []string{"tx1"},
		}
		mockedRepository.On("AddBlock", minedBlock).Return(nil)

		// perform test
		miningService.AddBlock(minedBlock)

		// test verification
		mockedRepository.AssertExpectations(t)
	})

	t.Run("replaces with nil chain", func(t *testing.T) {
		beforeEach()

		// perform test
		err := miningService.ReplaceChain(nil)

		// test verification
		assert.Equal(err, ErrInvalidChain)
		mockedRepository.AssertNotCalled(t, "ReplaceChain")
	})

	t.Run("replaces with shorter chain", func(t *testing.T) {
		beforeEach()
		mockedListing.On("GetBlockCount").Return(5)

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

		blockchain := &Blockchain{Chain: []Block{genesisBlock, blockA, blockB}}

		// perform test
		err := miningService.ReplaceChain(blockchain)

		// test verification
		assert.Equal(err, ErrShorterChain)
		mockedRepository.AssertNotCalled(t, "ReplaceChain")
	})

	t.Run("replaces with longer valid chain", func(t *testing.T) {
		beforeEach()
		mockedListing.On("GetBlockCount").Return(2)
		mockedValidating.On("IsValidChain", mock.Anything).Return(true)

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

		blockchain := &Blockchain{Chain: []Block{genesisBlock, blockA, blockB}}
		mockedRepository.On("ReplaceChain", blockchain).Return(nil)

		// perform test
		err := miningService.ReplaceChain(blockchain)

		// test verification
		assert.Nil(err)
		mockedRepository.AssertCalled(t, "ReplaceChain", blockchain)
	})

	t.Run("replaces with longer invalid chain", func(t *testing.T) {
		beforeEach()
		mockedListing.On("GetBlockCount").Return(2)
		mockedValidating.On("IsValidChain", mock.Anything).Return(false)

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
			Data:      []string{"txB", "double spend"},
		}
		blockBHash := hashing.SHA256Hash(timestamp2, blockAHash, blockB.Data)
		blockB.Hash = &blockBHash

		blockchain := &Blockchain{Chain: []Block{genesisBlock, blockA, blockB}}
		mockedRepository.On("ReplaceChain", blockchain).Return(nil)

		// perform test
		err := miningService.ReplaceChain(blockchain)

		// test verification
		assert.Equal(err, ErrInvalidChain)
		mockedRepository.AssertNotCalled(t, "ReplaceChain")
	})
}
