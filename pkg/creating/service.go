package creating

type Service interface {
	CreateGenesisBlock()
	CreateBlock(lastHash *string, data []string)
}
