package listing

// Repository provides access to the blockchain
type Repository interface {
	GetLastBlock() Block
}

// Service provides block listing operations
type Service interface {
	GetLastBlock() Block
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
