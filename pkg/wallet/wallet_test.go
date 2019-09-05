package wallet

import (
	"testing"

	"github.com/knd/kndchain/pkg/listing"

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

func TestWallet_CreateTransaction(t *testing.T) {
	assert := assert.New(t)
	secp256k1 := crypto.NewSecp256k1Generator()
	senderWallet := NewWallet(secp256k1)
	receiverWallet := NewWallet(secp256k1)
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

func TestWallet_CalculateBalance(t *testing.T) {
	assert := assert.New(t)
	secp265k1 := crypto.NewSecp256k1Generator()
	calculatingWallet := NewWallet(secp265k1)
	var mockedLister *MockedListing

	beforeEach := func() {
		mockedLister = new(MockedListing)
	}

	addEachTxToEachBlock := func(txs ...Transaction) []listing.Block {
		var blocks []listing.Block
		for _, tx := range txs {
			blocks = append(blocks, listing.Block{
				Data: []listing.Transaction{
					listing.Transaction{
						ID:     tx.GetID(),
						Output: tx.GetOutput(),
						Input: listing.Input{
							Timestamp: tx.GetInput().Timestamp,
							Amount:    tx.GetInput().Amount,
							Address:   tx.GetInput().Address,
							Signature: tx.GetInput().Signature,
						},
					},
				},
			})
		}
		return blocks
	}

	addTxsToOneBlock := func(txs ...Transaction) listing.Block {
		var block listing.Block
		for _, tx := range txs {
			block.Data = append(block.Data, listing.Transaction{
				ID:     tx.GetID(),
				Output: tx.GetOutput(),
				Input: listing.Input{
					Timestamp: tx.GetInput().Timestamp,
					Amount:    tx.GetInput().Amount,
					Address:   tx.GetInput().Address,
					Signature: tx.GetInput().Signature,
				},
			})

		}
		return block
	}

	t.Run("equals initial balance when no outputs for wallet", func(t *testing.T) {
		beforeEach()
		blockchain := &listing.Blockchain{Chain: []listing.Block{}}
		mockedLister.On("GetBlockchain").Return(blockchain)

		// perform test
		receivedBalance := CalculateBalance(mockedLister, calculatingWallet.PubKeyHex())

		// test verification
		assert.Equal(InitialBalance, receivedBalance)
	})

	t.Run("updates wallet balance when there are outputs for wallet", func(t *testing.T) {
		beforeEach()
		mockedLister.On("GetBlockchain").Return(nil)
		txA, _ := NewWallet(secp265k1).CreateTransaction(calculatingWallet.PubKeyHex(), 50, mockedLister)
		txB, _ := NewWallet(secp265k1).CreateTransaction(calculatingWallet.PubKeyHex(), 10, mockedLister)
		blockchain := &listing.Blockchain{Chain: addEachTxToEachBlock(txA, txB)}
		mockedLister = new(MockedListing)
		mockedLister.On("GetBlockchain").Return(blockchain)

		// perform test
		receivedBlance := CalculateBalance(mockedLister, calculatingWallet.PubKeyHex())

		// test verification
		assert.Equal(1060, int(receivedBlance))
	})

	t.Run("updates wallet balance when there are recent transaction from this wallet", func(t *testing.T) {
		beforeEach()
		mockedLister.On("GetBlockchain").Return(nil)
		txA, _ := NewWallet(secp265k1).CreateTransaction(calculatingWallet.PubKeyHex(), 50, mockedLister)

		mockedLister = new(MockedListing)
		blockchain := &listing.Blockchain{Chain: addEachTxToEachBlock(txA)}
		mockedLister.On("GetBlockchain").Return(blockchain)
		txB, _ := NewWallet(secp265k1).CreateTransaction(calculatingWallet.PubKeyHex(), 10, mockedLister)

		mockedLister = new(MockedListing)
		blockchain = &listing.Blockchain{Chain: addEachTxToEachBlock(txA, txB)}
		mockedLister.On("GetBlockchain").Return(blockchain)
		txC, _ := calculatingWallet.CreateTransaction(NewWallet(secp265k1).PubKeyHex(), 100, mockedLister)

		mockedLister = new(MockedListing)
		blockchain = &listing.Blockchain{Chain: addEachTxToEachBlock(txA, txB, txC)}
		mockedLister.On("GetBlockchain").Return(blockchain)

		// perform test
		receivedBalance := CalculateBalance(mockedLister, calculatingWallet.PubKeyHex())

		// test verification
		assert.Equal(960, int(receivedBalance))
	})

	t.Run("updates wallet balance when there are outputs after recent transaction from the wallet", func(t *testing.T) {
		beforeEach()
		mockedLister.On("GetBlockchain").Return(nil)
		txA, _ := NewWallet(secp265k1).CreateTransaction(calculatingWallet.PubKeyHex(), 50, mockedLister)

		mockedLister = new(MockedListing)
		blockchain := &listing.Blockchain{Chain: addEachTxToEachBlock(txA)}
		mockedLister.On("GetBlockchain").Return(blockchain)
		txB, _ := NewWallet(secp265k1).CreateTransaction(calculatingWallet.PubKeyHex(), 10, mockedLister)

		mockedLister = new(MockedListing)
		blockchain = &listing.Blockchain{Chain: addEachTxToEachBlock(txA, txB)}
		mockedLister.On("GetBlockchain").Return(blockchain)
		txC, _ := calculatingWallet.CreateTransaction(NewWallet(secp265k1).PubKeyHex(), 100, mockedLister)

		mockedLister = new(MockedListing)
		blockchain = &listing.Blockchain{Chain: addEachTxToEachBlock(txA, txB, txC)}
		mockedLister.On("GetBlockchain").Return(blockchain)
		txD, _ := calculatingWallet.CreateTransaction(NewWallet(secp265k1).PubKeyHex(), 100, mockedLister)

		mockedLister = new(MockedListing)
		rewardTx, _ := CreateRewardTransaction(calculatingWallet)
		blockchain.Chain = append(blockchain.Chain, addTxsToOneBlock(txD, rewardTx))
		mockedLister.On("GetBlockchain").Return(blockchain)
		txE, _ := NewWallet(secp265k1).CreateTransaction(calculatingWallet.PubKeyHex(), 75, mockedLister)

		mockedLister = new(MockedListing)
		blockchain.Chain = append(blockchain.Chain, addTxsToOneBlock(txE))
		mockedLister.On("GetBlockchain").Return(blockchain)

		// perform test
		receivedBalance := CalculateBalance(mockedLister, calculatingWallet.PubKeyHex())

		// test verification
		assert.Equal(940, int(receivedBalance))
	})
}
