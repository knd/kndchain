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

// Verify pubkey, msg, and signature
func (m *MockedKeyPairGenerator) Verify(pubKey, msg, signature []byte) bool {
	args := m.Called(pubKey, msg, signature)
	return args.Bool(0)
}

// Sign msg with private key to obtain signature
func (m *MockedKeyPairGenerator) Sign(msg, privKey []byte) ([]byte, error) {
	args := m.Called(msg, privKey)
	return args.Get(0).([]byte), args.Error(0)
}
