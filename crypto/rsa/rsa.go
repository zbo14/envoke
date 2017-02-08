package rsa

import (
	gocrypto "crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	. "github.com/zbo14/envoke/common"
	"github.com/zbo14/envoke/crypto/crypto"
)

const (
	E              = 65537
	KEY_SIZE       = 128
	PRIVKEY        = "RSA PRIVATE KEY"
	PUBKEY         = "RSA PUBLIC KEY"
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
	p []byte
}

func NewPrivateKey(inner rsa.PrivateKey) *PrivateKey {
	if len(inner.N.Bytes()) != KEY_SIZE {
		panic(ErrInvalidSize.Error())
	}
	// TODO: check private exponent?
	inner.E = E
	return &PrivateKey{inner}
}

func NewPublicKey(inner rsa.PublicKey) *PublicKey {
	if len(inner.N.Bytes()) != KEY_SIZE {
		panic(ErrInvalidSize.Error())
	}
	inner.E = E
	return &PublicKey{inner}
}

func NewSignature(inner []byte) *Signature {
	if len(inner) != SIGNATURE_SIZE {
		panic(ErrInvalidSize.Error())
	}
	return &Signature{inner}
}

func GenerateKeypair() (*PrivateKey, *PublicKey) {
	inner, err := rsa.GenerateKey(rand.Reader, KEY_SIZE*8)
	Check(err)
	priv := NewPrivateKey(*inner)
	pub := NewPublicKey(inner.PublicKey)
	return priv, pub
}

func NewPSSOptions() *rsa.PSSOptions {
	opts := new(rsa.PSSOptions)
	opts.SaltLength = SALT_SIZE
	opts.Hash = gocrypto.SHA256
	return opts
}

// PrivKey

func (_ *PrivateKey) IsPrivateKey() {}

func (priv *PrivateKey) Sign(message []byte) crypto.Signature {
	hash := NewSha256()
	hash.Write(message)
	hashed := hash.Sum(nil)
	opts := NewPSSOptions()
	inner, err := rsa.SignPSS(rand.Reader, &priv.inner, gocrypto.SHA256, hashed, opts)
	Check(err)
	return NewSignature(inner)
}

func (priv *PrivateKey) MarshalPEM() []byte {
	p := x509.MarshalPKCS1PrivateKey(&priv.inner)
	b := BlockPEM(p, PRIVKEY)
	return EncodePEM(b)
}

func (priv *PrivateKey) UnmarshalPEM(pem []byte) error {
	b, _ := DecodePEM(pem)
	if b.Type != PRIVKEY {
		panic(ErrInvalidType.Error())
	}
	inner, err := x509.ParsePKCS1PrivateKey(b.Bytes)
	Check(err)
	priv.inner = *inner
	return nil
}

func (priv *PrivateKey) Public() crypto.PublicKey {
	inner := priv.inner.PublicKey
	return NewPublicKey(inner)
}

// PubKey

func (_ *PublicKey) IsPublicKey() {}

func (pub *PublicKey) Verify(message []byte, sig crypto.Signature) bool {
	hash := NewSha256()
	hash.Write(message)
	hashed := hash.Sum(nil)
	opts := NewPSSOptions()
	err := rsa.VerifyPSS(&pub.inner, gocrypto.SHA256, hashed, sig.Bytes(), opts)
	return err == nil
}

func (pub *PublicKey) MarshalPEM() []byte {
	p, err := x509.MarshalPKIXPublicKey(&pub.inner)
	Check(err)
	b := BlockPEM(p, PUBKEY)
	return EncodePEM(b)
}

func (pub *PublicKey) UnmarshalPEM(pem []byte) error {
	b, _ := DecodePEM(pem)
	if b.Type != PUBKEY {
		panic(ErrInvalidType.Error())
	}
	inner, err := x509.ParsePKIXPublicKey(b.Bytes)
	Check(err)
	pub.inner = *inner.(*rsa.PublicKey)
	return nil
}

// Returns value of public modulus as a big-endian byte slice
func (pub *PublicKey) Bytes() []byte {
	if pub == nil {
		return nil
	}
	return pub.inner.N.Bytes()
}

func (pub *PublicKey) FromBytes(p []byte) error {
	if len(p) != KEY_SIZE {
		return ErrInvalidSize
	}
	pub.inner.E = E
	pub.inner.N = BigIntFromBytes(p)
	return nil
}

func (pub *PublicKey) String() string {
	return BytesToB58(pub.Bytes())
}

func (pub *PublicKey) FromString(str string) error {
	p := BytesFromB58(str)
	if size := len(p); size != KEY_SIZE {
		return ErrInvalidSize
	}
	return pub.FromBytes(p)
}

func (pub *PublicKey) MarshalJSON() ([]byte, error) {
	if pub == nil {
		return nil, nil
	}
	str := pub.String()
	p := MustMarshalJSON(str)
	return p, nil
}

func (pub *PublicKey) UnmarshalJSON(inner []byte) error {
	var str string
	if err := UnmarshalJSON(inner, &str); err != nil {
		return err
	}
	return pub.FromString(str)
}

// Signature

func (_ *Signature) IsSignature() {}

func (sig *Signature) Bytes() []byte {
	if sig == nil {
		return nil
	}
	return sig.p
}

func (sig *Signature) FromBytes(p []byte) error {
	if len(p) != SIGNATURE_SIZE {
		return ErrInvalidSize
	}
	sig.p = make([]byte, SIGNATURE_SIZE)
	copy(sig.p, p)
	return nil
}

func (sig *Signature) String() string {
	return BytesToB58(sig.Bytes())
}

func (sig *Signature) FromString(str string) error {
	p := BytesFromB58(str)
	if len(p) != SIGNATURE_SIZE {
		return ErrInvalidSize
	}
	return sig.FromBytes(p)
}

func (sig *Signature) MarshalJSON() ([]byte, error) {
	if sig == nil {
		return nil, nil
	}
	str := sig.String()
	p := MustMarshalJSON(str)
	return p, nil
}

func (sig *Signature) UnmarshalJSON(inner []byte) error {
	var str string
	if err := UnmarshalJSON(inner, &str); err != nil {
		return err
	}
	return sig.FromString(str)
}
