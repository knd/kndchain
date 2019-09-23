package memory

// Input of transaction
type Input struct {
	Timestamp int64
	Amount    uint64
	Address   string
	Signature string
}

// Transaction in data
type Transaction struct {
	ID     string
	Input  Input
	Output map[string]uint64
}

// Block represents a block in blockchain
type Block struct {
	Timestamp  int64
	LastHash   *string
	Hash       *string
	Data       []Transaction
	Nonce      uint32
	Difficulty uint32
}
