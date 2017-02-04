package aes_gcm

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	. "github.com/zbo14/envoke/common"
)

const NONCE_SIZE = 12 //what should this be?

// AES-GCM should be used because the operation is an authenticated encryption
// algorithm designed to provide both data authenticity (integrity) as well as
// confidentiality.

// TODO: add additional (unencrypted) data

func Encrypt(key, plaintext []byte) []byte {
	block, err := aes.NewCipher(key)
	Check(err)
	nonce := make([]byte, NONCE_SIZE)
	ReadFull(rand.Reader, nonce)
	gcm, err := cipher.NewGCM(block)
	Check(err)
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext
}

func Decrypt(key, ciphertext []byte) []byte {
	block, err := aes.NewCipher(key)
	Check(err)
	gcm, err := cipher.NewGCM(block)
	Check(err)
	nonce := make([]byte, NONCE_SIZE)
	nonce = ciphertext[:NONCE_SIZE]
	ciphertext = ciphertext[NONCE_SIZE:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	Check(err)
	return plaintext
}
