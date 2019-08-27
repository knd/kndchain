package syncing

// Blockchain represents a chain of mined blocks
type Blockchain struct {
	Chain []Block `json:"chain"`
}
