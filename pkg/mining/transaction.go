package mining

import (
	"github.com/knd/kndchain/pkg/validating"
)

// Input of transaction
type Input struct {
	Timestamp int64  `json:"timestamp"`
	Amount    uint64 `json:"amount"`
	Address   string `json:"address"`
	Signature string `json:"sig"`
}

// Transaction in data
type Transaction struct {
	ID     string            `json:"id"`
	Input  Input             `json:"input"`
	Output map[string]uint64 `json:"output"`
}

func toValidatingTransactions(data []Transaction) []validating.Transaction {
	var vTxs []validating.Transaction
	for _, transaction := range data {
		vTxs = append(vTxs, validating.Transaction{
			ID:     transaction.ID,
			Output: transaction.Output,
			Input: validating.Input{
				Timestamp: transaction.Input.Timestamp,
				Amount:    transaction.Input.Amount,
				Address:   transaction.Input.Address,
				Signature: transaction.Input.Signature,
			},
		})
	}
	return vTxs
}

func toValidatingChain(newChain *Blockchain) *validating.Blockchain {
	vBlockchain := &validating.Blockchain{}
	for _, block := range newChain.Chain {
		vBlock := &validating.Block{
			Timestamp:  block.Timestamp,
			LastHash:   block.LastHash,
			Hash:       block.Hash,
			Data:       toValidatingTransactions(block.Data),
			Nonce:      block.Nonce,
			Difficulty: block.Difficulty,
		}
		vBlockchain.Chain = append(vBlockchain.Chain, *vBlock)
	}
	return vBlockchain
}
