package wallet

import (
	"github.com/stretchr/testify/mock"
)

// MockedKeyPairGenerator simulates any key pair generator
type MockedKeyPairGenerator struct {
	mock.Mock
}

// Generate returns generated key pair pubKey, privKey
func (m *MockedKeyPairGenerator) Generate() (pubKey, privKey []byte) {
	args := m.Called()
	return args.Get(0).([]byte), args.Get(1).([]byte)
}
