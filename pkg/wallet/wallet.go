package wallet

import (
	"encoding/hex"
	"errors"
	"log"
)

// InitialBalance is the balance when wallet is created
const InitialBalance uint64 = 1000

// ErrTxAmountExceedsBalance indicates tx amount exceeds current balance
var ErrTxAmountExceedsBalance = errors.New("Tx amount exceeds balance")

// KeyPairGenerator provides access to key pair generating operations
type KeyPairGenerator interface {
	Generate() (pubKey, privKey []byte)
	Sign(msg, privKey []byte) ([]byte, error)
	Verify(pubKey, msg, signature []byte) bool
}

// Wallet provides access to wallet operations
type Wallet interface {
	PubKey() []byte
	PubKeyHex() string
	Balance() uint64
	Sign(data []byte) []byte
	CreateTransaction(receiver string, amount uint64) (Transaction, error)
}

type wallet struct {
	gen        KeyPairGenerator
	balance    uint64
	publicKey  []byte
	privateKey []byte
}

// NewWallet creates a wallet with necessary dependencies
func NewWallet(kpg KeyPairGenerator) Wallet {
	pubKey, privKey := kpg.Generate()

	w := &wallet{
		gen:        kpg,
		balance:    InitialBalance,
		publicKey:  pubKey,
		privateKey: privKey,
	}

	return w
}

// PubKey returns public key in bytes
func (w *wallet) PubKey() []byte {
	return w.publicKey
}

// PubKeyHex returns public key in hex string
func (w *wallet) PubKeyHex() string {
	return hex.EncodeToString(w.publicKey)
}

// Balance returns the current balance of wallet
func (w *wallet) Balance() uint64 {
	return w.balance
}

// Sign returns a signed signature of input string
func (w *wallet) Sign(data []byte) []byte {
	b, err := w.gen.Sign(data, w.privateKey)
	if err != nil {
		log.Fatal(err)
	}

	return b
}

// CreateTransaction creates a new transaction from this wallet
func (w *wallet) CreateTransaction(receiver string, amount uint64) (Transaction, error) {
	if amount > w.Balance() {
		return nil, ErrTxAmountExceedsBalance
	}

	return NewTransaction(w, receiver, amount), nil
}
