package listing

// Repository provides access to the blockchain
type Repository interface {
	GetLastBlock() Block
	GetBlockCount() uint32
	GetBlocks() []Block
}

// Service provides block listing operations
type Service interface {
	GetLastBlock() Block
	GetBlockCount() uint32
	GetBlocks() []Block
}

type service struct {
	r Repository
}

// NewService creates a listing service with necessary dependencies
func NewService(r Repository) Service {
	return &service{r}
}

// GetLastBlock returns the last mined block in the blockchain
func (s *service) GetLastBlock() Block {
	return s.r.GetLastBlock()
}

// GetBlockCount returns the latest number of blocks in blockchain
func (s *service) GetBlockCount() uint32 {
	return s.r.GetBlockCount()
}

// GetBlocks returns a list of blocks from genesis block
func (s *service) GetBlocks() []Block {
	return s.r.GetBlocks()
}
