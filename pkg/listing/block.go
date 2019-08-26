package listing

import (
	"time"
)

// Block represents a block in blockchain
type Block struct {
	Timestamp  time.Time `json:"timestamp"`
	LastHash   *string   `json:"lastHash"`
	Hash       *string   `json:"hash"`
	Data       []string  `json:"data"`
	Nonce      uint32    `json:"nonce"`
	Difficulty uint32    `json:"difficulty"`
}
