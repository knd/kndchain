package main

import (
	"log"
	"os"

	"github.com/knd/kndchain/pkg/networking/pubsub"

	"github.com/knd/kndchain/pkg/crypto"
	"github.com/knd/kndchain/pkg/listing"
	"github.com/knd/kndchain/pkg/mining"
	"github.com/knd/kndchain/pkg/storage/leveldb"
	"github.com/knd/kndchain/pkg/validating"
	"github.com/knd/kndchain/pkg/wallet"

	"github.com/knd/kndchain/pkg/calculating"
)

func main() {
	calculator := calculating.NewService(1000)
	repository := leveldb.NewRepository("/home/ubuntu/kndchainDatadir")
	lister := listing.NewService(repository)
	validator := validating.NewService(lister, calculator, "MINER_REWARD", 5)
	miningService := mining.NewService(repository, lister, validator, 10*1000)

	// Load wallet
	minerWallet := wallet.LoadWallet(
		crypto.NewSecp256k1Generator(),
		calculator,
		1000,
		"/home/ubuntu/kndchainKeys",
		os.Args[1])

	// Open Redis connection
	transactionPool := wallet.NewTransactionPool(lister)
	p2pComm := pubsub.NewService(
		lister,
		miningService,
		transactionPool,
		"kndchain",
		"kndchaintransactions",
		"http://ec2-52-27-194-217.us-west-2.compute.amazonaws.com:6379")
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

	for {
		miner.Mine()
	}
}
