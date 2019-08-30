package wallet

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"

	"github.com/google/uuid"
)

// TxInput of transaction which composes of `timestamp`, `sender current balance amount`, `sender address`, `signature` of transaction output
type TxInput struct {
	Timestamp int64
	Amount    uint64
	Address   string
	Signature []byte
}

// TxOutput of transaction which composes of `receiver amount` and `remaining sender balance`
type TxOutput map[string]uint64

// Transaction provides access to transacting info
type Transaction interface {
	GetID() string
	GetInput() TxInput
	GetOutput() TxOutput
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

func (t *transaction) GetOutput() TxOutput {
	o := make(map[string]uint64)

	o[t.receiver] = t.amount
	o[t.senderWallet.PubKeyHex()] = t.senderWallet.Balance() - t.amount

	return o
}

func (t *transaction) GetInput() TxInput {
	ob, err := GetBytes(t.GetOutput())
	if err != nil {
		log.Fatal(err)
	}

	return TxInput{
		Timestamp: time.Now().Unix(),
		Amount:    t.senderWallet.Balance(),
		Address:   t.senderWallet.PubKeyHex(),
		Signature: t.senderWallet.Sign(ob),
	}
}

// GetBytes returns bytes of any Go inteface
func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
