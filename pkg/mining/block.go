package mining

import (
	"time"
)

type Block struct {
	Timestamp time.Time
	LastHash  *string
	Hash      *string
	Data      []string
}
