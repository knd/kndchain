package wallet

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/knd/kndchain/pkg/crypto"
	"github.com/knd/kndchain/pkg/hashing"
	"github.com/stretchr/testify/assert"
)

func TestTransaction_HasAutoGeneratedID(t *testing.T) {
	m := new(MockedKeyPairGenerator)
	m.On("Generate").Return([]byte{}, []byte{})

	// perform test
	tx := NewTransaction(NewWallet(m), "receiver", 0)

	// test verification
	assert.NotEmpty(t, tx.GetID())
}

func TestTransaction_OutputHasAmountToReceiver(t *testing.T) {
	m := new(MockedKeyPairGenerator)
	m.On("Generate").Return([]byte{}, []byte{})

	// perform test
	tx := NewTransaction(NewWallet(m), "receiver", 1)

	// test verification
	assert.Equal(t, uint64(1), tx.GetOutput()["receiver"])
}

func TestTransaction_OutputHasRemainingBalanceOfSenderWallet(t *testing.T) {
	m := new(MockedKeyPairGenerator)
	senderPubKey := []byte("pubkey-sender")
	m.On("Generate").Return(senderPubKey, []byte{})
	w := NewWallet(m)

	// perform test
	tx := NewTransaction(w, "receiver", 1)

	// test verification
	assert.Equal(t, uint64(999), tx.GetOutput()[w.PubKeyHex()])
}

func TestTransaction_Input(t *testing.T) {
	assert := assert.New(t)
	secp256k1 := crypto.NewSecp256k1Generator()
	senderWallet := NewWallet(secp256k1)
	receiverWallet := NewWallet(secp256k1)

	// perform test
	tx := NewTransaction(senderWallet, receiverWallet.PubKeyHex(), 99)

	t.Run("has timestamp", func(t *testing.T) {
		assert.NotZero(tx.GetInput().Timestamp)
	})

	t.Run("sets `amount` to the `senderWallet` balance", func(t *testing.T) {
		assert.Equal(InitialBalance, tx.GetInput().Amount)
	})

	t.Run("sets `address` to the `senderWallet` pubKey", func(t *testing.T) {
		assert.Equal(senderWallet.PubKeyHex(), tx.GetInput().Address)
	})

	t.Run("signs the input with senderWallet privKey", func(t *testing.T) {
		ob, _ := hex.DecodeString(hashing.SHA256Hash(tx.GetOutput()))

		assert.True(secp256k1.Verify(senderWallet.PubKey(), ob, tx.GetInput().Signature))
	})
}

func TestIsValidTransaction(t *testing.T) {
	assert := assert.New(t)

	t.Run("returns true if tx is valid", func(t *testing.T) {
		secp256k1 := crypto.NewSecp256k1Generator()
		senderWallet := NewWallet(secp256k1)
		receiverWallet := NewWallet(secp256k1)
		tx := NewTransaction(senderWallet, receiverWallet.PubKeyHex(), 99)

		// perform test
		valid, _ := IsValidTransaction(tx)

		// test verification
		assert.True(valid)
	})

	t.Run("returns false if tx ouptut is invalid", func(t *testing.T) {
		iT := time.Now().Unix()
		senderPubKeyHex := "0x123"
		receiverPubKeyHex := "0x456"
		tx := new(MockedTransaction)
		var s [65]byte
		copy(s[:], []byte("data"))
		tx.On("GetInput").Return(TxInput{
			Timestamp: iT,
			Amount:    1000,
			Address:   senderPubKeyHex,
			Signature: s[:],
		})
		tx.On("GetOutput").Return(TxOutput{
			senderPubKeyHex:   991,
			receiverPubKeyHex: 10,
		})

		// perform test
		valid, err := IsValidTransaction(tx)

		// test verification
		assert.False(valid)
		assert.Equal(ErrInvalidOutputTotalBalance, err)
	})

	t.Run("returns false if tx input signature invalid", func(t *testing.T) {
		iT := time.Now().Unix()
		secp256k1 := crypto.NewSecp256k1Generator()
		senderWallet := NewWallet(secp256k1)
		receiverWallet := NewWallet(secp256k1)
		tx := new(MockedTransaction)
		var s [65]byte
		copy(s[:], []byte("data"))
		tx.On("GetInput").Return(TxInput{
			Timestamp: iT,
			Amount:    1000,
			Address:   senderWallet.PubKeyHex(),
			Signature: s[:],
		})
		tx.On("GetOutput").Return(TxOutput{
			senderWallet.PubKeyHex():   uint64(990),
			receiverWallet.PubKeyHex(): uint64(10),
		})

		// perform test
		valid, err := IsValidTransaction(tx)

		// test verification
		assert.False(valid)
		assert.Equal(ErrInvalidSignature, err)
	})

	t.Run("returns false if tx input signature is signed by different key", func(t *testing.T) {
		iT := time.Now().Unix()
		secp256k1 := crypto.NewSecp256k1Generator()
		senderWallet := NewWallet(secp256k1)
		receiverWallet := NewWallet(secp256k1)
		tx := new(MockedTransaction)
		output := TxOutput{
			senderWallet.PubKeyHex():   uint64(990),
			receiverWallet.PubKeyHex(): uint64(10),
		}
		tx.On("GetOutput").Return(output)
		outputBytes, _ := hex.DecodeString(hashing.SHA256Hash(output))
		tx.On("GetInput").Return(TxInput{
			Timestamp: iT,
			Amount:    1000,
			Address:   senderWallet.PubKeyHex(),
			Signature: receiverWallet.Sign(outputBytes),
		})

		// perform test
		valid, err := IsValidTransaction(tx)

		// test verification
		assert.False(valid)
		assert.Equal(ErrInvalidSignature, err)
	})
}
