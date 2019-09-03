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
func (m *MockedTransaction) GetInput() Input {
	args := m.Called()
	return args.Get(0).(Input)
}

// GetOutput returns tx output
func (m *MockedTransaction) GetOutput() Output {
	args := m.Called()
	return args.Get(0).(Output)
}

// Append updates another tx receiver with another amount
func (m *MockedTransaction) Append(w Wallet, receiver string, amount uint64) error {
	args := m.Called(w, receiver, amount)
	return args.Error(0)
}
