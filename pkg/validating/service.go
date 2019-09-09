package validating

import (
	"errors"
	"math"

	"github.com/knd/kndchain/pkg/config"

	"github.com/knd/kndchain/pkg/hashing"
)

// Service provides blockchain validating operations
type Service interface {
	IsValidChain(bc *Blockchain) bool
	ContainsValidTransactions(bc *Blockchain) (bool, error)
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

// ErrMinerRewardExceedsLimit indicates when miner reward is more than 1
var ErrMinerRewardExceedsLimit = errors.New("Miner reward exceeds limit")

// ErrInvalidMinerRewardAmount indicates when miner reward tx amount is not same as config
var ErrInvalidMinerRewardAmount = errors.New("Miner reward amount is invalid")

// ContainsValidTransactions returns true if all chain transactions are valid
func (s *service) ContainsValidTransactions(bc *Blockchain) (bool, error) {
	if bc == nil {
		return false, errors.New("Empty blockchain")
	}

	for i := 0; i < len(bc.Chain); i++ {
		block := bc.Chain[i]
		rewardTransactionCount := 0

		for _, transaction := range block.Data {
			if transaction.Input.Address == config.RewardTxInputAddress {
				rewardTransactionCount++

				if rewardTransactionCount > 1 {
					return false, ErrMinerRewardExceedsLimit
				}

				if len(transaction.Output) > 1 || getFirstValueOfMap(transaction.Output) != config.MiningReward {
					return false, ErrInvalidMinerRewardAmount
				}
			}
		}
	}
	return true, nil
}

func getFirstValueOfMap(m map[string]uint64) uint64 {
	for _, val := range m {
		return val
	}

	return 0
}
