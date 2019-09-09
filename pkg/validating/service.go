package validating

import (
	"encoding/hex"
	"errors"
	"math"

	"github.com/knd/kndchain/pkg/calculating"

	"github.com/knd/kndchain/pkg/config"
	"github.com/knd/kndchain/pkg/crypto"

	"github.com/knd/kndchain/pkg/hashing"
)

// Service provides blockchain validating operations
type Service interface {
	IsValidChain(bc *Blockchain) bool
	ContainsValidTransactions(bc *Blockchain) (bool, error)
}

type service struct {
	calculator calculating.Service
}

// NewService creates a validating service with necessary dependencies
func NewService(c calculating.Service) Service {
	return &service{c}
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

// ErrInvalidOutputTotalBalance invalid output total balance compared with input amount
var ErrInvalidOutputTotalBalance = errors.New("Output has invalid total balance")

// ErrInvalidSignature invalid signature
var ErrInvalidSignature = errors.New("Signature is invalid")

// ErrInvalidPubKey invalid public key
var ErrInvalidPubKey = errors.New("Invalid public key")

// ErrCannotGetOutputBytes indicates error obtaining output bytes
var ErrCannotGetOutputBytes = errors.New("Cannot obtain output bytes")

// IsValidTransaction returns true if transaction itself contains
// valid input and output information
func IsValidTransaction(tx Transaction) (bool, error) {
	i := tx.Input
	o := tx.Output

	var oBalance uint64
	for _, oAmount := range o {
		oBalance += oAmount
	}

	if i.Amount != oBalance {
		return false, ErrInvalidOutputTotalBalance
	}

	pubKeyInByte, err := hex.DecodeString(i.Address)
	if err != nil {
		return false, ErrInvalidPubKey
	}

	outputBytes, err := hex.DecodeString(hashing.SHA256Hash(tx.Output))
	if err != nil {
		return false, ErrCannotGetOutputBytes
	}

	sigBytes, _ := hex.DecodeString(i.Signature)
	if !crypto.NewSecp256k1Generator().Verify(pubKeyInByte, outputBytes, sigBytes) {
		return false, ErrInvalidSignature
	}

	return true, nil
}

// ErrMinerRewardExceedsLimit indicates when miner reward is more than 1
var ErrMinerRewardExceedsLimit = errors.New("Miner reward exceeds limit")

// ErrInvalidMinerRewardAmount indicates when miner reward tx amount is not same as config
var ErrInvalidMinerRewardAmount = errors.New("Miner reward amount is invalid")

// ErrInvalidInputBalance indicates when the sender has invalid input balance
var ErrInvalidInputBalance = errors.New("Invalid input balance")

// ContainsValidTransactions returns true if all chain transactions are valid
func (s *service) ContainsValidTransactions(bc *Blockchain) (bool, error) {
	if bc == nil {
		return false, errors.New("Empty blockchain")
	}

	cBlockchain := toCalculatingBlockchain(bc)
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
			} else {
				if valid, err := IsValidTransaction(transaction); !valid && err != nil {
					return valid, ErrInvalidMinerRewardAmount
				}

				senderBalance := s.calculator.Balance(transaction.Input.Address, cBlockchain)
				if transaction.Input.Amount != senderBalance {
					return false, ErrInvalidInputBalance
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

func toCalculatingBlockchain(bc *Blockchain) *calculating.Blockchain {
	if bc == nil {
		return nil
	}

	result := &calculating.Blockchain{}
	for _, block := range bc.Chain {
		cTransactions := []calculating.Transaction{}
		for _, transaction := range block.Data {
			cTx := calculating.Transaction{
				ID:     transaction.ID,
				Output: transaction.Output,
				Input: calculating.Input{
					Timestamp: transaction.Input.Timestamp,
					Amount:    transaction.Input.Amount,
					Address:   transaction.Input.Address,
					Signature: transaction.Input.Signature,
				},
			}
			cTransactions = append(cTransactions, cTx)
		}
		cBlock := calculating.Block{
			Timestamp:  block.Timestamp,
			LastHash:   block.LastHash,
			Hash:       block.Hash,
			Data:       cTransactions,
			Nonce:      block.Nonce,
			Difficulty: block.Difficulty,
		}
		result.Chain = append(result.Chain, cBlock)
	}

	return result
}
