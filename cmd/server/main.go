package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/knd/kndchain/pkg/crypto"

	"github.com/knd/kndchain/pkg/wallet"

	"github.com/knd/kndchain/pkg/http/rest"
	"github.com/knd/kndchain/pkg/listing"
	"github.com/knd/kndchain/pkg/mining"
	"github.com/knd/kndchain/pkg/networking/pubsub"
	"github.com/knd/kndchain/pkg/storage/memory"
	"github.com/knd/kndchain/pkg/syncing"
	"github.com/knd/kndchain/pkg/validating"
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
	// PubSub will communicate with peers via Redis pubsub
	PubSub Type = iota
)

// BeaconNodePort is port on which the beacon node is running
const BeaconNodePort int = 3000

func main() {
	storageType := Memory    // hard-coded
	networkingType := PubSub // hard-coded

	var miner mining.Service
	var lister listing.Service
	var validator validating.Service
	var comm pubsub.Service

	switch storageType {
	case Memory:
		r := memory.NewRepository()

		validator = validating.NewService()
		lister = listing.NewService(r)
		miner = mining.NewService(r, lister, validator, nil)
	}

	switch networkingType {
	case PubSub:
		comm = pubsub.NewService(lister, miner)

		comm.Connect()
		defer comm.Disconnect()

		err := comm.SubscribePeers()
		if err != nil {
			log.Fatal(err)
		}
	}

	w := wallet.NewWallet(crypto.NewSecp256k1Generator())
	transactionPool := wallet.NewTransactionPool()
	router := rest.Handler(lister, miner, comm, transactionPool, w)

	var port int
	if len(os.Args) > 1 && os.Args[1] == "beacon" {
		port = BeaconNodePort

		// Add genesis block then broadcast to peers
		genesisBlock, _ := mining.CreateGenesisBlock(nil)
		miner.AddBlock(genesisBlock)
		comm.BroadcastBlockchain(lister.GetBlockchain())

	} else {
		// Generate port from 3000 - 4000
		rand.Seed(time.Now().UnixNano())
		portShuffle := rand.Intn(1000)
		port = 3000 + portShuffle

		log.Printf("Syncing blockchain. Current chain len: %d", lister.GetBlockCount())
		syncer := syncing.NewService(miner)
		err := syncer.SyncBlockchain(fmt.Sprintf("http://localhost:%d/api/blocks", BeaconNodePort))
		if err != nil {
			log.Println(err)
		}
		log.Printf("Syncing done. Synced chain len: %d", lister.GetBlockCount())
	}

	fmt.Printf("Serving now on http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
