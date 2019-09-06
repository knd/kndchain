package validating

// Input of transaction
type Input struct {
	Timestamp int64  `json:"timestamp"`
	Amount    uint64 `json:"amount"`
	Address   string `json:"address"`
	Signature string `json:"sig"`
}

// Transaction in data
type Transaction struct {
	ID     string            `json:"id"`
	Input  Input             `json:"input"`
	Output map[string]uint64 `json:"output"`
}
