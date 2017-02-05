package ed25519

import (
	"bytes"
	. "github.com/zbo14/envoke/common"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ed25519"
)

const (
	PRIVKEY_SIZE   = ed25519.PrivateKeySize
	PUBKEY_SIZE    = ed25519.PublicKeySize
	SIGNATURE_SIZE = ed25519.SignatureSize
)

type PublicKey struct {
	inner ed25519.PublicKey
}

type PrivateKey struct {
	inner ed25519.PrivateKey
}

type Signature struct {
	inner []byte
}

func NewPrivateKey(inner ed25519.PrivateKey) (*PrivateKey, error) {
	if size := len(inner); size != PRIVKEY_SIZE {
		return nil, Errorf("Expected privkey with size=%d; got size=%d\n", PRIVKEY_SIZE, size)
	}
	return &PrivateKey{inner}, nil
}

func NewPublicKey(inner ed25519.PublicKey) (*PublicKey, error) {
	if size := len(inner); size != PUBKEY_SIZE {
		return nil, Errorf("Expected pubkey with size=%d; got size=%d\n", PUBKEY_SIZE, size)
	}
	return &PublicKey{inner}, nil
}

func NewSignature(inner []byte) (*Signature, error) {
	if size := len(inner); size != SIGNATURE_SIZE {
		return nil, Errorf("Expected signature with size=%d; got size=%d\n", SIGNATURE_SIZE, size)
	}
	return &Signature{inner}, nil
}

func GenerateKeypair(password string) (*PrivateKey, *PublicKey) {
	secret := GenerateSecret(password)
	buf := new(bytes.Buffer)
	buf.Write(secret)
	pubInner, privInner, err := ed25519.GenerateKey(buf)
	Check(err)
	priv, err := NewPrivateKey(privInner)
	Check(err)
	pub, err := NewPublicKey(pubInner)
	Check(err)
	return priv, pub
}

// Private Key

func (priv *PrivateKey) Sign(message []byte) *Signature {
	p := ed25519.Sign(priv.inner, message)
	sig, err := NewSignature(p)
	Check(err)
	return sig
}

func (priv *PrivateKey) Bytes() []byte {
	return priv.inner[:]
}

func (priv *PrivateKey) Public() *PublicKey {
	p := priv.inner.Public().(ed25519.PublicKey)
	pub, err := NewPublicKey(p)
	Check(err)
	return pub
}

func (priv *PrivateKey) ToB58() string {
	return BytesToB58(priv.Bytes())
}

func (priv *PrivateKey) FromB58(b58 string) error {
	inner := BytesFromB58(b58)
	if size := len(inner); size != PRIVKEY_SIZE {
		return Errorf("Expected privkey with size=%d; got size=%d\n", PRIVKEY_SIZE, size)
	}
	priv.inner = make([]byte, PRIVKEY_SIZE)
	copy(priv.inner, inner)
	return nil
}

func (priv *PrivateKey) ToHex() string {
	return BytesToHex(priv.Bytes())
}

func (priv *PrivateKey) FromHex(hex string) error {
	inner := BytesFromHex(hex)
	if size := len(inner); size != PRIVKEY_SIZE {
		return Errorf("Expected privkey with size=%d; got size=%d\n", PRIVKEY_SIZE, size)
	}
	priv.inner = make([]byte, PRIVKEY_SIZE)
	copy(priv.inner, inner)
	return nil
}

func (priv *PrivateKey) MarshalJSON() ([]byte, error) {
	b58 := priv.ToB58()
	p := MustMarshalJSON(b58)
	return p, nil
}

func (priv *PrivateKey) UnmarshalJSON(inner []byte) error {
	var b58 string
	if err := UnmarshalJSON(inner, &b58); err != nil {
		return err
	}
	return priv.FromB58(b58)
}

// Public Key

func (pub *PublicKey) Verify(message []byte, sig *Signature) bool {
	return ed25519.Verify(pub.inner, message, sig.inner)
}

func (pub *PublicKey) Bytes() []byte {
	return pub.inner[:]
}

func (pub *PublicKey) ToB58() string {
	return BytesToB58(pub.Bytes())
}

func (pub *PublicKey) FromB58(b58 string) error {
	inner := BytesFromB58(b58)
	if size := len(inner); size != PUBKEY_SIZE {
		return Errorf("Expected pubkey with size=%d; got size=%d\n", PUBKEY_SIZE, size)
	}
	pub.inner = make([]byte, PUBKEY_SIZE)
	copy(pub.inner, inner)
	return nil
}

func (pub *PublicKey) ToHex() string {
	return BytesToHex(pub.Bytes())
}

func (pub *PublicKey) FromHex(hex string) error {
	inner := BytesFromHex(hex)
	if size := len(inner); size != PUBKEY_SIZE {
		return Errorf("Expected pubkey with size=%d; got size=%d\n", PUBKEY_SIZE, size)
	}
	pub.inner = make([]byte, PUBKEY_SIZE)
	copy(pub.inner, inner)
	return nil
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
	return sig.inner[:]
}

func (sig *Signature) ToB58() string {
	return BytesToB58(sig.Bytes())
}

func (sig *Signature) FromB58(b58 string) error {
	inner := BytesFromB58(b58)
	if size := len(inner); size != SIGNATURE_SIZE {
		return Errorf("Expected signature with size=%d; got size=%d\n", SIGNATURE_SIZE, size)
	}
	sig.inner = make([]byte, SIGNATURE_SIZE)
	copy(sig.inner, inner)
	return nil
}

func (sig *Signature) ToHex() string {
	return BytesToHex(sig.Bytes())
}

func (sig *Signature) FromHex(hex string) error {
	inner := BytesFromHex(hex)
	if size := len(inner); size != SIGNATURE_SIZE {
		return Errorf("Expected signature with size=%d; got size=%d\n", SIGNATURE_SIZE, size)
	}
	sig.inner = make([]byte, SIGNATURE_SIZE)
	copy(sig.inner, inner)
	return nil
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

func GenerateSecret(password string) []byte {
	secret, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	Check(err)
	return secret
}
