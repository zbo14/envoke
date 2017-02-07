package crypto

import (
	. "github.com/zbo14/envoke/common"
	"golang.org/x/crypto/bcrypt"
)

// Interfaces

type PrivateKey interface {
	IsPrivateKey()
	Public() PublicKey
	Sign([]byte) Signature
}

type PublicKey interface {
	IsPublicKey()
	Bytes() []byte
	FromBytes([]byte) error
	FromString(string) error
	MarshalJSON() ([]byte, error)
	String() string
	UnmarshalJSON([]byte) error
	Verify([]byte, Signature) bool
}

type Signature interface {
	IsSignature()
	Bytes() []byte
	FromBytes([]byte) error
	FromString(string) error
	// MarshalJSON() ([]byte, error)
	String() string
	// UnmarshalJSON([]byte) error
}

// Generate secret from password using bcrypt

func GenerateSecret(password string) []byte {
	secret, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	Check(err)
	return secret
}
