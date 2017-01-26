package util

import (
	"golang.org/x/crypto/sha3"
)

var NewSha256 = sha3.New256

func Checksum256(data []byte) []byte {
	digest := sha3.Sum256(data)
	return digest[:]
}
