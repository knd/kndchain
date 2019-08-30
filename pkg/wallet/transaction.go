package wallet

import (
	"encoding/hex"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/knd/kndchain/pkg/crypto"
	"github.com/knd/kndchain/pkg/hashing"
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
	ob, err := hex.DecodeString(hashing.SHA256Hash(t.GetOutput()))
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
	i := tx.GetInput()
	o := tx.GetOutput()

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

	outputBytes, err := hex.DecodeString(hashing.SHA256Hash(tx.GetOutput()))
	if err != nil {
		return false, ErrCannotGetOutputBytes
	}

	if !crypto.NewSecp256k1Generator().Verify(pubKeyInByte, outputBytes, i.Signature) {
		return false, ErrInvalidSignature
	}

	return true, nil
}
