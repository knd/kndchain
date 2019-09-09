package validating

import (
	"testing"
	"time"

	"github.com/knd/kndchain/pkg/hashing"
	"github.com/stretchr/testify/assert"
)

func TestService_IsInvalidChainWhenGenesisBlockIsInvalid(t *testing.T) {
	validatingService := NewService()
	lastHash := "0x123"
	hash := "0x456"
	blockchain := &Blockchain{
		Chain: []Block{
			Block{
				Timestamp:  time.Now(),
				LastHash:   &lastHash,
				Hash:       &hash,
				Data:       []Transaction{Transaction{}},
				Nonce:      1,
				Difficulty: 1,
			},
		},
	}

	// perform test & verification
	assert.False(t, validatingService.IsValidChain(blockchain))
}

func TestService_IsInvalidChainWhenLastHashIsTampered(t *testing.T) {
	validatingService := NewService()
	genesisTimestamp := time.Now()
	lastHash := "0x123"
	hash := "0x456"
	tamperedLashHash := "tampered"
	blockchain := &Blockchain{
		Chain: []Block{
			Block{
				Timestamp:  genesisTimestamp,
				LastHash:   &lastHash,
				Hash:       &hash,
				Data:       []Transaction{},
				Nonce:      0,
				Difficulty: 1,
			},
			Block{
				Timestamp:  time.Now(),
				LastHash:   &tamperedLashHash,
				Hash:       &hash,
				Data:       []Transaction{},
				Nonce:      1,
				Difficulty: 1,
			},
		},
	}

	// perform test & verification
	assert.False(t, validatingService.IsValidChain(blockchain))
}

func TestService_IsInvalidChainWhenTimestampIsNotInOrder(t *testing.T) {
	validatingService := NewService()

	genesisLastHash := "0x123"
	genesisHash := "0x456"
	genesisTimestamp := time.Now()
	timestamp1 := time.Now().Add(time.Duration(100))
	timestamp2 := time.Now().Add(time.Duration(200))

	genesisBlock := Block{
		Timestamp:  genesisTimestamp,
		LastHash:   &genesisLastHash,
		Hash:       &genesisHash,
		Data:       []Transaction{},
		Nonce:      0,
		Difficulty: 1,
	}

	txA := Transaction{ID: "txA"}
	blockA := Block{
		Timestamp:  timestamp2,
		LastHash:   &genesisHash,
		Data:       []Transaction{txA},
		Nonce:      1,
		Difficulty: 1,
	}
	blockAHash := hashing.SHA256Hash(timestamp2, genesisHash, []Transaction{txA}, blockA.Nonce, blockA.Difficulty)
	blockA.Hash = &blockAHash

	txB := Transaction{ID: "txB"}
	blockB := Block{
		Timestamp:  timestamp1,
		LastHash:   &blockAHash,
		Data:       []Transaction{txB},
		Nonce:      2,
		Difficulty: 1,
	}
	blockBHash := hashing.SHA256Hash(timestamp1, blockAHash, []Transaction{txB}, blockB.Nonce, blockB.Difficulty)
	blockB.Hash = &blockBHash

	blockchain := &Blockchain{Chain: []Block{genesisBlock, blockA, blockB}}

	// perform test & verification
	assert.False(t, validatingService.IsValidChain(blockchain))
}

func TestService_IsValidChainWhenChainContainsOnlyValidBlocks(t *testing.T) {
	validatingService := NewService()

	genesisLastHash := "0x123"
	genesisHash := "0x456"
	genesisTimestamp := time.Now()
	timestamp1 := time.Now().Add(time.Duration(100))
	timestamp2 := time.Now().Add(time.Duration(200))

	genesisBlock := Block{
		Timestamp:  genesisTimestamp,
		LastHash:   &genesisLastHash,
		Hash:       &genesisHash,
		Data:       []Transaction{},
		Nonce:      0,
		Difficulty: 1,
	}

	txA := Transaction{ID: "txA"}
	blockA := Block{
		Timestamp:  timestamp1,
		LastHash:   &genesisHash,
		Data:       []Transaction{txA},
		Nonce:      1,
		Difficulty: 1,
	}
	blockAHash := hashing.SHA256Hash(timestamp1.Unix(), genesisHash, blockA.Data, blockA.Nonce, blockA.Difficulty)
	blockA.Hash = &blockAHash

	txB := Transaction{ID: "txB"}
	blockB := Block{
		Timestamp:  timestamp2,
		LastHash:   &blockAHash,
		Data:       []Transaction{txB},
		Nonce:      2,
		Difficulty: 1,
	}
	blockBHash := hashing.SHA256Hash(timestamp2.Unix(), blockAHash, blockB.Data, blockB.Nonce, blockB.Difficulty)
	blockB.Hash = &blockBHash

	blockchain := &Blockchain{Chain: []Block{genesisBlock, blockA, blockB}}

	// perform test & verification
	assert.True(t, validatingService.IsValidChain(blockchain))
}

func TestService_IsInvalidChainWhenLastBlockJumpsDifficulty(t *testing.T) {
	validatingService := NewService()

	genesisLastHash := "0x123"
	genesisHash := "0x456"
	genesisTimestamp := time.Now()
	timestamp1 := time.Now().Add(time.Duration(100))
	timestamp2 := time.Now().Add(time.Duration(200))

	genesisBlock := Block{
		Timestamp:  genesisTimestamp,
		LastHash:   &genesisLastHash,
		Hash:       &genesisHash,
		Data:       []Transaction{},
		Nonce:      0,
		Difficulty: 5,
	}

	txA := Transaction{ID: "txA"}
	blockA := Block{
		Timestamp:  timestamp1,
		LastHash:   &genesisHash,
		Data:       []Transaction{txA},
		Nonce:      1,
		Difficulty: 4,
	}
	blockAHash := hashing.SHA256Hash(timestamp1, genesisHash, blockA.Data, blockA.Nonce, blockA.Difficulty)
	blockA.Hash = &blockAHash

	txB := Transaction{ID: "txB"}
	blockB := Block{
		Timestamp:  timestamp2,
		LastHash:   &blockAHash,
		Data:       []Transaction{txB},
		Nonce:      2,
		Difficulty: 2,
	}
	blockBHash := hashing.SHA256Hash(timestamp2, blockAHash, blockB.Data, blockB.Nonce, blockB.Difficulty)
	blockB.Hash = &blockBHash

	blockchain := &Blockchain{Chain: []Block{genesisBlock, blockA, blockB}}

	// perform test & verification
	assert.False(t, validatingService.IsValidChain(blockchain))
}

func TestService_ContainsValidTransactions(t *testing.T) {
	assert := assert.New(t)
	var validator Service
	var bc *Blockchain

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

	createBlock := func(timestamp time.Time, lastHash *string, hash *string, data []Transaction, nonce uint32, difficulty uint32) Block {
		return Block{
			Timestamp:  timestamp,
			LastHash:   lastHash,
			Hash:       hash,
			Data:       data,
			Nonce:      nonce,
			Difficulty: difficulty,
		}
	}

	beforeEach := func() {
		validator = NewService()
		bc = &Blockchain{}
	}

	t.Run("returns true if blockchain contains all valid transactions", func(t *testing.T) {
		beforeEach()

		blockTs, _ := time.Parse(time.RFC3339, "2019-09-06T14:18:44.226857+07:00")
		lastHash := "0x000"
		hash := "0x000"
		data := []Transaction{}
		nonce := uint32(0)
		difficulty := uint32(3)

		block := createBlock(blockTs, &lastHash, &hash, data, nonce, difficulty)
		bc.Chain = append(bc.Chain, block)

		blockTs, _ = time.Parse(time.RFC3339, "2019-09-06T14:50:04.265389+07:00")
		lastHash = "0x000"
		hash = "153bdcdd6dcb3d7c4746f91489305275efe324128d235b6d315b6d4118691184"
		nonce = uint32(7)
		difficulty = uint32(2)
		data = []Transaction{
			createTransaction("75b3d287-386d-4633-bea6-681b226dcbe5", map[string]uint64{"04c1bc492c403e1484c81316c7ac789353beb57e620a4c15536fcc668830b79dbcdca2a6cf4e01a2be88f9e617016d06c89f8a45a9e1550b29f6d182b9308113fa": 810, "0x893": 100, "0x89333": 90}, 1567756159, 1000, "04c1bc492c403e1484c81316c7ac789353beb57e620a4c15536fcc668830b79dbcdca2a6cf4e01a2be88f9e617016d06c89f8a45a9e1550b29f6d182b9308113fa", "027b67184af6964e3b1605c29e854cce41e1c1dfbbbfb4c0a8b3a271f3f9723f5267e7e70695a0b6ff547883d44b4e992d5c46453f92ddc8b9028185bf002dec01"),
			createTransaction("6ff6a803-6500-44f7-89f7-dbcf53b7b701", map[string]uint64{"0444e8eb4de7752fbcbdc28082b63f36b0d372e06952bd6382e3ef3232946e9f44cd641076458acaa2549725b7e41d4f204ef15f3071d1bc2e3b298d00b5a532d1": 999, "0x1233": 1}, 1567756188, 1000, "0444e8eb4de7752fbcbdc28082b63f36b0d372e06952bd6382e3ef3232946e9f44cd641076458acaa2549725b7e41d4f204ef15f3071d1bc2e3b298d00b5a532d1", "fa2266c9c3f3c6c11f08715d2eb32faeae3e64545a1244a077001b0b4cabe3743b0ab01745f8cbe5b5561711724b323dba250b24a17e6a86dd18fd23346858a600"),
			createTransaction("43b0982e-bda0-4726-a686-78b6628b2b19", map[string]uint64{"0444e8eb4de7752fbcbdc28082b63f36b0d372e06952bd6382e3ef3232946e9f44cd641076458acaa2549725b7e41d4f204ef15f3071d1bc2e3b298d00b5a532d1": 5}, 0, 0, "MINER_REWARD", ""),
		}
		block = createBlock(blockTs, &lastHash, &hash, data, nonce, difficulty)
		bc.Chain = append(bc.Chain, block)

		// perform test
		valid, err := validator.ContainsValidTransactions(bc)

		// test verification
		assert.Nil(err)
		assert.True(valid)
	})

	t.Run("returns false if block has more than 1 reward transaction", func(t *testing.T) {
		beforeEach()

		blockTs, _ := time.Parse(time.RFC3339, "2019-09-06T14:18:44.226857+07:00")
		lastHash := "0x000"
		hash := "0x000"
		data := []Transaction{}
		nonce := uint32(0)
		difficulty := uint32(3)

		block := createBlock(blockTs, &lastHash, &hash, data, nonce, difficulty)
		bc.Chain = append(bc.Chain, block)

		blockTs, _ = time.Parse(time.RFC3339, "2019-09-06T14:50:04.265389+07:00")
		lastHash = "0x000"
		hash = "153bdcdd6dcb3d7c4746f91489305275efe324128d235b6d315b6d4118691184"
		nonce = uint32(7)
		difficulty = uint32(2)
		data = []Transaction{
			createTransaction("75b3d287-386d-4633-bea6-681b226dcbe5", map[string]uint64{"04c1bc492c403e1484c81316c7ac789353beb57e620a4c15536fcc668830b79dbcdca2a6cf4e01a2be88f9e617016d06c89f8a45a9e1550b29f6d182b9308113fa": 810, "0x893": 100, "0x89333": 90}, 1567756159, 1000, "04c1bc492c403e1484c81316c7ac789353beb57e620a4c15536fcc668830b79dbcdca2a6cf4e01a2be88f9e617016d06c89f8a45a9e1550b29f6d182b9308113fa", "027b67184af6964e3b1605c29e854cce41e1c1dfbbbfb4c0a8b3a271f3f9723f5267e7e70695a0b6ff547883d44b4e992d5c46453f92ddc8b9028185bf002dec01"),
			createTransaction("6ff6a803-6500-44f7-89f7-dbcf53b7b701", map[string]uint64{"0444e8eb4de7752fbcbdc28082b63f36b0d372e06952bd6382e3ef3232946e9f44cd641076458acaa2549725b7e41d4f204ef15f3071d1bc2e3b298d00b5a532d1": 999, "0x1233": 1}, 1567756188, 1000, "0444e8eb4de7752fbcbdc28082b63f36b0d372e06952bd6382e3ef3232946e9f44cd641076458acaa2549725b7e41d4f204ef15f3071d1bc2e3b298d00b5a532d1", "fa2266c9c3f3c6c11f08715d2eb32faeae3e64545a1244a077001b0b4cabe3743b0ab01745f8cbe5b5561711724b323dba250b24a17e6a86dd18fd23346858a600"),
			createTransaction("43b0982e-bda0-4726-a686-78b6628b2b19", map[string]uint64{"0444e8eb4de7752fbcbdc28082b63f36b0d372e06952bd6382e3ef3232946e9f44cd641076458acaa2549725b7e41d4f204ef15f3071d1bc2e3b298d00b5a532d1": 5}, 0, 0, "MINER_REWARD", ""),
			createTransaction("43b0982e", map[string]uint64{"0444e8eb4de7752fbcbdc28082b63f36b0d372e06952bd6382e3ef3232946e9f44cd641076458acaa2549725b7e41d4f204ef15f3071d1bc2e3b298d00b5a532d1": 5}, 0, 0, "MINER_REWARD", ""), // 2nd reward transaction
		}
		block = createBlock(blockTs, &lastHash, &hash, data, nonce, difficulty)
		bc.Chain = append(bc.Chain, block)

		// perform test & verification
		valid, err := validator.ContainsValidTransactions(bc)

		// test verification
		assert.Equal(ErrMinerRewardExceedsLimit, err)
		assert.False(valid)
	})

	t.Run("returns false if a transaction has malformed output", func(t *testing.T) {
		beforeEach()

		blockTs, _ := time.Parse(time.RFC3339, "2019-09-06T14:18:44.226857+07:00")
		lastHash := "0x000"
		hash := "0x000"
		data := []Transaction{}
		nonce := uint32(0)
		difficulty := uint32(3)

		block := createBlock(blockTs, &lastHash, &hash, data, nonce, difficulty)
		bc.Chain = append(bc.Chain, block)

		blockTs, _ = time.Parse(time.RFC3339, "2019-09-06T14:50:04.265389+07:00")
		lastHash = "0x000"
		hash = "153bdcdd6dcb3d7c4746f91489305275efe324128d235b6d315b6d4118691184"
		nonce = uint32(7)
		difficulty = uint32(2)
		data = []Transaction{
			createTransaction("75b3d287-386d-4633-bea6-681b226dcbe5", map[string]uint64{"04c1bc492c403e1484c81316c7ac789353beb57e620a4c15536fcc668830b79dbcdca2a6cf4e01a2be88f9e617016d06c89f8a45a9e1550b29f6d182b9308113fa": 810, "0x893": 100, "0x89333": 90}, 1567756159, 1000, "04c1bc492c403e1484c81316c7ac789353beb57e620a4c15536fcc668830b79dbcdca2a6cf4e01a2be88f9e617016d06c89f8a45a9e1550b29f6d182b9308113fa", "027b67184af6964e3b1605c29e854cce41e1c1dfbbbfb4c0a8b3a271f3f9723f5267e7e70695a0b6ff547883d44b4e992d5c46453f92ddc8b9028185bf002dec01"),
			createTransaction("6ff6a803-6500-44f7-89f7-dbcf53b7b701", map[string]uint64{"0444e8eb4de7752fbcbdc28082b63f36b0d372e06952bd6382e3ef3232946e9f44cd641076458acaa2549725b7e41d4f204ef15f3071d1bc2e3b298d00b5a532d1": 99999, "0x1233": 1}, 1567756188, 1000, "0444e8eb4de7752fbcbdc28082b63f36b0d372e06952bd6382e3ef3232946e9f44cd641076458acaa2549725b7e41d4f204ef15f3071d1bc2e3b298d00b5a532d1", "fa2266c9c3f3c6c11f08715d2eb32faeae3e64545a1244a077001b0b4cabe3743b0ab01745f8cbe5b5561711724b323dba250b24a17e6a86dd18fd23346858a600"), // transaction with malformed output
			createTransaction("43b0982e-bda0-4726-a686-78b6628b2b19", map[string]uint64{"0444e8eb4de7752fbcbdc28082b63f36b0d372e06952bd6382e3ef3232946e9f44cd641076458acaa2549725b7e41d4f204ef15f3071d1bc2e3b298d00b5a532d1": 5}, 0, 0, "MINER_REWARD", ""),
		}
		block = createBlock(blockTs, &lastHash, &hash, data, nonce, difficulty)
		bc.Chain = append(bc.Chain, block)

		// perform test
		valid, err := validator.ContainsValidTransactions(bc)

		// test verification
		assert.Equal(ErrInvalidMinerRewardAmount, err)
		assert.False(valid)
	})

	t.Run("returns false if a reward transaction has malformed output", func(t *testing.T) {
		beforeEach()

		blockTs, _ := time.Parse(time.RFC3339, "2019-09-06T14:18:44.226857+07:00")
		lastHash := "0x000"
		hash := "0x000"
		data := []Transaction{}
		nonce := uint32(0)
		difficulty := uint32(3)

		block := createBlock(blockTs, &lastHash, &hash, data, nonce, difficulty)
		bc.Chain = append(bc.Chain, block)

		blockTs, _ = time.Parse(time.RFC3339, "2019-09-06T14:50:04.265389+07:00")
		lastHash = "0x000"
		hash = "153bdcdd6dcb3d7c4746f91489305275efe324128d235b6d315b6d4118691184"
		nonce = uint32(7)
		difficulty = uint32(2)
		data = []Transaction{
			createTransaction("75b3d287-386d-4633-bea6-681b226dcbe5", map[string]uint64{"04c1bc492c403e1484c81316c7ac789353beb57e620a4c15536fcc668830b79dbcdca2a6cf4e01a2be88f9e617016d06c89f8a45a9e1550b29f6d182b9308113fa": 810, "0x893": 100, "0x89333": 90}, 1567756159, 1000, "04c1bc492c403e1484c81316c7ac789353beb57e620a4c15536fcc668830b79dbcdca2a6cf4e01a2be88f9e617016d06c89f8a45a9e1550b29f6d182b9308113fa", "027b67184af6964e3b1605c29e854cce41e1c1dfbbbfb4c0a8b3a271f3f9723f5267e7e70695a0b6ff547883d44b4e992d5c46453f92ddc8b9028185bf002dec01"),
			createTransaction("6ff6a803-6500-44f7-89f7-dbcf53b7b701", map[string]uint64{"0444e8eb4de7752fbcbdc28082b63f36b0d372e06952bd6382e3ef3232946e9f44cd641076458acaa2549725b7e41d4f204ef15f3071d1bc2e3b298d00b5a532d1": 999, "0x1233": 1}, 1567756188, 1000, "0444e8eb4de7752fbcbdc28082b63f36b0d372e06952bd6382e3ef3232946e9f44cd641076458acaa2549725b7e41d4f204ef15f3071d1bc2e3b298d00b5a532d1", "fa2266c9c3f3c6c11f08715d2eb32faeae3e64545a1244a077001b0b4cabe3743b0ab01745f8cbe5b5561711724b323dba250b24a17e6a86dd18fd23346858a600"),
			createTransaction("43b0982e-bda0-4726-a686-78b6628b2b19", map[string]uint64{"0x123": 5, "0444e8eb4de7752fbcbdc28082b63f36b0d372e06952bd6382e3ef3232946e9f44cd641076458acaa2549725b7e41d4f204ef15f3071d1bc2e3b298d00b5a532d1": 10}, 0, 0, "MINER_REWARD", ""), // malformed reward transaction output
		}
		block = createBlock(blockTs, &lastHash, &hash, data, nonce, difficulty)
		bc.Chain = append(bc.Chain, block)

		// perform test
		valid, err := validator.ContainsValidTransactions(bc)

		// test verification
		assert.Equal(ErrInvalidMinerRewardAmount, err)
		assert.False(valid)
	})

	t.Run("returns false if a transaction has malformed input", func(t *testing.T) {
		beforeEach()

		blockTs, _ := time.Parse(time.RFC3339, "2019-09-06T14:18:44.226857+07:00")
		lastHash := "0x000"
		hash := "0x000"
		data := []Transaction{}
		nonce := uint32(0)
		difficulty := uint32(3)

		block := createBlock(blockTs, &lastHash, &hash, data, nonce, difficulty)
		bc.Chain = append(bc.Chain, block)

		blockTs, _ = time.Parse(time.RFC3339, "2019-09-06T14:50:04.265389+07:00")
		lastHash = "0x000"
		hash = "153bdcdd6dcb3d7c4746f91489305275efe324128d235b6d315b6d4118691184"
		nonce = uint32(7)
		difficulty = uint32(2)
		data = []Transaction{
			createTransaction("75b3d287-386d-4633-bea6-681b226dcbe5", map[string]uint64{"04c1bc492c403e1484c81316c7ac789353beb57e620a4c15536fcc668830b79dbcdca2a6cf4e01a2be88f9e617016d06c89f8a45a9e1550b29f6d182b9308113fa": 810, "0x893": 100, "0x89333": 90}, 1567756159, 1000, "04c1bc492c403e1484c81316c7ac789353beb57e620a4c15536fcc668830b79dbcdca2a6cf4e01a2be88f9e617016d06c89f8a45a9e1550b29f6d182b9308113fa", "027b67184af6964e3b1605c29e854cce41e1c1dfbbbfb4c0a8b3a271f3f9723f5267e7e70695a0b6ff547883d44b4e992d5c46453f92ddc8b9028185bf002dec01"),
			createTransaction("6ff6a803-6500-44f7-89f7-dbcf53b7b701", map[string]uint64{"0444e8eb4de7752fbcbdc28082b63f36b0d372e06952bd6382e3ef3232946e9f44cd641076458acaa2549725b7e41d4f204ef15f3071d1bc2e3b298d00b5a532d1": 999, "0x1233": 1}, 1567756188, 1000, "0444e8eb4de7752fbcbdc28082b63f36b0d372e06952bd6382e3ef3232946e9f44cd641076458acaa2549725b7e41d4f204ef15f3071d1bc2e3b298d00b5a532d1", "fa2266c9c3f3c6c11f08715d2eb32faeae3e64545a1244a077001b0b4cabe3743b0ab01745f8cbe5b5561711724b323dba250b24a17e6a86dd18fd23346858a600"),
			createTransaction("43b0982e-bda0-4726-a686-78b6628b2b19", map[string]uint64{"0444e8eb4de7752fbcbdc28082b63f36b0d372e06952bd6382e3ef3232946e9f44cd641076458acaa2549725b7e41d4f204ef15f3071d1bc2e3b298d00b5a532d1": 5}, 0, 0, "MINER_REWARD", ""),
		}
		block = createBlock(blockTs, &lastHash, &hash, data, nonce, difficulty)
		bc.Chain = append(bc.Chain, block)

		// block that contains invalid balance input
		blockTs, _ = time.Parse(time.RFC3339, "2019-09-09T11:43:14.056+07:00")
		lastHash = "153bdcdd6dcb3d7c4746f91489305275efe324128d235b6d315b6d4118691184"
		hash = "2cb58d8cdd8d21288546b86d98c7ffbf9f8d2a458726d943472589105d509547"
		nonce = uint32(1)
		difficulty = uint32(1)
		data = []Transaction{
			createTransaction("42be10af-e50d-4f2e-a8dd-6b245738f695", map[string]uint64{"0439c221797ebc2acee2ab71194a125b35c9de7b5880e0fed65e0d9602cfc3206de51048f17c2b39d65bd677a8d0300b72260bfb361bfd109ad9d4d3f24c9e2b6c": 901, "0x893999": 100}, 1568004167, 1001, "0439c221797ebc2acee2ab71194a125b35c9de7b5880e0fed65e0d9602cfc3206de51048f17c2b39d65bd677a8d0300b72260bfb361bfd109ad9d4d3f24c9e2b6c", "aedb4dd132cb57aac82fa18725953d785556ef5505952fcbab25ad0e1fd414890d98ebdc7349668e3058372ea60ddb3dd6229a9fb73ab5fc6b17899553b864db01"),
			createTransaction("838842c2-c34f-4947-9991-3af6491577d9", map[string]uint64{"0444e8eb4de7752fbcbdc28082b63f36b0d372e06952bd6382e3ef3232946e9f44cd641076458acaa2549725b7e41d4f204ef15f3071d1bc2e3b298d00b5a532d1": 5}, 0, 0, "MINER_REWARD", ""),
		}
		block = createBlock(blockTs, &lastHash, &hash, data, nonce, difficulty)
		bc.Chain = append(bc.Chain, block)

		// perform test & verification
		valid, err := validator.ContainsValidTransactions(bc)

		// test verification
		assert.Equal(ErrInvalidInputBalance, err)
		assert.False(valid)
	})

	/*

		t.Run("returns false if a block contains identical transactions", func(t *testing.T) {

			// perform test & verification
			assert.False(validator.ContainsValidTransactions(bc))
		})
	*/
}
