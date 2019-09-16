package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/knd/kndchain/pkg/calculating"
	"github.com/knd/kndchain/pkg/crypto"
	"github.com/knd/kndchain/pkg/http/rest"
	"github.com/knd/kndchain/pkg/listing"
	"github.com/knd/kndchain/pkg/miner"
	"github.com/knd/kndchain/pkg/mining"
	"github.com/knd/kndchain/pkg/networking/pubsub"
	"github.com/knd/kndchain/pkg/storage/leveldb"
	"github.com/knd/kndchain/pkg/storage/memory"
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

const (
	// PubSub will p2pCommunicate with peers via Redis pubsub
	PubSub Type = iota
)

// BeaconNodePort is port on which the beacon node is running
const BeaconNodePort int = 3000

type mine struct {
	MineRate          int64  `json:"mineRate"`
	GenesisLastHash   string `json:"genesisLastHash"`
	GenesisHash       string `json:"genesisHash"`
	GenesisDifficulty uint32 `json:"genesisDifficulty"`
	GenesisNonce      uint32 `json:"genesisNonce"`
}

type sync struct {
	ChannelPubSub       string `json:"channelPubSub"`
	ChannelTransactions string `json:"channelTransactions"`
	URLPubSub           string `json:"urlPubSub"`
}

type transaction struct {
	RewardTxInputAddress string `json:"rewardTxInputAddress"`
	MiningReward         uint64 `json:"miningReward"`
}

type wal struct {
	InitialBalance uint64 `json:"initialBalance"`
}

// Config feeds config from json
var Config struct {
	Mining        mine        `json:"mining"`
	Syncing       sync        `json:"syncing"`
	Transaction   transaction `json:"transaction"`
	Wallet        wal         `json:"wallet"`
	PathToDatadir string      `json:"pathToDatadir"`
	PathToKeysdir string      `json:"pathToKeysdir"`
}

func initConfig() {
	configFile, err := os.Open("./cmd/develop/config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer configFile.Close()

	byteValue, _ := ioutil.ReadAll(configFile)

	json.Unmarshal(byteValue, &Config)
}

func main() {
	initConfig()
	storageType := LevelDB   // hard-coded
	networkingType := PubSub // hard-coded

	var miningService mining.Service
	var listingService listing.Service
	var validatingService validating.Service
	var calculatingService calculating.Service
	var p2pComm pubsub.Service

	switch storageType {
	case Memory:
		repository := memory.NewRepository()
		listingService = listing.NewService(repository)
		calculatingService = calculating.NewService(Config.Wallet.InitialBalance)
		validatingService = validating.NewService(
			listingService,
			calculatingService,
			Config.Transaction.RewardTxInputAddress,
			Config.Transaction.MiningReward)
		miningService = mining.NewService(
			repository,
			listingService,
			validatingService,
			Config.Mining.MineRate)
	case LevelDB:
		repository := leveldb.NewRepository(Config.PathToDatadir)
		defer repository.Close()

		listingService = listing.NewService(repository)
		calculatingService = calculating.NewService(Config.Wallet.InitialBalance)
		validatingService = validating.NewService(
			listingService,
			calculatingService,
			Config.Transaction.RewardTxInputAddress,
			Config.Transaction.MiningReward)
		miningService = mining.NewService(
			repository,
			listingService,
			validatingService,
			Config.Mining.MineRate)
	}

	minerWallet := wallet.NewWallet(crypto.NewSecp256k1Generator(), calculatingService, Config.Wallet.InitialBalance, &Config.PathToKeysdir)
	transactionPool := wallet.NewTransactionPool(listingService)

	switch networkingType {
	case PubSub:
		p2pComm = pubsub.NewService(listingService, miningService, transactionPool, Config.Syncing.ChannelPubSub, Config.Syncing.ChannelTransactions, Config.Syncing.URLPubSub)

		p2pComm.Connect()
		defer p2pComm.Disconnect()

		err := p2pComm.SubscribePeers()
		if err != nil {
			log.Fatal(err)
		}
	}

	miner := miner.NewMiner(miningService, listingService, transactionPool, minerWallet, p2pComm, Config.Transaction.RewardTxInputAddress, Config.Transaction.MiningReward)
	router := rest.Handler(listingService, miningService, p2pComm, transactionPool, minerWallet, miner, calculatingService, true)

	var port int
	if len(os.Args) > 1 && os.Args[1] == "beacon" {
		port = BeaconNodePort

		// Add genesis block then broadcast to peers
		genesisBlock, _ := mining.CreateGenesisBlock(Config.Mining.GenesisLastHash, Config.Mining.GenesisHash, Config.Mining.GenesisDifficulty, Config.Mining.GenesisNonce)
		miningService.AddBlock(genesisBlock)
		p2pComm.BroadcastBlockchain(listingService.GetBlockchain())

	} else {
		// Generate port from 3000 - 4000
		rand.Seed(time.Now().UnixNano())
		portShuffle := rand.Intn(1000)
		port = 3000 + portShuffle

		log.Printf("Syncing blockchain. Current chain len: %d", listingService.GetBlockCount())
		syncer := syncing.NewService(miningService, transactionPool)
		err := syncer.SyncBlockchain(fmt.Sprintf("http://localhost:%d/api/blocks", BeaconNodePort))
		if err != nil {
			log.Println(err)
		}
		log.Printf("Blockchain synced. Synced chain len: %d", listingService.GetBlockCount())

		log.Println("Syncing transaction transactionPool...")
		err = syncer.SyncTransactionPool(fmt.Sprintf("http://localhost:%d/api/transactions", BeaconNodePort))
		if err != nil {
			log.Println(err)
		}
		log.Printf("Transaction transactionPool synced")
	}

	fmt.Printf("Serving now on http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
