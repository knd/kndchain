package calculating

// Input of transaction
type Input struct {
	Timestamp int64  `json:"timestamp"`
	Amount    uint64 `json:"amount"`
	Address   string `json:"address"`
	Signature string `json:"sig"`
}

// Transaction in block data
type Transaction struct {
	ID     string            `json:"id"`
	Input  Input             `json:"input"`
	Output map[string]uint64 `json:"output"`
}

// Block represents a block in blockchain
type Block struct {
	Timestamp  int64         `json:"timestamp"`
	LastHash   *string       `json:"lastHash"`
	Hash       *string       `json:"hash"`
	Data       []Transaction `json:"data"`
	Nonce      uint32        `json:"nonce"`
	Difficulty uint32        `json:"difficulty"`
}
