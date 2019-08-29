package wallet

import (
	"github.com/google/uuid"
)

// Output of transaction which composes of receiver amount and remaining sender balance
type Output map[string]uint64

// Transaction provides access to transacting info
type Transaction interface {
	GetID() string
	GetOutput() Output
}

type transaction struct {
	ID           string
	senderWallet Wallet
	receiver     string
	amount       uint64
}

// NewTransaction creates a transaction
func NewTransaction(s Wallet, r string, amount uint64) Transaction {
	return &transaction{
		ID:           uuid.New().String(),
		senderWallet: s,
		receiver:     r,
		amount:       amount,
	}
}

func (t *transaction) GetID() string {
	return t.ID
}

func (t *transaction) GetOutput() Output {
	o := make(map[string]uint64)

	o[t.receiver] = t.amount
	o[t.senderWallet.PubKeyHex()] = t.senderWallet.Balance() - t.amount

	return o
}
