package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

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
	URLBeaconNode       string `json:"urlBeaconNode"`
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

	// TODO: create a wallet and mining service here

	/*
		// Add genesis block then broadcast to peers
		genesisBlock, _ := mining.CreateGenesisBlock(Config.Mining.GenesisLastHash, Config.Mining.GenesisHash, Config.Mining.GenesisDifficulty, Config.Mining.GenesisNonce)
		miningService.AddBlock(genesisBlock)
		p2pComm.BroadcastBlockchain(listingService.GetBlockchain())
	*/
}
