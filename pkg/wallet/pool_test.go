package wallet

import (
	"encoding/hex"
	"testing"

	"github.com/knd/kndchain/pkg/listing"

	"github.com/knd/kndchain/pkg/hashing"

	"github.com/knd/kndchain/pkg/crypto"
	"github.com/stretchr/testify/assert"
)

func TestTransactionPool(t *testing.T) {
	assert := assert.New(t)
	var transactionPool TransactionPool
	var txA, txB, txC Transaction
	secp256k1 := crypto.NewSecp256k1Generator()
	walletA := NewWallet(secp256k1)
	walletB := NewWallet(secp256k1)
	walletC := NewWallet(secp256k1)
	var validTransactions []Transaction
	mockedListing := new(MockedListing)

	beforeEach := func() {
		transactionPool = NewTransactionPool(mockedListing)
		txA = NewTransaction(walletA, walletB.PubKeyHex(), 100)
		txB = NewTransaction(walletB, walletC.PubKeyHex(), 1)
		txC = NewTransaction(walletC, walletA.PubKeyHex(), 99)
		validTransactions = []Transaction{}
	}

	t.Run("adds 1 transaction", func(t *testing.T) {
		beforeEach()

		// perform test
		transactionPool.Add(txA)

		// test verification
		assert.Equal(txA, transactionPool.Get(txA.GetID()))
	})

	t.Run("exists transaction", func(t *testing.T) {
		beforeEach()

		// perform test
		transactionPool.Add(txB)

		// test verification
		assert.True(transactionPool.Exists(walletB.PubKeyHex()))

		receivedTx := transactionPool.GetTransaction(walletB.PubKeyHex())
		assert.Equal(txB, receivedTx)
	})

	t.Run("get valid transactions", func(t *testing.T) {
		beforeEach()

		outputInBytes, _ := hex.DecodeString(hashing.SHA256Hash(txC.GetOutput()))
		invalidSig := walletB.Sign(outputInBytes)
		invalidTx := &Tx{
			ID:     txC.GetID(),
			Output: txC.GetOutput(),
			Input: Input{
				Timestamp: txC.GetInput().Timestamp,
				Amount:    txC.GetInput().Amount,
				Address:   txC.GetInput().Address,
				Signature: hex.EncodeToString(invalidSig),
			},
		}

		transactionPool.Add(txA)
		transactionPool.Add(invalidTx)
		transactionPool.Add(txB)

		// perform test
		validTransactions = transactionPool.ValidTransactions()

		// test verification
		assert.Contains(validTransactions, txA)
		assert.Contains(validTransactions, txB)
	})

	t.Run("clears transaction pool", func(t *testing.T) {
		beforeEach()
		transactionPool.Add(txA)
		transactionPool.Add(txB)

		// perform test
		transactionPool.Clear()

		// test verification
		assert.Empty(transactionPool.All())
	})

	t.Run("clears blockchain transaction", func(t *testing.T) {
		beforeEach()

		bc := &listing.Blockchain{Chain: []listing.Block{
			listing.Block{
				// Data: []string{"dfd"},
			},
		}}
		mockedListing.On("GetBlockchain").Return(bc)

	})
}
