package mining

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestService_CreateDefaultGenesisBlock(t *testing.T) {
	creatingService := NewService()
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
	creatingService := NewService()
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
	creatingService := NewService()
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
	assert.NotEmpty(t, *newBlock.Hash)
	assert.Equal(t, data, newBlock.Data)
}
