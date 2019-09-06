package validating

import (
	"math"

	"github.com/knd/kndchain/pkg/hashing"
)

// Service provides blockchain validating operations
type Service interface {
	IsValidChain(bc *Blockchain) bool
	ContainsValidTransactions(bc *Blockchain) bool
}

type service struct{}

// NewService creates a validating service with necessary dependencies
func NewService() Service {
	return &service{}
}

// IsValidChain returns true if list of blocks compose valid blockchain
func (s *service) IsValidChain(bc *Blockchain) bool {
	if bc == nil || len(bc.Chain) == 0 {
		return false
	}
	if len(bc.Chain) == 1 {
		// the only constrant for valid genesis block is that data is empty
		return len(bc.Chain[0].Data) == 0
	}

	genesisBlock := bc.Chain[0]
	prevTimestamp := genesisBlock.Timestamp
	prevHash := genesisBlock.Hash
	prevBlockDifficulty := genesisBlock.Difficulty

	for i := 1; i < len(bc.Chain); i++ {
		currBlock := bc.Chain[i]

		if prevTimestamp.Equal(currBlock.Timestamp) ||
			prevTimestamp.After(currBlock.Timestamp) {
			return false
		}

		if *prevHash != *currBlock.LastHash {
			return false
		}

		// Prevent difficulty jump
		if math.Abs(float64(prevBlockDifficulty-currBlock.Difficulty)) > 1 {
			return false
		}

		if hashing.SHA256Hash(currBlock.Timestamp.Unix(), *currBlock.LastHash, currBlock.Data, currBlock.Nonce, currBlock.Difficulty) != *currBlock.Hash {
			return false
		}

		prevTimestamp = currBlock.Timestamp
		prevHash = currBlock.Hash
		prevBlockDifficulty = currBlock.Difficulty
	}

	return true
}

// ContainsValidTransactions returns true if all chain transactions are valid
func (s *service) ContainsValidTransactions(bc *Blockchain) bool {
	// TODO: Implement this
	return false
}
