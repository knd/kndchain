package mining

import (
	"time"
)

// Block represents a block in blockchain
type Block struct {
	Timestamp time.Time
	LastHash  *string
	Hash      *string
	Data      []string
}
