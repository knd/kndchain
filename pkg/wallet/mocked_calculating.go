package wallet

import (
	"github.com/knd/kndchain/pkg/calculating"
	"github.com/stretchr/testify/mock"
)

// MockedCalculating provides access to calculating service
type MockedCalculating struct {
	mock.Mock
}

// Balance returns balance of address based on given blockchain history
func (m *MockedCalculating) Balance(address string, bc *calculating.Blockchain) uint64 {
	args := m.Called(address, bc)
	return args.Get(0).(uint64)
}

// BalanceByBlockIndex returns balance of address based on given blockchain history at block index
func (m *MockedCalculating) BalanceByBlockIndex(address string, bc *calculating.Blockchain, index int) uint64 {
	args := m.Called(address, bc, index)
	return args.Get(0).(uint64)
}
