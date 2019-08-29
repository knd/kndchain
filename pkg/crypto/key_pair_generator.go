package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"log"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/knd/kndchain/pkg/hashing"
)

// Secp256k1Generator provides secp256k1 operations
type Secp256k1Generator struct {
}

// NewSecp256k1Generator creates a secp256k1 generator
func NewSecp256k1Generator() *Secp256k1Generator {
	return &Secp256k1Generator{}
}

// Generate returns generated pubKey and privKey
func (s *Secp256k1Generator) Generate() (pubKey, privKey []byte) {
	key, err := ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	pubKey = elliptic.Marshal(secp256k1.S256(), key.X, key.Y)

	privKey = make([]byte, 32)
	blob := key.D.Bytes()
	copy(privKey[32-len(blob):], blob)

	return pubKey, privKey
}

// Verify checks that the given pubKey created signature over msg
func (s *Secp256k1Generator) Verify(pubKey, msg, signature []byte) bool {
	msgHash, err := sha256Hash(msg)
	if err != nil {
		return false
	}

	// VerifySignature checks that the given pubkey created signature over message. The signature should be in [R || S] format.
	return secp256k1.VerifySignature(pubKey, msgHash, signature[:64])
}

// Sign creates a recoverable ECDSA signature
func (s *Secp256k1Generator) Sign(msg, privKey []byte) ([]byte, error) {
	msgHash, err := sha256Hash(msg)
	if err != nil {
		return nil, err
	}

	// Sign creates a recoverable ECDSA signature. The produced signature is in the 65-byte [R || S || V] format where V is 0 or 1
	return secp256k1.Sign(msgHash, privKey)
}

func sha256Hash(data []byte) ([]byte, error) {
	dataHash := hashing.SHA256Hash(data)
	return hex.DecodeString(dataHash)
}
