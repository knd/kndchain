package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/knd/kndchain/pkg/calculating"
	"github.com/knd/kndchain/pkg/crypto"
	"github.com/knd/kndchain/pkg/http/rest"
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

type mine struct {
	MineRate          int64  `json:"mineRate"`
	GenesisLastHash   string `json:"genesisLastHash"`
	GenesisHash       string `json:"genesisHash"`
	GenesisDifficulty uint32 `json:"genesisDifficulty"`
	GenesisNonce      uint32 `json:"genesisNonce"`
}

type sync struct {
	URLPubSub     string `json:"urlPubSub"`
	URLBeaconNode string `json:"urlBeaconNode"`
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
	Mining      mine        `json:"mining"`
	Syncing     sync        `json:"syncing"`
	Transaction transaction `json:"transaction"`
	Wallet      wal         `json:"wallet"`
}

func initConfig() {
	configFile, err := os.Open("./cmd/production/config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer configFile.Close()

	byteValue, _ := ioutil.ReadAll(configFile)

	json.Unmarshal(byteValue, &Config)
}

func main() {
	initConfig()

	repository := leveldb.NewRepository("~/kndchainDatadir")
	defer repository.Close()

	listingService := listing.NewService(repository)
	calculatingService := calculating.NewService(Config.Wallet.InitialBalance)

	validatingService := validating.NewService(
		listingService,
		calculatingService,
		Config.Transaction.RewardTxInputAddress,
		Config.Transaction.MiningReward)

	miningService := mining.NewService(
		repository,
		listingService,
		validatingService,
		Config.Mining.MineRate)

	minerWallet := wallet.NewWallet(
		crypto.NewSecp256k1Generator(),
		calculatingService,
		Config.Wallet.InitialBalance)

	transactionPool := wallet.NewTransactionPool(listingService)

	p2pComm := pubsub.NewService(
		listingService,
		miningService,
		transactionPool,
		"kndchain",
		"kndchaintransactions",
		Config.Syncing.URLPubSub)

	p2pComm.Connect()
	defer p2pComm.Disconnect()

	err := p2pComm.SubscribePeers()
	if err != nil {
		log.Fatal(err)
	}

	miner := miner.NewMiner(
		miningService,
		listingService,
		transactionPool,
		minerWallet,
		p2pComm,
		Config.Transaction.RewardTxInputAddress,
		Config.Transaction.MiningReward)

	router := rest.Handler(
		listingService,
		miningService,
		p2pComm,
		transactionPool,
		minerWallet,
		miner,
		calculatingService)

	log.Printf(
		"Syncing blockchain. Current chain len: %d", listingService.GetBlockCount())

	syncer := syncing.NewService(miningService, transactionPool)
	err = syncer.SyncBlockchain(
		fmt.Sprintf("%s/api/blocks", Config.Syncing.URLBeaconNode))
	if err != nil {
		log.Println(err)
	}

	log.Printf(
		"Blockchain synced. Synced chain len: %d", listingService.GetBlockCount())

	log.Println("Syncing transaction transactionPool...")

	err = syncer.SyncTransactionPool(
		fmt.Sprintf("%s/api/transactions", Config.Syncing.URLBeaconNode))
	if err != nil {
		log.Println(err)
	}

	log.Printf("Transaction transactionPool synced")

	log.Println("Serving now on http://localhost:3001")
	log.Fatal(http.ListenAndServe(":3001", router))
}
