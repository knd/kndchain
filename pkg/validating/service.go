package validating

import (
	"github.com/knd/kndchain/pkg/hashing"
)

// Service provides blockchain validating operations
type Service interface {
	IsValidChain(bc *Blockchain) bool
}

type service struct{}

// NewService creates a validating service with necessary dependencies
func NewService() Service {
	return &service{}
}

// IsValidChain returns true if list of blocks compose valid blockchain
func (s *service) IsValidChain(bc *Blockchain) bool {
	if bc == nil || len(bc.chain) == 0 {
		return false
	}
	if len(bc.chain) == 1 {
		return len(bc.chain[0].Data) == 0
	}

	genesisBlock := bc.chain[0]
	prevTimestamp := genesisBlock.Timestamp
	prevHash := genesisBlock.Hash

	for i := 1; i < len(bc.chain); i++ {
		currBlock := bc.chain[i]

		if prevTimestamp.Equal(currBlock.Timestamp) ||
			prevTimestamp.After(currBlock.Timestamp) {
			return false
		}

		if *prevHash != *currBlock.LastHash {
			return false
		}

		if hashing.SHA256Hash(currBlock.Timestamp, *currBlock.LastHash, currBlock.Data) != *currBlock.Hash {
			return false
		}

		prevTimestamp = currBlock.Timestamp
		prevHash = currBlock.Hash
	}

	return true
}
