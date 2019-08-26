package mining

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/knd/kndchain/pkg/hashing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateGenesisBlock(t *testing.T) {
	assert := assert.New(t)

	t.Run("creates default genesis block", func(t *testing.T) {
		jsonData := "{}"

		var genesisConfig GenesisConfig
		if err := json.Unmarshal([]byte(jsonData), &genesisConfig); err != nil {
			t.FailNow()
		}

		// perform test
		genesisBlock, err := CreateGenesisBlock(&genesisConfig)

		// test verification
		assert.Nil(err)
		assert.NotEmpty(genesisBlock.Timestamp)
		assert.Equal(*genesisBlock.LastHash, DefaultGenesisLastHash)
		assert.Equal(*genesisBlock.Hash, DefaultGenesisHash)
		assert.Empty(genesisBlock.Data)
	})

	t.Run("creates genesis block with given config", func(t *testing.T) {
		jsonData := `{ "lastHash": "0x123", "hash": "0x456", "data": ["tx1", "tx2"] }`

		var genesisConfig GenesisConfig
		if err := json.Unmarshal([]byte(jsonData), &genesisConfig); err != nil {
			t.FailNow()
		}

		// perform test
		genesisBlock, err := CreateGenesisBlock(&genesisConfig)

		// test verification
		assert.Nil(err)
		assert.NotEmpty(genesisBlock.Timestamp)
		assert.Equal("0x123", *genesisBlock.LastHash)
		assert.Equal("0x456", *genesisBlock.Hash)
		assert.Equal([]string{"tx1", "tx2"}, genesisBlock.Data)
	})
}

func TestAdjustBlockDifficulty(t *testing.T) {
	assert := assert.New(t)
	lastHash, hash := "0x123", "0x456"
	lastBlock := &Block{
		Timestamp:  time.Now(),
		LastHash:   &lastHash,
		Hash:       &hash,
		Data:       []string{"tx1"},
		Nonce:      1,
		Difficulty: 2,
	}

	t.Run("raises block difficulty if mining rate is faster than MINE_RATE milliseconds", func(t *testing.T) {
		blockTimestamp := (*lastBlock).Timestamp.Add(time.Millisecond * 999)

		// perform test
		difficulty := AdjustBlockDifficulty(*lastBlock, blockTimestamp)

		// test verification
		assert.Equal(uint32(3), difficulty)
	})

	t.Run("lowers block difficulty if mining rate is slower than threshold MINE_RATE milliseconds", func(t *testing.T) {
		blockTimestamp := (*lastBlock).Timestamp.Add(time.Millisecond * 1001)

		// perform test
		difficulty := AdjustBlockDifficulty(*lastBlock, blockTimestamp)

		// test verficiation
		assert.Equal(uint32(1), difficulty)
	})

	t.Run("lowers block difficulty if mining rate is equal to threshold MINE_RATE milliseconds", func(t *testing.T) {
		blockTimestamp := (*lastBlock).Timestamp.Add(time.Millisecond * 1000)

		// perform test
		difficulty := AdjustBlockDifficulty(*lastBlock, blockTimestamp)

		// test verficiation
		assert.Equal(uint32(2), difficulty)
	})

	t.Run("has minimum difficulty of 1 no matter what", func(t *testing.T) {
		lastBlock.Difficulty = 0
		blockTimestamp := (*lastBlock).Timestamp.Add(time.Millisecond * 1001)

		// perform test
		difficulty := AdjustBlockDifficulty(*lastBlock, blockTimestamp)

		// test verficiation
		assert.Equal(uint32(1), difficulty)
	})
}

func TestHexStringToBinary(t *testing.T) {
	// https://www.binaryhexconverter.com/hex-to-binary-converter
	assert.Equal(t, "0000000100100011010001011000100100010000101010111100110111101111", HexStringToBinary("0123458910abcdef"))
}

func TestService(t *testing.T) {
	assert := assert.New(t)
	var miningService Service
	var mockedRepository *MockedRepository
	var mockedListing *MockedListing
	var mockedValidating *MockedValidating

	beforeEach := func() {
		mockedRepository = new(MockedRepository)
		mockedListing = new(MockedListing)
		mockedValidating = new(MockedValidating)
		miningService = NewService(mockedRepository, mockedListing, mockedValidating, nil)
	}

	t.Run("mines new block", func(t *testing.T) {
		beforeEach()
		mockedListing.On("GetBlockCount").Return(1)
		lastHash := "0x123"
		hash := "0x456"
		var nonce uint32 = 1
		var difficulty uint32 = 1
		lastBlock := Block{
			Timestamp:  time.Now(),
			LastHash:   &lastHash,
			Hash:       &hash,
			Data:       []string{"tx1"},
			Nonce:      nonce,
			Difficulty: difficulty,
		}
		data := []string{"tx2"}

		// perform test
		newBlock, err := miningService.MineNewBlock(&lastBlock, data)

		// test verification
		assert.Nil(err)
		assert.NotEmpty(newBlock.Timestamp)
		assert.Equal("0x456", *newBlock.LastHash)
		assert.Equal("0", (*newBlock.Hash)[:difficulty])
		assert.Equal(hashing.SHA256Hash(data, newBlock.Timestamp, *lastBlock.Hash, newBlock.Nonce, 2), *newBlock.Hash)
		assert.Equal(data, newBlock.Data)
	})

	t.Run("adds block to chain", func(t *testing.T) {
		beforeEach()
		mockedListing.On("GetBlockCount").Return(1)
		LastHash := "0x123"
		Hash := "0x456"
		minedBlock := &Block{
			Timestamp: time.Now(),
			LastHash:  &LastHash,
			Hash:      &Hash,
			Data:      []string{"tx1"},
		}
		mockedRepository.On("AddBlock", minedBlock).Return(nil)

		// perform test
		miningService.AddBlock(minedBlock)

		// test verification
		mockedRepository.AssertExpectations(t)
	})

	t.Run("replaces with nil chain", func(t *testing.T) {
		beforeEach()

		// perform test
		err := miningService.ReplaceChain(nil)

		// test verification
		assert.Equal(err, ErrInvalidChain)
		mockedRepository.AssertNotCalled(t, "ReplaceChain")
	})

	t.Run("replaces with shorter chain", func(t *testing.T) {
		beforeEach()
		mockedListing.On("GetBlockCount").Return(5)

		genesisLastHash := "0x123"
		genesisHash := "0x456"
		genesisTimestamp := time.Now()
		timestamp1 := time.Now().Add(time.Duration(100))
		timestamp2 := time.Now().Add(time.Duration(200))

		genesisBlock := Block{
			Timestamp: genesisTimestamp,
			LastHash:  &genesisLastHash,
			Hash:      &genesisHash,
			Data:      []string{},
		}

		blockA := Block{
			Timestamp: timestamp1,
			LastHash:  &genesisHash,
			Data:      []string{"txA"},
		}
		blockAHash := hashing.SHA256Hash(timestamp1, genesisHash, blockA.Data)
		blockA.Hash = &blockAHash

		blockB := Block{
			Timestamp: timestamp2,
			LastHash:  &blockAHash,
			Data:      []string{"txB"},
		}
		blockBHash := hashing.SHA256Hash(timestamp2, blockAHash, blockB.Data)
		blockB.Hash = &blockBHash

		blockchain := &Blockchain{Chain: []Block{genesisBlock, blockA, blockB}}

		// perform test
		err := miningService.ReplaceChain(blockchain)

		// test verification
		assert.Equal(err, ErrShorterChain)
		mockedRepository.AssertNotCalled(t, "ReplaceChain")
	})

	t.Run("replaces with longer valid chain", func(t *testing.T) {
		beforeEach()
		mockedListing.On("GetBlockCount").Return(2)
		mockedValidating.On("IsValidChain", mock.Anything).Return(true)

		genesisLastHash := "0x123"
		genesisHash := "0x456"
		genesisTimestamp := time.Now()
		timestamp1 := time.Now().Add(time.Duration(100))
		timestamp2 := time.Now().Add(time.Duration(200))

		genesisBlock := Block{
			Timestamp: genesisTimestamp,
			LastHash:  &genesisLastHash,
			Hash:      &genesisHash,
			Data:      []string{},
		}

		blockA := Block{
			Timestamp: timestamp1,
			LastHash:  &genesisHash,
			Data:      []string{"txA"},
		}
		blockAHash := hashing.SHA256Hash(timestamp1, genesisHash, blockA.Data)
		blockA.Hash = &blockAHash

		blockB := Block{
			Timestamp: timestamp2,
			LastHash:  &blockAHash,
			Data:      []string{"txB"},
		}
		blockBHash := hashing.SHA256Hash(timestamp2, blockAHash, blockB.Data)
		blockB.Hash = &blockBHash

		blockchain := &Blockchain{Chain: []Block{genesisBlock, blockA, blockB}}
		mockedRepository.On("ReplaceChain", blockchain).Return(nil)

		// perform test
		err := miningService.ReplaceChain(blockchain)

		// test verification
		assert.Nil(err)
		mockedRepository.AssertCalled(t, "ReplaceChain", blockchain)
	})

	t.Run("replaces with longer invalid chain", func(t *testing.T) {
		beforeEach()
		mockedListing.On("GetBlockCount").Return(2)
		mockedValidating.On("IsValidChain", mock.Anything).Return(false)

		genesisLastHash := "0x123"
		genesisHash := "0x456"
		genesisTimestamp := time.Now()
		timestamp1 := time.Now().Add(time.Duration(100))
		timestamp2 := time.Now().Add(time.Duration(200))

		genesisBlock := Block{
			Timestamp: genesisTimestamp,
			LastHash:  &genesisLastHash,
			Hash:      &genesisHash,
			Data:      []string{},
		}

		blockA := Block{
			Timestamp: timestamp1,
			LastHash:  &genesisHash,
			Data:      []string{"txA"},
		}
		blockAHash := hashing.SHA256Hash(timestamp1, genesisHash, blockA.Data)
		blockA.Hash = &blockAHash

		blockB := Block{
			Timestamp: timestamp2,
			LastHash:  &blockAHash,
			Data:      []string{"txB", "double spend"},
		}
		blockBHash := hashing.SHA256Hash(timestamp2, blockAHash, blockB.Data)
		blockB.Hash = &blockBHash

		blockchain := &Blockchain{Chain: []Block{genesisBlock, blockA, blockB}}
		mockedRepository.On("ReplaceChain", blockchain).Return(nil)

		// perform test
		err := miningService.ReplaceChain(blockchain)

		// test verification
		assert.Equal(err, ErrInvalidChain)
		mockedRepository.AssertNotCalled(t, "ReplaceChain")
	})
}
