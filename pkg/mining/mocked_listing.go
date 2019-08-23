package mining

import (
	"github.com/knd/kndchain/pkg/listing"
	"github.com/stretchr/testify/mock"
)

// MockedListing is a mocked object that implememnts listing.Service
type MockedListing struct {
	mock.Mock
}

// GetLastBlock adds mined block to blockchain
func (m *MockedListing) GetLastBlock() listing.Block {
	args := m.Called()
	return args.Get(0).(listing.Block)
}

// GetBlockCount returns the latest block count in blockchain
func (m *MockedListing) GetBlockCount() uint32 {
	args := m.Called()
	return uint32(args.Int(0))
}
