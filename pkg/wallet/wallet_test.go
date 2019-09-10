package wallet

import (
	"testing"

	"github.com/knd/kndchain/pkg/crypto"
	"github.com/stretchr/testify/assert"
)

func TestWallet_PublicKeyIsGenerated(t *testing.T) {
	// perform test
	w := NewWallet(crypto.NewSecp256k1Generator(), new(MockedCalculating), 1000)

	// test verification
	assert.NotEmpty(t, w.PubKeyHex())
}

func TestWallet_InitialBalanceOf1000(t *testing.T) {
	// perform test
	w := NewWallet(crypto.NewSecp256k1Generator(), new(MockedCalculating), 1000)

	// test verification
	assert.Equal(t, 1000, int(w.Balance()))
}

func TestWallet_SigningData(t *testing.T) {
	assert := assert.New(t)

	t.Run("verifies signing is done properly", func(t *testing.T) {
		secp256k1 := crypto.NewSecp256k1Generator()
		w := NewWallet(secp256k1, new(MockedCalculating), 1000)
		data := []byte("hello world")

		// perform test
		signature := w.Sign(data)

		// test verification
		assert.True(secp256k1.Verify(w.PubKey(), data, signature))
	})

	t.Run("verifies signing is NOT done properly", func(t *testing.T) {
		secp256k1 := crypto.NewSecp256k1Generator()
		w := NewWallet(secp256k1, new(MockedCalculating), 1000)
		data := []byte("hello world")

		// perform test
		signature := NewWallet(secp256k1, new(MockedCalculating), 1000).Sign(data)

		// test verification
		assert.False(secp256k1.Verify(w.PubKey(), data, signature))
	})
}

func TestWallet_CreateTransaction(t *testing.T) {
	assert := assert.New(t)
	secp256k1 := crypto.NewSecp256k1Generator()
	senderWallet := NewWallet(secp256k1, new(MockedCalculating), 1000)
	receiverWallet := NewWallet(secp256k1, new(MockedCalculating), 1000)
	mockedLister := new(MockedListing)
	mockedLister.On("GetBlockchain").Return(nil)
	txA, errA := senderWallet.CreateTransaction(receiverWallet.PubKeyHex(), 99, mockedLister)
	txB, errB := senderWallet.CreateTransaction(receiverWallet.PubKeyHex(), 1001, mockedLister)

	t.Run("created transaction with input matched wallet", func(t *testing.T) {
		assert.Nil(errA)
		assert.Equal(senderWallet.PubKeyHex(), txA.GetInput().Address)
		assert.Equal(senderWallet.Balance(), txA.GetInput().Amount)
	})

	t.Run("created transaction with output containing receiver amount", func(t *testing.T) {
		assert.Equal(uint64(99), txA.GetOutput()[receiverWallet.PubKeyHex()])
	})

	t.Run("fails to create transaction with amount exceeding balance", func(t *testing.T) {
		assert.Equal(errB, ErrTxAmountExceedsBalance)
		assert.Nil(txB)
	})
}
