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

func TestHashing_HashYieldsDeterministicResultOfMap(t *testing.T) {
	mapA := make(map[string]uint64)
	mapA["test0"] = 1
	mapA["test3"] = 2

	mapB := make(map[string]uint64)
	mapB["test3"] = 2
	mapB["test0"] = 1

	assert.Equal(t, SHA256Hash(mapA), SHA256Hash(mapB))
}

func TestHashing_HashYeildsDeterministicResultOfStruct(t *testing.T) {
	type SA struct {
		Str string
	}
	type SB struct {
		Num int
	}
	structA := SA{Str: "test"}
	structB := SB{Num: 1}

	assert.Equal(t, SHA256Hash(structA, structB), SHA256Hash(SB{Num: 1}, SA{Str: "test"}))
}
