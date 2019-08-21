package creating

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_CreateDefautlGenesisBlock(t *testing.T) {
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

func TestService_CreateBlockWithGivenInput(t *testing.T) {
	creatingService := NewService()
	lastBlockHash := "0x789"
	data := []string{"tx1", "tx2"}

	// perform test
	block, err := creatingService.CreateBlock(&lastBlockHash, data)

	// test verification
	assert.Nil(t, err)
	assert.NotEmpty(t, block.Timestamp)
	assert.Equal(t, lastBlockHash, *block.LastHash)
	assert.NotEmpty(t, *block.Hash)
	assert.Equal(t, data, block.Data)
}
