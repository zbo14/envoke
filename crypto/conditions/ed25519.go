package conditions

import (
	"github.com/zballs/go_resonate/crypto/ed25519"
	. "github.com/zballs/go_resonate/util"
)

// Ed25519

type fulfillmentEd25519 struct {
	pub    *ed25519.PublicKey
	sig    *ed25519.Signature
	weight int
}

func NewFulfillmentEd25519(msg []byte, priv *ed25519.PrivateKey, weight int) *fulfillmentEd25519 {
	pub := priv.Public()
	sig := priv.Sign(msg)
	return &fulfillmentEd25519{
		pub:    pub,
		sig:    sig,
		weight: weight,
	}
}

func (f *fulfillmentEd25519) Id() int { return ED25519_ID }

func (f *fulfillmentEd25519) Bitmask() int { return ED25519_BITMASK }

func (f *fulfillmentEd25519) Size() int { return ED25519_SIZE }

func (f *fulfillmentEd25519) Len() int { p, _ := f.MarshalBinary(); return len(p) }

func (f *fulfillmentEd25519) Weight() int { return f.weight }

func (f *fulfillmentEd25519) String() string {
	p := append(f.pub.Bytes(), f.sig.Bytes()...)
	b64 := Base64UrlEncode(p)
	return Sprintf("cf:%x:%s", ED25519_ID, b64)
}

func (f *fulfillmentEd25519) FromString(uri string) error {
	if !MatchString(FULFILLMENT_REGEX, uri) {
		return Error("Uri does not match fulfillment regex")
	}
	parts := Split(uri, ":")
	id, err := ParseUint16(parts[1])
	if err != nil {
		return err
	}
	if int(id) != ED25519_ID {
		return Errorf("Expected type_id=%d; got type_id=%d\n", ED25519_ID, id)
	}
	p, err := Base64UrlDecode(parts[3])
	if err != nil {
		return err
	}
	if size := len(p); size != ED25519_SIZE {
		return Errorf("Expected fulfillment_size=%d; got fulfillment_size=%d\n", ED25519_SIZE, size)
	}
	f.pub, _ = ed25519.NewPublicKey(p[:ed25519.PUBKEY_SIZE])
	f.sig, _ = ed25519.NewSignature(p[ed25519.PUBKEY_SIZE:])
	return nil
}

func (f *fulfillmentEd25519) Hash() []byte {
	return f.pub.Bytes()
}

func (f *fulfillmentEd25519) Condition() *Condition {
	return NewCondition(ED25519_ID, ED25519_BITMASK, f.Hash(), f.Size(), f.Weight())
}

func (f *fulfillmentEd25519) Validate(msg []byte) bool {
	return f.pub.Verify(msg, f.sig)
}

func (f *fulfillmentEd25519) MarshalBinary() ([]byte, error) {
	return append(Uint16Bytes(ED25519_ID), append(f.pub.Bytes(), f.sig.Bytes()...)...), nil
}

func (f *fulfillmentEd25519) UnmarshalBinary(p []byte) error {
	if id := MustUint16(p[:2]); id != ED25519_ID {
		return Errorf("Expected id=%d; got id=%d\n", ED25519_ID, id)
	}
	if size := len(p[2:]); size > ED25519_SIZE {
		return Error("Exceeds max payload size")
	}
	f.pub, _ = ed25519.NewPublicKey(p[2 : 2+ed25519.PUBKEY_SIZE])
	f.sig, _ = ed25519.NewSignature(p[2+ed25519.PUBKEY_SIZE:])
	return nil
}
