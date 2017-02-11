package conditions

import (
	"bytes"
	. "github.com/zbo14/envoke/common"
	"github.com/zbo14/envoke/crypto/crypto"
	"github.com/zbo14/envoke/crypto/ed25519"
	"github.com/zbo14/envoke/crypto/rsa"
)

// ILP crypto-conditions

const (
	// Params
	HASH_LENGTH       = 32
	CONDITION_SIZE    = 10 + HASH_LENGTH
	MAX_PAYLOAD_SIZE  = 0xfff
	SUPPORTED_BITMASK = 0x3f

	// Regex
	CONDITION_REGEX        = `^cc:([1-9a-f][0-9a-f]{0,3}|0):[1-9a-f][0-9a-f]{0,15}:[a-zA-Z0-9_-]{0,86}:([1-9][0-9]{0,17}|0)$`
	CONDITION_REGEX_STRICT = `^cc:([1-9a-f][0-9a-f]{0,3}|0):[1-9a-f][0-9a-f]{0,7}:[a-zA-Z0-9_-]{0,86}:([1-9][0-9]{0,17}|0)$`
	FULFILLMENT_REGEX      = `^cf:([1-9a-f][0-9a-f]{0,3}|0):[a-zA-Z0-9_-]*$`
	TIMESTAMP_REGEX        = `^\d{10}(\.\d+)?$`

	// Types
	FULFILLMENT_TYPE = "fulfillment"

	PREIMAGE_ID      = 0
	PREIMAGE_BITMASK = 0x03

	PREFIX_ID      = 1
	PREFIX_BITMASK = 0x05

	THRESHOLD_ID      = 2
	THRESHOLD_BITMASK = 0x09

	RSA_ID      = 3
	RSA_BITMASK = 0x11
	RSA_SIZE    = rsa.KEY_SIZE + rsa.SIGNATURE_SIZE

	ED25519_ID      = 4
	ED25519_BITMASK = 0x20
	ED25519_SIZE    = ed25519.PUBKEY_SIZE + ed25519.SIGNATURE_SIZE

	TIMEOUT_ID      = 5
	TIMEOUT_BITMASK = 0x00
)

// Fulfillment

type Fulfillment interface {
	Bitmask() int
	FromString(string) error
	Hash() []byte
	Id() int
	Init()
	IsCondition() bool
	MarshalBinary() ([]byte, error)
	MarshalJSON() ([]byte, error)
	PublicKey() crypto.PublicKey
	Signature() crypto.Signature
	Size() int
	String() string
	UnmarshalBinary([]byte) error
	UnmarshalJSON([]byte) error
	Validate([]byte) bool
	Weight() int
}

// Fufillment from key

func FulfillmentFromPrivKey(msg []byte, priv crypto.PrivateKey, weight int) Fulfillment {
	switch priv.(type) {
	case *ed25519.PrivateKey:
		privEd25519 := priv.(*ed25519.PrivateKey)
		pubEd25519 := privEd25519.Public().(*ed25519.PublicKey)
		sigEd25519 := privEd25519.Sign(msg).(*ed25519.Signature)
		return NewFulfillmentEd25519(pubEd25519, sigEd25519, weight)
	case *rsa.PrivateKey:
		privRSA := priv.(*rsa.PrivateKey)
		pubRSA := privRSA.Public().(*rsa.PublicKey)
		sigRSA := privRSA.Sign(msg).(*rsa.Signature)
		return NewFulfillmentRSA(pubRSA, sigRSA, weight)
	}
	panic(ErrInvalidKey.Error())
}

func FulfillmentFromPubKey(pub crypto.PublicKey, weight int) Fulfillment {
	switch pub.(type) {
	case *ed25519.PublicKey:
		pubEd25519 := pub.(*ed25519.PublicKey)
		return NewFulfillmentEd25519(pubEd25519, nil, weight)
	case *rsa.PublicKey:
		pubRSA := pub.(*rsa.PublicKey)
		return NewFulfillmentRSA(pubRSA, nil, weight)
	}
	panic(ErrInvalidKey.Error())
}

func FulfillWithPrivKey(f Fulfillment, msg []byte, priv crypto.PrivateKey) {
	sig := priv.Sign(msg)
	if !f.PublicKey().Verify(msg, sig) {
		panic(ErrInvalidSignature.Error())
	}
	switch sig.(type) {
	case *ed25519.Signature:
		ful := f.(*fulfillmentEd25519)
		ful.payload = append(ful.payload, sig.Bytes()...)
		ful.sig = sig.(*ed25519.Signature)
	case *rsa.Signature:
		ful := f.(*fulfillmentRSA)
		ful.payload = append(ful.payload, sig.Bytes()...)
		ful.sig = sig.(*rsa.Signature)
	default:
		panic(ErrInvalidSignature.Error())
	}
}

type Fulfillments []Fulfillment

func (fs Fulfillments) Len() int {
	return len(fs)
}

// sort in `descending` order by weights, then lexicographically
func (fs Fulfillments) Less(i, j int) bool {
	if fs[i].Weight() > fs[j].Weight() {
		return true
	}
	if fs[i].Weight() == fs[j].Weight() {
		pi, _ := fs[i].MarshalBinary()
		pj, _ := fs[j].MarshalBinary()
		return bytes.Compare(pi, pj) == 1
	}
	return false
}

func (fs Fulfillments) Swap(i, j int) {
	fs[i], fs[j] = fs[j], fs[i]
}

func GetCondition(f Fulfillment) *Condition {
	if f.IsCondition() {
		return f.(*Condition)
	}
	return NewCondition(f.Bitmask(), f.Hash(), f.Id(), f.PublicKey(), f.Size(), f.Weight())
}

func FulfillmentURI(p []byte) (string, error) {
	if len(p) <= 2 {
		return "", ErrInvalidSize
	}
	id := MustUint16(p[:2])
	payload64 := Base64UrlEncode(p[2:])
	return Sprintf("cf:%x:%s", id, payload64), nil
}

func ConditionURI(p []byte) (string, error) {
	if len(p) != CONDITION_SIZE {
		return "", ErrInvalidSize
	}
	id := MustUint16(p[:2])
	bitmask := MustUint32(p[2:6])
	hash64 := Base64UrlEncode(p[6 : 6+HASH_LENGTH])
	size := MustUint16(p[6+HASH_LENGTH : CONDITION_SIZE])
	return Sprintf("cc:%x:%x:%s:%d", id, bitmask, hash64, size), nil
}

func UnmarshalBinary(p []byte, weight int) (f Fulfillment, err error) {
	if uri, err := ConditionURI(p); err == nil {
		if MatchString(CONDITION_REGEX, uri) {
			c := NilCondition()
			c.weight = weight
			if err := c.UnmarshalBinary(p); err != nil {
				return nil, err
			}
			return c, nil
		}
	}
	if uri, err := FulfillmentURI(p); err == nil {
		if MatchString(FULFILLMENT_REGEX, uri) {
			ful := new(fulfillment)
			if err := ful.UnmarshalBinary(p); err != nil {
				return nil, err
			}
			ful.weight = weight
			switch ful.id {
			case PREIMAGE_ID:
				f = &fulfillmentPreImage{ful}
			case PREFIX_ID:
				f = &fulfillmentPrefix{
					fulfillment: ful,
				}
			case ED25519_ID:
				f = &fulfillmentEd25519{
					fulfillment: ful,
				}
			case RSA_ID:
				f = &fulfillmentRSA{
					fulfillment: ful,
				}
			case THRESHOLD_ID:
				f = &fulfillmentThreshold{
					fulfillment: ful,
				}
			}
			f.Init()
			if !ful.Validate(nil) {
				return nil, ErrInvalidFulfillment
			}
			return f, nil
		}
	}
	return nil, ErrInvalidRegex
}

func UnmarshalURI(uri string, weight int) (f Fulfillment, err error) {
	if MatchString(CONDITION_REGEX, uri) {
		// Try to parse condition
		parts := Split(uri, ":")
		c := NilCondition()
		c.id, err = ParseUint16(parts[1], 16)
		if err != nil {
			return nil, err
		}
		c.bitmask, err = ParseUint32(parts[2], 16)
		if err != nil {
			return nil, err
		}
		c.hash, err = Base64UrlDecode(parts[3])
		if err != nil {
			return nil, err
		}
		c.size, err = ParseUint16(parts[4], 10)
		if err != nil {
			return nil, err
		}
		c.weight = weight
		if !c.Validate(nil) {
			return nil, ErrInvalidCondition
		}
		return c, nil
	}
	if MatchString(FULFILLMENT_REGEX, uri) {
		// Try to parse non-condition fulfillment
		ful := new(fulfillment)
		parts := Split(uri, ":")
		ful.id, err = ParseUint16(parts[1], 16)
		if err != nil {
			return nil, err
		}
		ful.payload, err = Base64UrlDecode(parts[2])
		if err != nil {
			return nil, err
		}
		ful.weight = weight
		switch ful.id {
		case PREIMAGE_ID:
			f = &fulfillmentPreImage{ful}
		case PREFIX_ID:
			f = &fulfillmentPrefix{
				fulfillment: ful,
			}
		case ED25519_ID:
			f = &fulfillmentEd25519{
				fulfillment: ful,
			}
		case RSA_ID:
			f = &fulfillmentRSA{
				fulfillment: ful,
			}
		case THRESHOLD_ID:
			f = &fulfillmentThreshold{
				fulfillment: ful,
			}
		}
		f.Init()
		if !ful.Validate(nil) {
			return nil, ErrInvalidFulfillment
		}
		return f, nil
	}
	return nil, ErrInvalidRegex
}

type fulfillment struct {
	bitmask int
	hash    []byte
	id      int
	outer   Fulfillment
	payload []byte
	size    int
	weight  int
}

func NewFulfillment(id int, outer Fulfillment, payload []byte, weight int) *fulfillment {
	switch id {
	case PREIMAGE_ID, PREFIX_ID, ED25519_ID, RSA_ID, THRESHOLD_ID:
		// Ok..
	default:
		Panicf("Unexpected id=%d\n", id)
	}
	if len(payload) > MAX_PAYLOAD_SIZE {
		panic("Exceeds max payload size")
	}
	if weight < 1 {
		panic("Weight cannot be less than 1")
	}
	return &fulfillment{
		id:      id,
		outer:   outer,
		payload: payload,
		weight:  weight,
	}
}

func (f *fulfillment) Bitmask() int { return f.bitmask }

func (f *fulfillment) FromString(uri string) (err error) {
	if !MatchString(FULFILLMENT_REGEX, uri) {
		return ErrInvalidRegex
	}
	parts := Split(uri, ":")
	f.id, err = ParseUint16(parts[1], 16)
	if err != nil {
		return err
	}
	f.payload, err = Base64UrlDecode(parts[2])
	if err != nil {
		return err
	}
	if f.outer != nil {
		f.outer.Init()
		if !f.Validate(nil) {
			return ErrInvalidFulfillment
		}
	}
	return nil
}

func (f *fulfillment) Hash() []byte { return f.hash }

func (f *fulfillment) Id() int { return f.id }

func (f *fulfillment) Init() { /* no op */ }

func (f *fulfillment) IsCondition() bool { return false }

func (f *fulfillment) MarshalBinary() ([]byte, error) {
	return append(Uint16Bytes(f.id), f.payload...), nil
}

func (f *fulfillment) MarshalJSON() ([]byte, error) {
	if f == nil {
		return nil, nil
	}
	if !f.Validate(nil) {
		panic(ErrInvalidFulfillment.Error())
	}
	/*
		if f.outer.PublicKey() != nil {
			if f.outer.Signature() == nil {
				return MustMarshalJSON(struct {
					Bitmask   int              `json:"bitmask"`
					PubKey    crypto.PublicKey `json:"public_key"`
					Signature crypto.Signature `json:"signature"`
					Type      string           `json:"type"`
					TypeId    int              `json:"type_id"`
				}{
					Bitmask:   f.bitmask,
					PubKey:    f.outer.PublicKey(),
					Signature: nil,
					Type:      FULFILLMENT_TYPE,
					TypeId:    f.id,
				}), nil
			}
		}
	*/
	return MustMarshalJSON(f.String()), nil
}

func (f *fulfillment) PublicKey() crypto.PublicKey { return nil }

func (f *fulfillment) Signature() crypto.Signature { return nil }

func (f *fulfillment) Size() int { return f.size }

func (f *fulfillment) String() string {
	payload64 := Base64UrlEncode(f.payload)
	return Sprintf("cf:%x:%s", f.id, payload64)
}

func (f *fulfillment) UnmarshalBinary(p []byte) error {
	if len(p) <= 2 {
		return ErrInvalidSize
	}
	f.id = MustUint16(p[:2])
	f.payload = p[2:]
	if f.outer != nil {
		f.outer.Init()
		if !f.Validate(nil) {
			return ErrInvalidFulfillment
		}
	}
	return nil
}

func (f *fulfillment) UnmarshalJSON(p []byte) error {
	var uri string
	if err := UnmarshalJSON(p, &uri); err != nil {
		return err
	}
	if err := f.FromString(uri); err != nil {
		return err
	}
	return nil
}

func (f *fulfillment) Validate(p []byte) bool {
	switch {
	case
		f.id == PREIMAGE_ID && f.bitmask == PREIMAGE_BITMASK,
		f.id == PREFIX_ID && f.bitmask == PREFIX_BITMASK,
		f.id == THRESHOLD_ID && f.bitmask >= THRESHOLD_BITMASK:
		// Ok..
	case f.id == ED25519_ID && f.bitmask == ED25519_BITMASK:
		if f.size != ED25519_SIZE {
			Println(1)
			return false
		}
	case f.id == RSA_ID && f.bitmask == RSA_BITMASK:
		if f.size != RSA_SIZE {
			Println(2)
			return false
		}
	default:
		Println(3)
		return false
	}
	switch {
	case
		len(f.hash) != HASH_LENGTH,
		f.size > MAX_PAYLOAD_SIZE:
		// f.weight < 1
		return false
	}
	return true
}

func (f *fulfillment) Weight() int {
	return f.weight
}

// Condition

type Condition struct {
	*fulfillment
	pub crypto.PublicKey
}

func NilCondition() *Condition {
	return &Condition{
		fulfillment: &fulfillment{},
	}
}

func NewConditionWithPubKey(pub crypto.PublicKey, weight int) *Condition {
	switch pub.(type) {
	case *ed25519.PublicKey:
		return NewCondition(
			ED25519_BITMASK,
			pub.Bytes(),
			ED25519_ID,
			pub,
			ED25519_SIZE,
			weight)
	case *rsa.PublicKey:
		return NewCondition(
			RSA_BITMASK,
			Checksum256(pub.Bytes()),
			RSA_ID,
			pub,
			RSA_SIZE,
			weight)
	}
	panic(ErrInvalidKey.Error())
}

func NewCondition(bitmask int, hash []byte, id int, pub crypto.PublicKey, size, weight int) *Condition {
	c := &Condition{
		&fulfillment{
			bitmask: bitmask,
			hash:    hash,
			id:      id,
			size:    size,
			weight:  weight,
		}, pub,
	}
	if !c.Validate(nil) {
		panic(ErrInvalidCondition.Error())
	}
	return c
}

func (c *Condition) FromString(uri string) (err error) {
	if !MatchString(CONDITION_REGEX, uri) {
		return ErrInvalidRegex
	}
	parts := Split(uri, ":")
	c.id, err = ParseUint16(parts[1], 16)
	if err != nil {
		return err
	}
	c.bitmask, err = ParseUint32(parts[2], 16)
	if err != nil {
		return err
	}
	c.hash, err = Base64UrlDecode(parts[3])
	if err != nil {
		return err
	}
	c.size, err = ParseUint16(parts[4], 10)
	if err != nil {
		return err
	}
	if !c.Validate(nil) {
		return ErrInvalidCondition
	}
	return nil
}

func (c *Condition) IsCondition() bool { return true }

func (c *Condition) MarshalBinary() ([]byte, error) {
	p := make([]byte, CONDITION_SIZE)
	copy(p[:2], Uint16Bytes(c.id))
	copy(p[2:6], Uint32Bytes(c.bitmask))
	copy(p[6:6+HASH_LENGTH], c.hash)
	copy(p[6+HASH_LENGTH:CONDITION_SIZE], Uint16Bytes(c.size))
	return p, nil
}

// For bigchain txs..
// Get JSON-serializable details of condition
func (c *Condition) MarshalJSON() ([]byte, error) {
	if c == nil {
		return nil, nil
	}
	if !c.Validate(nil) {
		panic(ErrInvalidCondition.Error())
	}
	return MustMarshalJSON(struct {
		Details struct {
			Bitmask   int              `json:"bitmask"`
			PubKey    crypto.PublicKey `json:"public_key"`
			Signature interface{}      `json:"signature"`
			Type      string           `json:"type"`
			TypeId    int              `json:"type_id"`
		} `json:"details"`
		URI string `json:"uri"`
	}{
		Details: struct {
			Bitmask   int              `json:"bitmask"`
			PubKey    crypto.PublicKey `json:"public_key"`
			Signature interface{}      `json:"signature"`
			Type      string           `json:"type"`
			TypeId    int              `json:"type_id"`
		}{
			Bitmask:   c.bitmask,
			PubKey:    c.pub,
			Signature: nil,
			Type:      FULFILLMENT_TYPE,
			TypeId:    c.id,
		},
		URI: c.String(),
	}), nil
}

func (c *Condition) String() string {
	hash64 := Base64UrlEncode(c.hash)
	return Sprintf("cc:%x:%x:%s:%d", c.id, c.bitmask, hash64, c.size)
}

func (c *Condition) UnmarshalBinary(p []byte) error {
	if len(p) != CONDITION_SIZE {
		return ErrInvalidSize
	}
	c.id = MustUint16(p[:2])
	c.bitmask = MustUint32(p[2:6])
	c.hash = p[6 : 6+HASH_LENGTH]
	c.size = MustUint16(p[6+HASH_LENGTH:])
	if !c.Validate(nil) {
		return ErrInvalidCondition
	}
	return nil
}

func (c *Condition) UnmarshalJSON(p []byte) error {
	v := struct {
		Details struct {
			Bitmask   int    `json:"bitmask"`
			PubKey    string `json:"public_key"`
			Signature string `json:"signature"`
			Type      string `json:"type"`
			TypeId    int    `json:"type_id"`
		} `json:"details"`
		URI string `json:"uri"`
	}{}
	if err := UnmarshalJSON(p, &v); err != nil {
		return err
	}
	if err := c.FromString(v.URI); err != nil {
		return err
	}
	return nil
}
