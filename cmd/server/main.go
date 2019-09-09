package main

import (
	"fmt"
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
)

const (
	// PubSub will p2pCommunicate with peers via Redis pubsub
	PubSub Type = iota
)

// BeaconNodePort is port on which the beacon node is running
const BeaconNodePort int = 3000

func main() {
	storageType := Memory    // hard-coded
	networkingType := PubSub // hard-coded

	var miningService mining.Service
	var listingService listing.Service
	var validatingService validating.Service
	var calculatingService calculating.Service
	var p2pComm pubsub.Service

	switch storageType {
	case Memory:
		r := memory.NewRepository()

		validatingService = validating.NewService()
		calculatingService = calculating.NewService()
		listingService = listing.NewService(r)
		miningService = mining.NewService(r, listingService, validatingService, nil)
	}

	minerWallet := wallet.NewWallet(crypto.NewSecp256k1Generator(), calculatingService)
	transactionPool := wallet.NewTransactionPool(listingService)

	switch networkingType {
	case PubSub:
		p2pComm = pubsub.NewService(listingService, miningService, transactionPool)

		p2pComm.Connect()
		defer p2pComm.Disconnect()

		err := p2pComm.SubscribePeers()
		if err != nil {
			log.Fatal(err)
		}
	}

	miner := miner.NewMiner(miningService, listingService, transactionPool, minerWallet, p2pComm)
	router := rest.Handler(listingService, miningService, p2pComm, transactionPool, minerWallet, miner, calculatingService)

	var port int
	if len(os.Args) > 1 && os.Args[1] == "beacon" {
		port = BeaconNodePort

		// Add genesis block then broadcast to peers
		genesisBlock, _ := mining.CreateGenesisBlock(nil)
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
