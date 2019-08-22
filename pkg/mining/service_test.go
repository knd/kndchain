package mining

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/knd/kndchain/pkg/hashing"
	"github.com/stretchr/testify/assert"
)

func TestService_CreateDefaultGenesisBlock(t *testing.T) {
	creatingService := NewService(new(MockedRepository))
	jsonData := "{}"

	var genesisConfig GenesisConfig
	if err := json.Unmarshal([]byte(jsonData), &genesisConfig); err != nil {
		t.FailNow()
	}

	// perform test
	genesisBlock, err := creatingService.CreateGenesisBlock(genesisConfig)

	// test verification
	assert.Nil(t, err)
	assert.NotEmpty(t, genesisBlock.Timestamp)
	assert.Equal(t, *genesisBlock.LastHash, DefaultGenesisLastHash)
	assert.Equal(t, *genesisBlock.Hash, DefaultGenesisHash)
	assert.Empty(t, genesisBlock.Data)
}

func TestService_CreateGenesisBlockWithGivenInput(t *testing.T) {
	creatingService := NewService(new(MockedRepository))
	jsonData := `{ "lastHash": "0x123", "hash": "0x456", "data": ["tx1", "tx2"] }`

	var genesisConfig GenesisConfig
	if err := json.Unmarshal([]byte(jsonData), &genesisConfig); err != nil {
		t.FailNow()
	}

	// perform test
	genesisBlock, err := creatingService.CreateGenesisBlock(genesisConfig)

	// test verification
	assert.Nil(t, err)
	assert.NotEmpty(t, genesisBlock.Timestamp)
	assert.Equal(t, "0x123", *genesisBlock.LastHash)
	assert.Equal(t, "0x456", *genesisBlock.Hash)
	assert.Equal(t, []string{"tx1", "tx2"}, genesisBlock.Data)
}

func TestService_MineNewBlock(t *testing.T) {
	creatingService := NewService(new(MockedRepository))
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
	newBlock, err := creatingService.MineNewBlock(&lastBlock, data)

	// test verification
	assert.Nil(t, err)
	assert.NotEmpty(t, newBlock.Timestamp)
	assert.Equal(t, "0x456", *newBlock.LastHash)
	assert.Equal(t, *newBlock.Hash, hashing.SHA256Hash(data, newBlock.Timestamp, *lastBlock.Hash))
	assert.Equal(t, data, newBlock.Data)
}

func TestService_AddBlockToBlockchain(t *testing.T) {
	mockedRepository := new(MockedRepository)
	creatingService := NewService(mockedRepository)
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
	creatingService.AddBlock(minedBlock)

	// test verification
	mockedRepository.AssertExpectations(t)
}
