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
type Config struct {
	Mining        mine        `json:"mining"`
	Syncing       sync        `json:"syncing"`
	Transaction   transaction `json:"transaction"`
	Wallet        wal         `json:"wallet"`
	PathToDatadir string      `json:"pathToDatadir"`
	PathToKeysdir string      `json:"pathToKeysdir"`
}

func initConfig() *Config {
	configFile, err := os.Open("./cmd/production/config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer configFile.Close()

	byteValue, _ := ioutil.ReadAll(configFile)

	var config Config
	json.Unmarshal(byteValue, &config)

	return &config
}
