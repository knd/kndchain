package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/knd/kndchain/pkg/calculating"
	"github.com/knd/kndchain/pkg/crypto"
	"github.com/knd/kndchain/pkg/listing"
	"github.com/knd/kndchain/pkg/miner"
	"github.com/knd/kndchain/pkg/mining"
	"github.com/knd/kndchain/pkg/networking/pubsub"
	"github.com/knd/kndchain/pkg/storage/leveldb"
	"github.com/knd/kndchain/pkg/syncing"
	"github.com/knd/kndchain/pkg/validating"
	"github.com/knd/kndchain/pkg/wallet"
)

// Type indicates type of constants
type Type int

const (
	// JSON will store data in JSON files saved on disk
	JSON Type = iota

	// Memory will store data in memory
	Memory

	// LevelDB will store data in local key-value store
	LevelDB
)

func main() {
	config := initConfig()
	enableMining := len(os.Args) > 2 && os.Args[1] == "mining"
	var miningAddress string
	if len(os.Args) > 2 && enableMining {
		miningAddress = os.Args[2]
	}

	repository := leveldb.NewRepository(config.PathToDatadir)
	defer repository.Close()

	listingService := listing.NewService(repository)
	calculatingService := calculating.NewService(config.Wallet.InitialBalance)
	validatingService := validating.NewService(
		listingService,
		calculatingService,
		config.Transaction.RewardTxInputAddress,
		config.Transaction.MiningReward)
	miningService := mining.NewService(
		repository,
		listingService,
		validatingService,
		config.Mining.MineRate)

	transactionPool := wallet.NewTransactionPool(listingService)
	p2pComm := pubsub.NewService(
		listingService,
		miningService,
		transactionPool,
		"kndchain",
		"kndchaintransactions",
		config.Syncing.URLPubSub)

	p2pComm.Connect()
	defer p2pComm.Disconnect()
	err := p2pComm.SubscribePeers()
	if err != nil {
		log.Fatal(err)
	}

	var m miner.Miner
	if enableMining {
		minerWallet := wallet.LoadWallet(
			crypto.NewSecp256k1Generator(),
			calculatingService,
			config.Wallet.InitialBalance,
			config.PathToKeysdir,
			miningAddress)
		m = miner.NewMiner(
			miningService,
			listingService,
			transactionPool,
			minerWallet,
			p2pComm,
			config.Transaction.RewardTxInputAddress,
			config.Transaction.MiningReward)
	}

	// router := rest.Handler(
	// 	listingService,
	// 	miningService,
	// 	p2pComm,
	// 	transactionPool,
	// 	m,
	// 	calculatingService,
	// 	enableMining)

	log.Printf(
		"Syncing blockchain. Current chain len: %d", listingService.GetBlockCount())

	syncer := syncing.NewService(miningService, transactionPool)
	err = syncer.SyncBlockchain(
		fmt.Sprintf("%s/api/blocks", config.Syncing.URLBeaconNode))
	if err != nil {
		log.Println(err)
	}

	log.Printf(
		"Blockchain synced. Synced chain len: %d", listingService.GetBlockCount())

	if listingService.GetBlockCount() == 0 {
		// create genesis block
		genesisBlock, _ := mining.CreateGenesisBlock(config.Mining.GenesisLastHash, config.Mining.GenesisHash, config.Mining.GenesisDifficulty, config.Mining.GenesisNonce)
		miningService.AddBlock(genesisBlock)
		p2pComm.BroadcastBlockchain(listingService.GetBlockchain())
	}

	log.Println("Syncing transaction transactionPool...")

	err = syncer.SyncTransactionPool(
		fmt.Sprintf("%s/api/transactions", config.Syncing.URLBeaconNode))
	if err != nil {
		log.Println(err)
	}

	log.Printf("Transaction transactionPool synced")

	log.Println("Serving now on http://localhost:3001")
	log.Fatal(http.ListenAndServe(":3001", router))
}
