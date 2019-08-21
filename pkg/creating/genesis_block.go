package creating

import "time"

var genesisHash = "0"

var GenesisBlock = Block{
	Timestamp: time.Now(),
	LastHash:  nil,
	Hash:      &genesisHash,
	Data:      []string{},
}
