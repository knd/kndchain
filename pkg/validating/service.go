package validating

// Service provides blockchain validating operations
type Service interface {
	IsValidChain(bc *Blockchain) bool
}

type service struct{}

// NewService creates a validating service with necessary dependencies
func NewService() Service {
	return &service{}
}

// IsValidChain returns true if list of blocks compose valid blockchain
func (s *service) IsValidChain(bc *Blockchain) bool {
	// TODO: Impelment this method
	return false
}
