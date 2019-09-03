package wallet

import (
	"errors"
)

// TransactionPool provides access to tx pool operations
type TransactionPool interface {
	All() map[string]Transaction
	Get(id string) Transaction
	GetTransaction(inputAddress string) Transaction
	Add(tx Transaction) error
	Exists(inputAddress string) bool
	SetPool(newPool map[string]Transaction) error
}

type transactionPool struct {
	transactions map[string]Transaction
}

// NewTransactionPool creates an new transaction pool
func NewTransactionPool() TransactionPool {
	return &transactionPool{
		transactions: make(map[string]Transaction),
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
