package wallet

import (
	"testing"

	"github.com/knd/kndchain/pkg/crypto"
	"github.com/stretchr/testify/assert"
)

func TestWallet_PublicKeyIsGenerated(t *testing.T) {
	// perform test
	w := NewWallet(crypto.NewSecp256k1Generator())

	// test verification
	assert.NotEmpty(t, w.PubKeyHex())
}

func TestWallet_InitialBalanceOf1000(t *testing.T) {
	// perform test
	w := NewWallet(crypto.NewSecp256k1Generator())

	// test verification
	assert.Equal(t, InitialBalance, w.Balance())
}

func TestWallet_SigningData(t *testing.T) {
	assert := assert.New(t)

	t.Run("verifies signing is done properly", func(t *testing.T) {
		secp256k1 := crypto.NewSecp256k1Generator()
		w := NewWallet(secp256k1)
		data := []byte("hello world")

		// perform test
		signature := w.Sign(data)

		// test verification
		assert.True(secp256k1.Verify(w.PubKey(), data, signature))
	})

	t.Run("verifies signing is NOT done properly", func(t *testing.T) {
		secp256k1 := crypto.NewSecp256k1Generator()
		w := NewWallet(secp256k1)
		data := []byte("hello world")

		// perform test
		signature := NewWallet(secp256k1).Sign(data)

		// test verification
		assert.False(secp256k1.Verify(w.PubKey(), data, signature))
	})
}
