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
	HASH_SIZE         = 32
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

func DefaultFulfillmentsFromPrivKeys(msgs [][]byte, privs []crypto.PrivateKey) Fulfillments {
	n := len(privs)
	if n != len(msgs) {
		panic(ErrorAppend(ErrInvalidSize, "slices are different sizes"))
	}
	fulfillments := make(Fulfillments, n)
	for i, priv := range privs {
		fulfillments[i] = DefaultFulfillmentFromPrivKey(msgs[i], priv)
	}
	return fulfillments
}

func FulfillmentsFromPrivKeys(msgs [][]byte, privs []crypto.PrivateKey, weights []int) Fulfillments {
	n := len(privs)
	if n != len(msgs) || n != len(weights) {
		panic(ErrorAppend(ErrInvalidSize, "slices are different sizes"))
	}
	fulfillments := make(Fulfillments, n)
	for i, priv := range privs {
		fulfillments[i] = FulfillmentFromPrivKey(msgs[i], priv, weights[i])
	}
	return fulfillments
}

func DefaultFulfillmentFromPrivKey(msg []byte, priv crypto.PrivateKey) Fulfillment {
	return FulfillmentFromPrivKey(msg, priv, 1)
}

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
	panic(ErrInvalidKey)
}

func DefaultFulfillmentsFromPubKeys(pubs []crypto.PublicKey) Fulfillments {
	fulfillments := make(Fulfillments, len(pubs))
	for i, pub := range pubs {
		fulfillments[i] = DefaultFulfillmentFromPubKey(pub)
	}
	return fulfillments
}

func FulfillmentsFromPubKeys(pubs []crypto.PublicKey, weights []int) Fulfillments {
	n := len(pubs)
	if n != len(weights) {
		panic(ErrorAppend(ErrInvalidSize, "slices are different sizes"))
	}
	fulfillments := make(Fulfillments, n)
	for i, pub := range pubs {
		fulfillments[i] = FulfillmentFromPubKey(pub, weights[i])
	}
	return fulfillments
}

func DefaultFulfillmentFromPubKey(pub crypto.PublicKey) Fulfillment {
	return FulfillmentFromPubKey(pub, 1)
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
	panic(ErrInvalidKey)
}

func FulfillWithPrivKey(f Fulfillment, msg []byte, priv crypto.PrivateKey) {
	sig := priv.Sign(msg)
	if !f.PublicKey().Verify(msg, sig) {
		panic(ErrInvalidSignature)
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
		panic(ErrInvalidSignature)
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
		return bytes.Compare(pi, pj) == -1
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
	buf := bytes.NewBuffer(p)
	id, err := ReadUint16(buf)
	if err != nil {
		return "", err
	}
	hash, err := ReadVarOctet(buf)
	if err != nil {
		return "", err
	}
	payload64 := Base64UrlEncode(hash)
	return Sprintf("cf:%x:%s", id, payload64), nil
}

func ConditionURI(p []byte) (string, error) {
	buf := bytes.NewBuffer(p)
	id, err := ReadUint16(buf)
	if err != nil {
		return "", err
	}
	bitmask, err := ReadVarUint(buf)
	if err != nil {
		return "", err
	}
	hash, err := ReadVarOctet(buf)
	if err != nil {
		return "", err
	}
	hash64 := Base64UrlEncode(hash)
	size, err := ReadVarUint(buf)
	if err != nil {
		return "", err
	}
	return Sprintf("cc:%x:%x:%s:%d", id, bitmask, hash64, size), nil
}

func UnmarshalBinary(p []byte, weight int) (f Fulfillment, err error) {
	c := NilCondition()
	if err := c.UnmarshalBinary(p); err == nil {
		c.weight = weight
		return c, nil
	}
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

func UnmarshalURI(uri string, weight int) (f Fulfillment, err error) {
	if MatchStr(CONDITION_REGEX, uri) {
		// Try to parse condition
		parts := SplitStr(uri, ":")
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
	if MatchStr(FULFILLMENT_REGEX, uri) {
		// Try to parse non-condition fulfillment
		ful := new(fulfillment)
		parts := SplitStr(uri, ":")
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
	return nil, ErrInvalidFulfillment
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
		//..
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
	if !MatchStr(FULFILLMENT_REGEX, uri) {
		return ErrInvalidFulfillment
	}
	parts := SplitStr(uri, ":")
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
	buf := new(bytes.Buffer)
	WriteUint16(buf, f.id)
	WriteVarOctet(buf, f.payload)
	return buf.Bytes(), nil
}

func (f *fulfillment) MarshalJSON() ([]byte, error) {
	if f == nil {
		return nil, nil
	}
	if !f.Validate(nil) {
		panic(ErrInvalidFulfillment)
	}
	return MustMarshalJSON(f.String()), nil
}

func (f *fulfillment) PublicKey() crypto.PublicKey { return nil }

func (f *fulfillment) Signature() crypto.Signature { return nil }

func (f *fulfillment) Size() int { return f.size }

func (f *fulfillment) String() string {
	payload64 := Base64UrlEncode(f.payload)
	return Sprintf("cf:%x:%s", f.id, payload64)
}

func (f *fulfillment) UnmarshalBinary(p []byte) (err error) {
	c := NilCondition()
	c.UnmarshalBinary(p)
	buf := bytes.NewBuffer(p)
	f.id, err = ReadUint16(buf)
	if err != nil {
		return err
	}
	f.payload, err = ReadVarOctet(buf)
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
		//..
	case f.id == ED25519_ID && f.bitmask == ED25519_BITMASK:
		if f.size != ED25519_SIZE {
			return false
		}
	case f.id == RSA_ID && f.bitmask == RSA_BITMASK:
		if f.size != RSA_SIZE {
			return false
		}
	default:
		return false
	}
	if f.size > MAX_PAYLOAD_SIZE {
		return false
	}
	if len(f.hash) != HASH_SIZE {
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
		panic(ErrInvalidCondition)
	}
	return c
}

func (c *Condition) FromString(uri string) (err error) {
	if !MatchStr(CONDITION_REGEX, uri) {
		return ErrInvalidCondition
	}
	parts := SplitStr(uri, ":")
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
	buf := new(bytes.Buffer)
	WriteUint16(buf, c.id)
	WriteVarUint(buf, c.bitmask)
	WriteVarOctet(buf, c.hash)
	WriteUint16(buf, c.size)
	return buf.Bytes(), nil
}

func (c *Condition) MarshalJSON() ([]byte, error) {
	if c == nil {
		return nil, nil
	}
	if !c.Validate(nil) {
		panic(ErrInvalidCondition)
	}
	return MustMarshalJSON(c.String()), nil
}

func (c *Condition) String() string {
	hash64 := Base64UrlEncode(c.hash)
	return Sprintf("cc:%x:%x:%s:%d", c.id, c.bitmask, hash64, c.size)
}

func (c *Condition) UnmarshalBinary(p []byte) (err error) {
	buf := bytes.NewBuffer(p)
	c.id, err = ReadUint16(buf)
	if err != nil {
		return err
	}
	c.bitmask, err = ReadVarUint(buf)
	if err != nil {
		return err
	}
	c.hash, err = ReadVarOctet(buf)
	if err != nil {
		return err
	}
	c.size, err = ReadUint16(buf)
	if err != nil {
		return err
	}
	if !c.Validate(nil) {
		return ErrInvalidCondition
	}
	return nil
}

func (c *Condition) UnmarshalJSON(p []byte) error {
	var uri string
	if err := UnmarshalJSON(p, &uri); err != nil {
		return err
	}
	return c.FromString(uri)
}
