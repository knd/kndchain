package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/knd/kndchain/pkg/calculating"
	"github.com/knd/kndchain/pkg/crypto"
	"github.com/knd/kndchain/pkg/http/rest"
	"github.com/knd/kndchain/pkg/listing"
	"github.com/knd/kndchain/pkg/mining"
	"github.com/knd/kndchain/pkg/networking/pubsub"
	"github.com/knd/kndchain/pkg/storage/leveldb"
	"github.com/knd/kndchain/pkg/validating"
	"github.com/knd/kndchain/pkg/wallet"
)

const (
	initialBalance uint64 = 1000

	blockRewardAddress string = "MINER_REWARD"
	blockRewardAmount  uint64 = 5
	blockMiningRate    int64  = 10 * 1000 // 10 seconds

	p2pBlockChannel string = "kndchain"
	p2pTxChannel    string = "kndchaintransactions"
	p2pURI          string = "redis://@localhost:6379"
)

func main() {
	enableMining := flag.Bool("mining", false, "enable mining option")
	address := flag.String("address", "", "provide pubkeyhex/ address used for transactions or mining reward")
	chainDatadir := flag.String("chainDatadir", "/tmp/kndchainDatadir", "directory to store blockchain data")
	keysDatadir := flag.String("keysDatadir", "/tmp/kndchainKeys", "directory to store keys")
	flag.Parse()

	calculator := calculating.NewService(initialBalance)
	repository := leveldb.NewRepository(*chainDatadir)
	lister := listing.NewService(repository)
	validator := validating.NewService(lister, calculator, blockRewardAddress, blockRewardAmount)
	miningService := mining.NewService(repository, lister, validator, blockMiningRate)

	var wal wallet.Wallet
	if len(*address) != 0 {
		// Load wallet
		wal = wallet.LoadWallet(
			crypto.NewSecp256k1Generator(),
			calculator,
			lister,
			*keysDatadir,
			*address)
	} else {
		wal = wallet.NewWallet(
			crypto.NewSecp256k1Generator(),
			calculator,
			initialBalance,
			keysDatadir)
		log.Printf("Created new pubkey=%s, in %s", wal.PubKeyHex(), *keysDatadir)
	}

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

	if *enableMining {
		// Create miner
		miner := NewMiner(
			miningService,
			lister,
			transactionPool,
			wal,
			p2pComm,
			blockRewardAddress,
			blockRewardAmount)

		// Create genesis block
		if lister.GetBlockCount() == 0 {
			log.Println("Creating genesis block")
			genesisBlock, _ := mining.CreateGenesisBlock("0x000", "0x000", 20, 0)
			miningService.AddBlock(genesisBlock)
			p2pComm.BroadcastBlockchain(lister.GetBlockchain())
		}

		// Start mining
		var durations []float64
		go func() {
			for {
				lastBlock := lister.GetLastBlock()

				minedBlock, _ := miner.Mine()

				durationDiff := minedBlock.Timestamp - lastBlock.Timestamp
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
		}()
	}

	router := rest.Handler(lister, miningService, p2pComm, transactionPool, wal, calculator)
	log.Println("Serving now on http://localhost:3001")
	log.Fatal(http.ListenAndServe(":3001", router))
}
