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
