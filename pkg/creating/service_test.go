package creating

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_CreateDefautlGenesisBlock(t *testing.T) {
	addingService := NewService()
	jsonData := "{}"

	var genesisConfig GenesisConfig
	if err := json.Unmarshal([]byte(jsonData), &genesisConfig); err != nil {
		t.FailNow()
	}

	// perform test
	genesisBlock, err := addingService.CreateGenesisBlock(genesisConfig)

	// test verification
	assert.Nil(t, err)
	assert.NotEmpty(t, genesisBlock.Timestamp)
	assert.Equal(t, *genesisBlock.LastHash, DefaultGenesisLastHash)
	assert.Equal(t, *genesisBlock.Hash, DefaultGenesisHash)
	assert.Empty(t, genesisBlock.Data)
}
