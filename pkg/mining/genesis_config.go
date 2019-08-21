package mining

type GenesisConfig struct {
	LastHash *string   `json:"lastHash"`
	Hash     *string   `json:"hash"`
	Data     *[]string `json:"data"`
}
