package rsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	. "github.com/zbo14/envoke/common"
)

const (
	E              = 65537
	KEY_SIZE       = 128
	SALT_SIZE      = 32
	SIGNATURE_SIZE = KEY_SIZE
)

type PrivateKey struct {
	inner rsa.PrivateKey
}

type PublicKey struct {
	inner rsa.PublicKey
}

type Signature struct {
	inner []byte
}

func NewPrivateKey(inner rsa.PrivateKey) *PrivateKey {
	if size := len(inner.N.Bytes()); size != KEY_SIZE {
		Panicf("Expected key with size=%d; got size=%d\n", KEY_SIZE, size)
	}
	// TODO: check private exponent?
	inner.E = E
	return &PrivateKey{inner}
}

func NewPublicKey(inner rsa.PublicKey) *PublicKey {
	if size := len(inner.N.Bytes()); size != KEY_SIZE {
		Panicf("Expected key with size=%d; got size=%d\n", KEY_SIZE, size)
	}
	inner.E = E
	return &PublicKey{inner}
}

func NewSignature(inner []byte) *Signature {
	if size := len(inner); size != SIGNATURE_SIZE {
		Panicf("Expected signature with size=%d; got size=%d\n", SIGNATURE_SIZE, size)
	}
	return &Signature{inner}
}

func GenerateKey() *PrivateKey {
	inner, err := rsa.GenerateKey(rand.Reader, KEY_SIZE*8)
	Check(err)
	return NewPrivateKey(*inner)
}

func NewPSSOptions() *rsa.PSSOptions {
	opts := new(rsa.PSSOptions)
	opts.SaltLength = SALT_SIZE
	opts.Hash = crypto.SHA256
	return opts
}

// PrivKey

func (priv *PrivateKey) Sign(message []byte) *Signature {
	hash := NewSha256()
	hash.Write(message)
	hashed := hash.Sum(nil)
	opts := NewPSSOptions()
	inner, err := rsa.SignPSS(rand.Reader, &priv.inner, crypto.SHA256, hashed, opts)
	Check(err)
	return NewSignature(inner)
}

func (priv *PrivateKey) Public() *PublicKey {
	inner := priv.inner.PublicKey
	return NewPublicKey(inner)
}

func (pub *PublicKey) Verify(message []byte, sig *Signature) bool {
	hash := NewSha256()
	hash.Write(message)
	hashed := hash.Sum(nil)
	opts := NewPSSOptions()
	err := rsa.VerifyPSS(&pub.inner, crypto.SHA256, hashed, sig.inner, opts)
	return err == nil
}

// PubKey

// Returns value of public modulus as a big-endian byte slice
func (pub *PublicKey) Bytes() []byte {
	return pub.inner.N.Bytes()
}

func (pub *PublicKey) FromBytes(p []byte) error {
	if size := len(p); size != KEY_SIZE {
		return Errorf("Expected key with size=%d; got size=%d\n", KEY_SIZE, size)
	}
	pub.inner.E = E
	pub.inner.N = BigIntFromBytes(p)
	return nil
}

func (pub *PublicKey) ToB58() string {
	return BytesToB58(pub.Bytes())
}

func (pub *PublicKey) FromB58(b58 string) error {
	p := BytesFromB58(b58)
	if size := len(p); size != KEY_SIZE {
		return Errorf("Expected key with size=%d; got size=%d\n", KEY_SIZE, size)
	}
	return pub.FromBytes(p)
}

func (pub *PublicKey) ToHex() string {
	return BytesToHex(pub.Bytes())
}

func (pub *PublicKey) FromHex(hex string) error {
	p := BytesFromHex(hex)
	if size := len(p); size != KEY_SIZE {
		return Errorf("Expected key with size=%d; got size=%d\n", KEY_SIZE, size)
	}
	return pub.FromBytes(p)
}

func (pub *PublicKey) MarshalJSON() ([]byte, error) {
	b58 := pub.ToB58()
	p := MustMarshalJSON(b58)
	return p, nil
}

func (pub *PublicKey) UnmarshalJSON(inner []byte) error {
	var b58 string
	if err := UnmarshalJSON(inner, &b58); err != nil {
		return err
	}
	return pub.FromB58(b58)
}

// Signature

func (sig *Signature) Bytes() []byte {
	return sig.inner
}

func (sig *Signature) FromBytes(inner []byte) error {
	if size := len(inner); size != SIGNATURE_SIZE {
		return Errorf("Expected signature with size=%d; got size=%d\n", SIGNATURE_SIZE, size)
	}
	sig.inner = make([]byte, SIGNATURE_SIZE)
	copy(sig.inner, inner)
	return nil
}

func (sig *Signature) ToB58() string {
	return BytesToB58(sig.Bytes())
}

func (sig *Signature) FromB58(b58 string) error {
	inner := BytesFromB58(b58)
	if size := len(inner); size != SIGNATURE_SIZE {
		return Errorf("Expected signature with size=%d; got size=%d\n", SIGNATURE_SIZE, size)
	}
	return sig.FromBytes(inner)
}

func (sig *Signature) ToHex() string {
	return BytesToHex(sig.Bytes())
}

func (sig *Signature) FromHex(hex string) error {
	inner := BytesFromHex(hex)
	if size := len(inner); size != SIGNATURE_SIZE {
		return Errorf("Expected signature with size=%d; got size=%d\n", SIGNATURE_SIZE, size)
	}
	return sig.FromBytes(inner)
}

func (sig *Signature) MarshalJSON() ([]byte, error) {
	b58 := sig.ToB58()
	p := MustMarshalJSON(b58)
	return p, nil
}

func (sig *Signature) UnmarshalJSON(inner []byte) error {
	var b58 string
	if err := UnmarshalJSON(inner, &b58); err != nil {
		return err
	}
	return sig.FromB58(b58)
}
