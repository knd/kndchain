package wallet

import (
	"errors"

	"github.com/knd/kndchain/pkg/validating"

	"github.com/knd/kndchain/pkg/listing"
)

// TransactionPool provides access to tx pool operations
type TransactionPool interface {
	All() map[string]Transaction
	Get(id string) Transaction
	GetTransaction(inputAddress string) Transaction
	Add(tx Transaction) error
	Exists(inputAddress string) bool
	SetPool(newPool map[string]Transaction) error
	ValidTransactions() []Transaction
	Clear() error
	ClearBlockTransactions() error
}

type transactionPool struct {
	transactions map[string]Transaction
	lister       listing.Service
}

// NewTransactionPool creates an new transaction pool
func NewTransactionPool(l listing.Service) TransactionPool {
	return &transactionPool{
		transactions: make(map[string]Transaction),
		lister:       l,
	}
}

func (p *transactionPool) All() map[string]Transaction {
	return p.transactions
}

func (p *transactionPool) Get(id string) Transaction {
	if val, ok := p.transactions[id]; ok {
		return val
	}

	return nil
}

func (p *transactionPool) Add(tx Transaction) error {
	if tx == nil {
		return errors.New("Can't add nil transaction to pool")
	}

	p.transactions[tx.GetID()] = tx

	return nil
}

func (p *transactionPool) Exists(inputAddress string) bool {
	for _, tx := range p.transactions {
		if tx.GetInput().Address == inputAddress {
			return true
		}
	}

	return false
}

func (p *transactionPool) GetTransaction(inputAddress string) Transaction {
	for _, tx := range p.transactions {
		if tx.GetInput().Address == inputAddress {
			return tx
		}
	}

	return nil
}

func (p *transactionPool) SetPool(newPool map[string]Transaction) error {
	p.transactions = newPool
	return nil
}

func (p *transactionPool) ValidTransactions() []Transaction {
	var validTxs []Transaction
	for _, tx := range p.transactions {
		validatingTx := validating.Transaction{
			ID:     tx.GetID(),
			Output: tx.GetOutput(),
			Input: validating.Input{
				Timestamp: tx.GetInput().Timestamp,
				Amount:    tx.GetInput().Amount,
				Address:   tx.GetInput().Address,
				Signature: tx.GetInput().Signature,
			},
		}

		valid, err := validating.IsValidTransaction(validatingTx)
		if valid && err == nil {
			validTxs = append(validTxs, tx)
		}
	}
	return validTxs
}

func (p *transactionPool) Clear() error {
	p.transactions = make(map[string]Transaction)
	return nil
}

func (p *transactionPool) ClearBlockTransactions() error {
	for _, block := range p.lister.GetBlockchain().Chain {
		for _, transaction := range block.Data {
			if _, ok := p.transactions[transaction.ID]; ok {
				delete(p.transactions, transaction.ID)
			}
		}
	}

	return nil
}
