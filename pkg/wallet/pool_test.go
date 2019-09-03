package wallet

import (
	"testing"

	"github.com/knd/kndchain/pkg/crypto"
	"github.com/stretchr/testify/assert"
)

func TestTransactionPool(t *testing.T) {
	assert := assert.New(t)
	var transactionPool TransactionPool
	var txA Transaction
	secp256k1 := crypto.NewSecp256k1Generator()
	walletA := NewWallet(secp256k1)
	walletB := NewWallet(secp256k1)

	beforeEach := func() {
		transactionPool = NewTransactionPool()
		txA = NewTransaction(walletA, walletB.PubKeyHex(), 100)
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

		txB := NewTransaction(walletB, walletB.PubKeyHex(), 1)

		// perform test
		transactionPool.Add(txB)

		// test verification
		assert.True(transactionPool.Exists(walletB.PubKeyHex()))

		receivedTx := transactionPool.GetTransaction(walletB.PubKeyHex())
		assert.Equal(txB, receivedTx)
	})
}
