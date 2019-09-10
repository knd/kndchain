package sorting

import (
	"bytes"
	"log"
)

// SortByteArrays implements `Interface` in sort package.
type SortByteArrays [][]byte

func (b SortByteArrays) Len() int {
	return len(b)
}

func (b SortByteArrays) Less(i, j int) bool {
	// bytes package already implements Comparable for []byte.
	switch bytes.Compare(b[i], b[j]) {
	case -1:
		return true
	case 0, 1:
		return false
	default:
		log.Println("not fail-able with `bytes.Comparable` bounded [-1, 1].")
		return false
	}
}

func (b SortByteArrays) Swap(i, j int) {
	b[j], b[i] = b[i], b[j]
}
