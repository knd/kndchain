package mining

import (
	"github.com/knd/kndchain/pkg/validating"
	"github.com/stretchr/testify/mock"
)

// MockedValidating is a mocked object that implememnts validating.Service
type MockedValidating struct {
	mock.Mock
}

// IsValidChain returns true if list of blocks compose valid blockchain
func (m *MockedValidating) IsValidChain(bc *validating.Blockchain) bool {
	args := m.Called(bc)
	return args.Bool(0)
}