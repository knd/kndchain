package syncing

type input struct {
	Timestamp int64  `json:"timestamp"`
	Amount    uint64 `json:"amount"`
	Address   string `json:"address"`
	Signature string `json:"sig"`
}

type output map[string]uint64

// Transaction to marshall in syncing service
type Transaction struct {
	ID     string `json:"id"`
	Input  input  `json:"input"`
	Output output `json:"output"`
}

type TransactionPool map[string]Transaction
