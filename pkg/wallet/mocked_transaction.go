package wallet

import (
	"github.com/stretchr/testify/mock"
)

// MockedTransaction is a mocked transaction
type MockedTransaction struct {
	mock.Mock
}

// GetID returns tx id
func (m *MockedTransaction) GetID() string {
	args := m.Called()
	return args.String(0)
}

// GetInput returns tx input
func (m *MockedTransaction) GetInput() TxInput {
	args := m.Called()
	return args.Get(0).(TxInput)
}

// GetOutput returns tx output
func (m *MockedTransaction) GetOutput() TxOutput {
	args := m.Called()
	return args.Get(0).(TxOutput)
}
