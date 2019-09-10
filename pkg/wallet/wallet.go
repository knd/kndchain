package wallet

import (
	"encoding/hex"
	"errors"
	"log"

	"github.com/knd/kndchain/pkg/calculating"
	"github.com/knd/kndchain/pkg/listing"
)

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
	CreateTransaction(receiver string, amount uint64, lister listing.Service) (Transaction, error)
}

type wallet struct {
	gen        KeyPairGenerator
	balance    uint64
	publicKey  []byte
	privateKey []byte
	calculator calculating.Service
}

// NewWallet creates a wallet with necessary dependencies
func NewWallet(kpg KeyPairGenerator, c calculating.Service, initialBalance uint64) Wallet {
	pubKey, privKey := kpg.Generate()

	w := &wallet{
		gen:        kpg,
		balance:    initialBalance,
		publicKey:  pubKey,
		privateKey: privKey,
		calculator: c,
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
func (w *wallet) CreateTransaction(receiver string, amount uint64, lister listing.Service) (Transaction, error) {
	bc := toCalculatingBlockchain(lister.GetBlockchain())
	if bc != nil {
		w.balance = w.calculator.Balance(w.PubKeyHex(), bc)
	}

	if amount > w.Balance() {
		return nil, ErrTxAmountExceedsBalance
	}

	return NewTransaction(w, receiver, amount), nil
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
