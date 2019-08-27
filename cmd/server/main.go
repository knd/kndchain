package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/knd/kndchain/pkg/http/rest"
	"github.com/knd/kndchain/pkg/listing"
	"github.com/knd/kndchain/pkg/mining"
	"github.com/knd/kndchain/pkg/storage/memory"
	"github.com/knd/kndchain/pkg/validating"
)

type Type int

const (
	// JSON will store data in JSON files saved on disk
	JSON Type = iota

	// Memory will store data in memory
	Memory
)

func main() {
	storageType := Memory

	var miner mining.Service
	var lister listing.Service
	var validator validating.Service

	switch storageType {
	case Memory:
		r := memory.NewRepository()

		validator = validating.NewService()
		lister = listing.NewService(r)
		miner = mining.NewService(r, lister, validator, nil)
	}

	genesisBlock, _ := mining.CreateGenesisBlock(nil)
	miner.AddBlock(genesisBlock)

	router := rest.Handler(lister, miner)

	fmt.Println("Serving now on http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}
