package listing

// Repository provides access to the blockchain
type Repository interface {
	GetLastBlock() Block
	GetBlockCount() uint32
}

// Service provides block listing operations
type Service interface {
	GetLastBlock() Block
	GetBlockCount() uint32
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