package mining

// GenesisConfig contains config for the genesis block
type GenesisConfig struct {
	LastHash *string   `json:"lastHash"`
	Hash     *string   `json:"hash"`
	Data     *[]string `json:"data"`
}
