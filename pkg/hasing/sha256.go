package hasing

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"
)

// SHA256Hash returns a determistic 256 bit hash given a set of inputs
func SHA256Hash(inputs ...string) string {
	h := sha256.New()
	sort.Strings(inputs)
	h.Write([]byte(strings.Join(inputs, "")))
	return hex.EncodeToString(h.Sum(nil))
}
