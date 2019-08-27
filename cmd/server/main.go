package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/knd/kndchain/pkg/http/rest"
	"github.com/knd/kndchain/pkg/listing"
	"github.com/knd/kndchain/pkg/mining"
	"github.com/knd/kndchain/pkg/networking/pubsub"
	"github.com/knd/kndchain/pkg/storage/memory"
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

	// Add genesis block then broadcast to peers
	genesisBlock, _ := mining.CreateGenesisBlock(nil)
	miner.AddBlock(genesisBlock)
	comm.BroadcastBlockchain(lister.GetBlockchain())

	router := rest.Handler(lister, miner)

	// Generate port from 3000 - 4000
	rand.Seed(time.Now().UnixNano())
	portShuffle := rand.Intn(1000)
	randomPort := 3000 + portShuffle
	fmt.Printf("Serving now on http://localhost:%d\n", randomPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", randomPort), router))
}
