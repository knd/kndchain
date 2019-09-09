package calculating

import (
	"testing"
	"time"

	"github.com/knd/kndchain/pkg/config"

	"github.com/stretchr/testify/assert"
)

func TestService_Balance(t *testing.T) {
	assert := assert.New(t)
	address := "04c1bc492c403e1484c81316c7ac789353beb57e620a4c15536fcc668830b79dbcdca2a6cf4e01a2be88f9e617016d06c89f8a45a9e1550b29f6d182b9308113fa"
	var service Service

	beforeEach := func() {
		service = NewService()
	}

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

	t.Run("equals initial balance when no outputs for wallet", func(t *testing.T) {
		beforeEach()
		blockchain := &Blockchain{Chain: []Block{}}

		// perform test
		receivedBalance := service.Balance(address, blockchain)

		// test verification
		assert.Equal(config.InitialBalance, receivedBalance)
	})

	t.Run("updates wallet balance when there are outputs for wallet", func(t *testing.T) {
		beforeEach()
		blockchain := &Blockchain{Chain: []Block{
			Block{
				Timestamp:  time.Now(),
				LastHash:   nil,
				Hash:       nil,
				Nonce:      1,
				Difficulty: 1,
				Data: []Transaction{
					createTransaction("75b3d287-386d-4633-bea6-681b226dcbe5", map[string]uint64{"04c1bc492c403e1484c81316c7ac789353beb57e620a4c15536fcc668830b79dbcdca2a6cf4e01a2be88f9e617016d06c89f8a45a9e1550b29f6d182b9308113fa": 810, "0x893": 100, "0x89333": 90}, 1567756159, 1000, "04c1bc492c403e1484c81316c7ac789353beb57e620a4c15536fcc668830b79dbcdca2a6cf4e01a2be88f9e617016d06c89f8a45a9e1550b29f6d182b9308113fa", "027b67184af6964e3b1605c29e854cce41e1c1dfbbbfb4c0a8b3a271f3f9723f5267e7e70695a0b6ff547883d44b4e992d5c46453f92ddc8b9028185bf002dec01"),
				},
			},
		}}

		// perform test & verification
		assert.Equal(1100, int(service.Balance("0x893", blockchain)))
		assert.Equal(1090, int(service.Balance("0x89333", blockchain)))
		assert.Equal(810, int(service.Balance("04c1bc492c403e1484c81316c7ac789353beb57e620a4c15536fcc668830b79dbcdca2a6cf4e01a2be88f9e617016d06c89f8a45a9e1550b29f6d182b9308113fa", blockchain)))
	})

	t.Run("updates wallet balance when there are outputs after recent transaction from the wallet", func(t *testing.T) {
		beforeEach()
		blockchain := &Blockchain{Chain: []Block{
			Block{
				Timestamp:  time.Now(),
				LastHash:   nil,
				Hash:       nil,
				Nonce:      1,
				Difficulty: 1,
				Data: []Transaction{
					createTransaction("75b3d287-386d-4633-bea6-681b226dcbe5", map[string]uint64{"04c1bc492c403e1484c81316c7ac789353beb57e620a4c15536fcc668830b79dbcdca2a6cf4e01a2be88f9e617016d06c89f8a45a9e1550b29f6d182b9308113fa": 810, "0x893": 100, "0x89333": 90}, 1567756159, 1000, "04c1bc492c403e1484c81316c7ac789353beb57e620a4c15536fcc668830b79dbcdca2a6cf4e01a2be88f9e617016d06c89f8a45a9e1550b29f6d182b9308113fa", "027b67184af6964e3b1605c29e854cce41e1c1dfbbbfb4c0a8b3a271f3f9723f5267e7e70695a0b6ff547883d44b4e992d5c46453f92ddc8b9028185bf002dec01"),
					createTransaction("6ff6a803-6500-44f7-89f7-dbcf53b7b701", map[string]uint64{"0444e8eb4de7752fbcbdc28082b63f36b0d372e06952bd6382e3ef3232946e9f44cd641076458acaa2549725b7e41d4f204ef15f3071d1bc2e3b298d00b5a532d1": 999, "0x1233": 1}, 1567756188, 1000, "0444e8eb4de7752fbcbdc28082b63f36b0d372e06952bd6382e3ef3232946e9f44cd641076458acaa2549725b7e41d4f204ef15f3071d1bc2e3b298d00b5a532d1", "fa2266c9c3f3c6c11f08715d2eb32faeae3e64545a1244a077001b0b4cabe3743b0ab01745f8cbe5b5561711724b323dba250b24a17e6a86dd18fd23346858a600"),
					createTransaction("43b0982e-bda0-4726-a686-78b6628b2b19", map[string]uint64{"0444e8eb4de7752fbcbdc28082b63f36b0d372e06952bd6382e3ef3232946e9f44cd641076458acaa2549725b7e41d4f204ef15f3071d1bc2e3b298d00b5a532d1": 5}, 0, 0, "MINER_REWARD", ""),
				},
			},
		}}

		// perform test & verification
		assert.Equal(1004, int(service.Balance("0444e8eb4de7752fbcbdc28082b63f36b0d372e06952bd6382e3ef3232946e9f44cd641076458acaa2549725b7e41d4f204ef15f3071d1bc2e3b298d00b5a532d1", blockchain)))
	})

}
