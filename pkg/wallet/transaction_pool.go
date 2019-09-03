package wallet

// TransactionPool provides access to tx pool operations
type TransactionPool interface {
	Get(id string) Transaction
	Add(tx Transaction) error
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

func (p *transactionPool) Get(id string) Transaction {
	if val, ok := p.transactions[id]; ok {
		return val
	}

	return nil
}

func (p *transactionPool) Add(tx Transaction) error {
	p.transactions[tx.GetID()] = tx
	return nil
}
