package mining

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/knd/kndchain/pkg/hashing"
	"github.com/stretchr/testify/assert"
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
}
