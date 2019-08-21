package hashing

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHashing_HashSingleInput(t *testing.T) {
	assert.Equal(t, "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824", SHA256Hash("hello"))
}

func TestHashing_HashYieldsDeterministicResultRegardlessInputsOrder(t *testing.T) {
	now := time.Now()
	assert.Equal(t, SHA256Hash("hello", now, 2019), SHA256Hash(2019, "hello", now))
}
