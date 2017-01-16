package util

import (
	"golang.org/x/crypto/sha3"
)

func Checksum256(data []byte) []byte {
	hash := sha3.New256()
	hash.Write(data)
	return hash.Sum(nil)
}
