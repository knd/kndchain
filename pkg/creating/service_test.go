package creating

import (
	"enconding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_CreateDefautlGenesisBlock(t *testing.T) {
	addingService := NewService()
	jsonData = "{}"

	var genesisConfig GenesisConfig
	if err := json.Unmarshal([]byte(jsonData), genesisConfig); err != nil {
		t.FailNow()
	}

	// perform test
	genesisBlock := addingService.CreateGenesisBlock(genesisConfig)

	// test verification
	assert.Equal(1, 1)
}
