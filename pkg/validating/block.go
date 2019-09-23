package validating

// Block represents a block in blockchain
type Block struct {
	Timestamp  int64         `json:"timestamp"`
	LastHash   *string       `json:"lastHash"`
	Hash       *string       `json:"hash"`
	Data       []Transaction `json:"data"`
	Nonce      uint32        `json:"nonce"`
	Difficulty uint32        `json:"difficulty"`
}
