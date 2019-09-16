package main

import (
	"log"
	"os"
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

func main() {
	calculator := calculating.NewService(1000)
	repository := leveldb.NewRepository("/tmp/kndchainDatadir")
	lister := listing.NewService(repository)
	validator := validating.NewService(lister, calculator, "MINER_REWARD", 5)
	miningService := mining.NewService(repository, lister, validator, 10*1000)

	// Load wallet
	minerWallet := wallet.LoadWallet(
		crypto.NewSecp256k1Generator(),
		calculator,
		1000,
		"/tmp/kndchainKeys",
		os.Args[1])

	// Open Redis connection
	transactionPool := wallet.NewTransactionPool(lister)
	p2pComm := pubsub.NewService(
		lister,
		miningService,
		transactionPool,
		"kndchain",
		"kndchaintransactions",
		"redis://@localhost:6379")
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
		"MINER_REWARD",
		5)

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
