package memory

import (
	"time"
)

type input struct {
	Timestamp int64
	Amount    uint64
	Address   string
	Signature string
}

// Transaction in data
type Transaction struct {
	ID     string
	Input  input
	Output map[string]uint64
}

// Block represents a block in blockchain
type Block struct {
	Timestamp  time.Time
	LastHash   *string
	Hash       *string
	Data       []Transaction
	Nonce      uint32
	Difficulty uint32
}
