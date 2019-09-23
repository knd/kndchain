package validating

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsValidTransaction(t *testing.T) {
	assert := assert.New(t)

	createTransaction := func(id string, output map[string]uint64, timestamp int64, amount uint64, address string, signature string) Transaction {
		return Transaction{
			ID: id,
			Input: Input{
				Timestamp: timestamp,
				Amount:    amount,
				Address:   address,
				Signature: signature,
			},
			Output: output,
		}
	}

	t.Run("returns true if tx is valid", func(t *testing.T) {
		tx := createTransaction("75b3d287-386d-4633-bea6-681b226dcbe5", map[string]uint64{"04c1bc492c403e1484c81316c7ac789353beb57e620a4c15536fcc668830b79dbcdca2a6cf4e01a2be88f9e617016d06c89f8a45a9e1550b29f6d182b9308113fa": 810, "0x893": 100, "0x89333": 90}, 1567756159, 1000, "04c1bc492c403e1484c81316c7ac789353beb57e620a4c15536fcc668830b79dbcdca2a6cf4e01a2be88f9e617016d06c89f8a45a9e1550b29f6d182b9308113fa", "027b67184af6964e3b1605c29e854cce41e1c1dfbbbfb4c0a8b3a271f3f9723f5267e7e70695a0b6ff547883d44b4e992d5c46453f92ddc8b9028185bf002dec01")

		// perform test
		valid, _ := IsValidTransaction(tx)

		// test verification
		assert.True(valid)
	})

	t.Run("returns false if tx ouptut is invalid", func(t *testing.T) {
		iT := time.Now().UnixNano()
		senderPubKeyHex := "0x123"
		receiverPubKeyHex := "0x456"
		tx := Transaction{}
		var s [65]byte
		copy(s[:], []byte("data"))
		tx.Input = Input{
			Timestamp: iT,
			Amount:    1000,
			Address:   senderPubKeyHex,
			Signature: hex.EncodeToString(s[:]),
		}
		tx.Output = map[string]uint64{
			senderPubKeyHex:   991,
			receiverPubKeyHex: 10,
		}

		// perform test
		valid, err := IsValidTransaction(tx)

		// test verification
		assert.False(valid)
		assert.Equal(ErrInvalidOutputTotalBalance, err)
	})

	t.Run("returns false if tx input signature invalid", func(t *testing.T) {
		tx := createTransaction("75b3d287-386d-4633-bea6-681b226dcbe5", map[string]uint64{"04c1bc492c403e1484c81316c7ac789353beb57e620a4c15536fcc668830b79dbcdca2a6cf4e01a2be88f9e617016d06c89f8a45a9e1550b29f6d182b9308113fa": 810, "0x893": 100, "0x89333": 90}, 1567756159, 1000, "04c1bc492c403e1484c81316c7ac789353beb57e620a4c15536fcc668830b79dbcdca2a6cf4e01a2be88f9e617016d06c89f8a45a9e1550b29f6d182b9308113fa", "027b67184af6964e3b1605c29e854cce41e1c1dfbbbfb4c0a8b3a271f3f9723f5267e7e70695a0b6ff547883d44b4e992d5c46453f92ddc8b9028185bf002dec01")
		tx.Input.Signature = "abc"

		// perform test
		valid, err := IsValidTransaction(tx)

		// test verification
		assert.False(valid)
		assert.Equal(ErrInvalidSignature, err)
	})

	t.Run("returns false if tx input signature is signed by different key", func(t *testing.T) {
		tx := createTransaction("75b3d287-386d-4633-bea6-681b226dcbe5", map[string]uint64{"04c1bc492c403e1484c81316c7ac789353beb57e620a4c15536fcc668830b79dbcdca2a6cf4e01a2be88f9e617016d06c89f8a45a9e1550b29f6d182b9308113fa": 810, "0x893": 100, "0x89333": 90}, 1567756159, 1000, "04c1bc492c403e1484c81316c7ac789353beb57e620a4c15536fcc668830b79dbcdca2a6cf4e01a2be88f9e617016d06c89f8a45a9e1550b29f6d182b9308113fa", "fa2266c9c3f3c6c11f08715d2eb32faeae3e64545a1244a077001b0b4cabe3743b0ab01745f8cbe5b5561711724b323dba250b24a17e6a86dd18fd23346858a600")

		// perform test
		valid, err := IsValidTransaction(tx)

		// test verification
		assert.False(valid)
		assert.Equal(ErrInvalidSignature, err)
	})
}
