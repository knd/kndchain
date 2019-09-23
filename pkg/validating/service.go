package validating

import (
	"encoding/hex"
	"errors"
	"log"

	"github.com/knd/kndchain/pkg/calculating"
	"github.com/knd/kndchain/pkg/crypto"
	"github.com/knd/kndchain/pkg/hashing"
	"github.com/knd/kndchain/pkg/listing"
)

// Service provides blockchain validating operations
type Service interface {
	IsValidChain(bc *Blockchain) bool
	ContainsValidTransactions(bc *Blockchain) (bool, error)
}

type service struct {
	lister               listing.Service
	calculator           calculating.Service
	RewardTxInputAddress string
	MiningReward         uint64
}

// NewService creates a validating service with necessary dependencies
func NewService(l listing.Service, c calculating.Service, rewardInputAddress string, reward uint64) Service {
	return &service{l, c, rewardInputAddress, reward}
}

// IsValidChain returns true if list of blocks compose valid blockchain
func (s *service) IsValidChain(bc *Blockchain) bool {
	if bc == nil || len(bc.Chain) == 0 {
		log.Println("Not a valid chain. Chain length is nil or zero length")
		return false
	}
	if len(bc.Chain) == 1 {
		// the only constrant for valid genesis block is that data is empty
		log.Println("Not a valid chain. Genesis block should have zero data")
		return len(bc.Chain[0].Data) == 0
	}

	genesisBlock := bc.Chain[0]
	prevTimestamp := genesisBlock.Timestamp
	prevHash := genesisBlock.Hash
	prevBlockDifficulty := genesisBlock.Difficulty

	for i := 1; i < len(bc.Chain); i++ {
		currBlock := bc.Chain[i]
		log.Printf("PrevBlockHash=%s, CurrBlockHash=%s", *prevHash, *currBlock.Hash)

		if prevTimestamp >= currBlock.Timestamp {
			log.Println("Not a valid chain. Block timestamp is not chronological")
			return false
		}

		if *prevHash != *currBlock.LastHash {
			log.Printf("Not a valid chain. Last block hash is not inside current block's last hash. *prevHash=%s, *currBlock.LastHash=%s", *prevHash, *currBlock.LastHash)
			return false
		}

		// Prevent difficulty jump
		if (prevBlockDifficulty > currBlock.Difficulty && prevBlockDifficulty-currBlock.Difficulty > 1) || (currBlock.Difficulty > prevBlockDifficulty && currBlock.Difficulty-prevBlockDifficulty > 1) {
			log.Printf("Not a valid chain. Difficulty jump in blocks. prevBlockDifficulty=%d, currBlock.Difficulty=%d", prevBlockDifficulty, currBlock.Difficulty)
			return false
		}

		if hashing.SHA256Hash(currBlock.Timestamp, *currBlock.LastHash, currBlock.Data, currBlock.Nonce, currBlock.Difficulty) != *currBlock.Hash {
			log.Println("Not a valid chain. Current block hash is not correct SHA256")
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

// ErrDuplicateTransaction indicates when the sender has duplicate transactions in same block
var ErrDuplicateTransaction = errors.New("Duplicate transaction in same block")

// ContainsValidTransactions returns true if all chain transactions are valid
func (s *service) ContainsValidTransactions(bc *Blockchain) (bool, error) {
	if bc == nil {
		return false, errors.New("Empty blockchain")
	}

	cBlockchain := toCalculatingBlockchain(s.lister.GetBlockchain())
	for i := 0; i < len(bc.Chain); i++ {
		block := bc.Chain[i]
		rewardTransactionCount := 0
		senderTransactions := map[string]bool{}

		for _, transaction := range block.Data {
			if transaction.Input.Address == s.RewardTxInputAddress {
				rewardTransactionCount++

				if rewardTransactionCount > 1 {
					return false, ErrMinerRewardExceedsLimit
				}

				if len(transaction.Output) > 1 || getFirstValueOfMap(transaction.Output) != s.MiningReward {
					return false, ErrInvalidMinerRewardAmount
				}
			} else {
				if valid, err := IsValidTransaction(transaction); !valid && err != nil {
					return valid, ErrInvalidMinerRewardAmount
				}

				senderBalance := s.calculator.BalanceByBlockIndex(transaction.Input.Address, cBlockchain, i-1)
				if transaction.Input.Amount != senderBalance {
					return false, ErrInvalidInputBalance
				}

				if _, present := senderTransactions[transaction.Input.Address]; present {
					return false, ErrDuplicateTransaction
				}
				senderTransactions[transaction.Input.Address] = true
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

func toCalculatingBlockchain(bc *listing.Blockchain) *calculating.Blockchain {
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
