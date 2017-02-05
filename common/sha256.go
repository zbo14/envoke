package common

import (
	"golang.org/x/crypto/sha3"
)

var NewSha256 = sha3.New256

func Checksum256(p []byte) []byte {
	digest := sha3.Sum256(p)
	return digest[:]
}

func Shake256(p []byte, size int) []byte {
	hash := make([]byte, size)
	sha3.ShakeSum256(hash, p)
	return hash
}
