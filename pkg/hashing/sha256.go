package hashing

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/knd/kndchain/pkg/sorting"
)

// SHA256Hash returns a determistic 256 bit hash given a set of inputs
func SHA256Hash(inputs ...interface{}) string {
	h := sha256.New()

	var inputByteArray [][]byte
	for _, i := range inputs {
		inputByteArray = append(inputByteArray, []byte(fmt.Sprintf("%v", i)))
	}

	sortedByteArray := sorting.SortByteArrays(inputByteArray)
	sort.Sort(sortedByteArray)

	var unifiedBytes []byte
	for _, s := range sortedByteArray {
		unifiedBytes = append(unifiedBytes, s...)
	}

	h.Write(unifiedBytes)
	return hex.EncodeToString(h.Sum(nil))
}
