package main

import (
	"log"
	"time"

	"github.com/knd/kndchain/pkg/calculating"
	"github.com/knd/kndchain/pkg/crypto"
	"github.com/knd/kndchain/pkg/listing"
	"github.com/knd/kndchain/pkg/mining"
	"github.com/knd/kndchain/pkg/networking/pubsub"
	"github.com/knd/kndchain/pkg/storage/leveldb"
	"github.com/knd/kndchain/pkg/validating"
	"github.com/knd/kndchain/pkg/wallet"
)

const (
	chainDatadir string = "/tmp/kndchainDatadir"
	keysDatadir  string = "/tmp/kndchainKeys"

	initialBalance uint64 = 1000

	blockRewardAddress string = "MINER_REWARD"
	blockRewardAmount  uint64 = 5
	blockMiningRate    int64  = 10 * 1000 // 10 seconds

	p2pBlockChannel string = "kndchain"
	p2pTxChannel    string = "kndchaintransactions"
	p2pURI          string = "redis://@localhost:6379"

	rewardPubKey string = "04749d91026def10f5c55170115b119291c2c9ddc9f8e009808a93fd8c7e4f3753d74ce24e89b7341c634d58e47765c2be8fdf2e9ccaca78f36c1aa7f6ca33b615"
)

func main() {
	calculator := calculating.NewService(initialBalance)
	repository := leveldb.NewRepository(chainDatadir)
	lister := listing.NewService(repository)
	validator := validating.NewService(lister, calculator, blockRewardAddress, blockRewardAmount)
	miningService := mining.NewService(repository, lister, validator, blockMiningRate)

	// Load wallet
	minerWallet := wallet.LoadWallet(
		crypto.NewSecp256k1Generator(),
		calculator,
		lister,
		keysDatadir,
		rewardPubKey)

	// Open Redis connection
	transactionPool := wallet.NewTransactionPool(lister)
	p2pComm := pubsub.NewService(
		lister,
		miningService,
		transactionPool,
		p2pBlockChannel,
		p2pTxChannel,
		p2pURI)
	p2pComm.Connect()
	defer p2pComm.Disconnect()
	err := p2pComm.SubscribePeers()
	if err != nil {
		log.Fatal(err)
	}

	// Create miner
	miner := NewMiner(
		miningService,
		lister,
		transactionPool,
		minerWallet,
		p2pComm,
		blockRewardAddress,
		blockRewardAmount)

	// Create genesis block
	if lister.GetBlockCount() == 0 {
		genesisBlock, _ := mining.CreateGenesisBlock("0x000", "0x000", 10, 0)
		miningService.AddBlock(genesisBlock)
		p2pComm.BroadcastBlockchain(lister.GetBlockchain())
	}

	var durations []float64
	for {
		lastBlock := lister.GetLastBlock()

		minedBlock, _ := miner.Mine()

		durationDiff := minedBlock.Timestamp.Sub(lastBlock.Timestamp)
		durationDiffInMillis := float64(durationDiff) / float64(time.Millisecond)

		durations = append(durations, durationDiffInMillis)
		var sumDuration float64
		for _, duration := range durations {
			sumDuration = sumDuration + duration
		}
		averageDuration := float64(sumDuration) / float64(len(durations))

		log.Printf(
			"Time to mine block: %.2f ms. Difficulty: %d. Average time: %.2f ms", durationDiffInMillis,
			minedBlock.Difficulty,
			averageDuration,
		)
	}
}
