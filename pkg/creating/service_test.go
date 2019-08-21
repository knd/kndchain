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
	genesisBlock, _ := addingService.CreateGenesisBlock(genesisConfig)

	// test verification
	lastHash := *genesisBlock.LastHash
	assert.Equal(t, lastHash, DefaultGenesisLastHash)
}
