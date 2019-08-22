package mining

import (
	"github.com/stretchr/testify/mock"
)

// MockedRepository is a mocked object that implememnts mining.Repository
type MockedRepository struct {
	mock.Mock
}

// AddBlock adds mined block to blockchain
func (m *MockedRepository) AddBlock(minedBlock *Block) error {
	args := m.Called(minedBlock)
	return args.Error(0)
}
