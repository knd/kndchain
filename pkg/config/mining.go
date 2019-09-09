package config

const (
	// MineRate (1000 milliseconds) adjusts the difficulty of mining operation
	MineRate int = 1000
)

const (
	// DefaultGenesisLastHash is default last genesis block hash if not given from genesis config
	DefaultGenesisLastHash string = "0x000"

	// DefaultGenesisHash is default genesis hash if not given from genesis config
	DefaultGenesisHash string = "0x000"

	// DefaultGenesisDifficulty is default difficulty in genesis block
	DefaultGenesisDifficulty uint32 = 3

	// DefaultGenesisNonce is default nonce in genesis block
	DefaultGenesisNonce uint32 = 0
)
