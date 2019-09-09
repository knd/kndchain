package wallet

import (
	"encoding/hex"
	"errors"
	"log"
	"time"

	"github.com/knd/kndchain/pkg/config"

	"github.com/google/uuid"
	"github.com/knd/kndchain/pkg/crypto"
	"github.com/knd/kndchain/pkg/hashing"
)

// Input of transaction
type Input struct {
	Timestamp int64  `json:"timestamp"`
	Amount    uint64 `json:"amount"`
	Address   string `json:"address"`
	Signature string `json:"sig"`
}

// Output of transaction
type Output map[string]uint64

// Transaction provides access to transacting info
type Transaction interface {
	GetID() string
	GetInput() Input
	GetOutput() Output
	Append(w Wallet, r string, amount uint64) error
}

// Tx encapsulates necessary transaction info
type Tx struct {
	ID     string `json:"id"`
	Input  Input  `json:"input"`
	Output Output `json:"output"`
}

// ErrAmountExceedsBalance indicates amount to be sent exceeds the sender remaining balance
var ErrAmountExceedsBalance = errors.New("Amount exceeds sender balance")

// NewTransaction creates a transaction
func NewTransaction(w Wallet, r string, amount uint64) Transaction {
	tx := &Tx{ID: uuid.New().String()}
	tx.Output = tx.generateOutput(w, r, amount)
	tx.Input = tx.generateInput(w, tx.Output)

	return tx
}

// Append adds more amount and receiver
func (t *Tx) Append(w Wallet, receiver string, amount uint64) error {
	if amount > t.Output[w.PubKeyHex()] {
		return ErrAmountExceedsBalance
	}

	if _, ok := t.Output[receiver]; ok {
		t.Output[receiver] += amount
	} else {
		t.Output[receiver] = amount
	}

	t.Output[w.PubKeyHex()] -= amount
	t.Input = t.generateInput(w, t.Output)

	return nil
}

// GetInput returns input
func (t *Tx) GetInput() Input {
	return t.Input
}

// GetOutput returns output
func (t *Tx) GetOutput() Output {
	return t.Output
}

// GetID returns tx ID
func (t *Tx) GetID() string {
	return t.ID
}

func (t *Tx) generateOutput(w Wallet, receiver string, amount uint64) Output {
	o := Output{}
	o[receiver] = amount
	o[w.PubKeyHex()] = w.Balance() - amount
	return o
}

func (t *Tx) generateInput(w Wallet, op Output) Input {
	ob, err := hex.DecodeString(hashing.SHA256Hash(op))
	if err != nil {
		log.Fatal(err)
	}

	return Input{
		Timestamp: time.Now().Unix(),
		Amount:    w.Balance(),
		Address:   w.PubKeyHex(),
		Signature: hex.EncodeToString(w.Sign(ob)),
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

	sigBytes, _ := hex.DecodeString(i.Signature)
	if !crypto.NewSecp256k1Generator().Verify(pubKeyInByte, outputBytes, sigBytes) {
		return false, ErrInvalidSignature
	}

	return true, nil
}

// GetRewardTransactionInput returns the special input in the reward tx to miner
func GetRewardTransactionInput() Input {
	return Input{
		Address: config.RewardTxInputAddress,
	}
}

// CreateRewardTransaction creates a reward transaction to the miner to seals block
func CreateRewardTransaction(mw Wallet) (Transaction, error) {
	tx := &Tx{ID: uuid.New().String()}
	tx.Input = GetRewardTransactionInput()

	o := Output{}
	o[mw.PubKeyHex()] = config.MiningReward
	tx.Output = o

	return tx, nil
}
