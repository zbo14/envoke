package conditions

import (
	"bytes"
	"github.com/zballs/go_resonate/crypto/ed25519"
	. "github.com/zballs/go_resonate/util"
)

// crypto-conditions

const (
	// Params
	HASH_SIZE         = 32
	CONDITION_SIZE    = 10 + HASH_SIZE
	MAX_PAYLOAD_SIZE  = 0xfff
	SUPPORTED_BITMASK = 0x3f

	// Regex
	CONDITION_REGEX        = `^cc:([1-9a-f][0-9a-f]{0,3}|0):[1-9a-f][0-9a-f]{0,15}:[a-zA-Z0-9_-]{0,86}:([1-9][0-9]{0,17}|0)$`
	CONDITION_REGEX_STRICT = `^cc:([1-9a-f][0-9a-f]{0,3}|0):[1-9a-f][0-9a-f]{0,7}:[a-zA-Z0-9_-]{0,86}:([1-9][0-9]{0,17}|0)$`
	FULFILLMENT_REGEX      = `^cf:([1-9a-f][0-9a-f]{0,3}|0):[a-zA-Z0-9_-]*$`

	// Types
	PREIMAGE_ID      = 0
	PREIMAGE_BITMASK = 0x03

	THRESHOLD_ID      = 2
	THRESHOLD_BITMASK = 0x09

	ED25519_ID      = 4
	ED25519_BITMASK = 0x20
	ED25519_SIZE    = ed25519.PUBKEY_SIZE + ed25519.SIGNATURE_SIZE
)

// Fulfillment

type Fulfillment interface {
	Id() int
	Bitmask() int
	Size() int // size of serialized fulfillment
	Len() int  // length of serialized condition
	Weight() int
	String() string          // serialize fulfillment to uri
	FromString(string) error // parse fulfillment from uri
	Hash() []byte
	Condition() *Condition
	Validate([]byte) bool
	MarshalBinary() ([]byte, error)
	UnmarshalBinary([]byte) error
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

func UnmarshalFulfillment(p []byte, weight int) (Fulfillment, error) {
	c := new(Condition)
	if err := c.UnmarshalBinary(p); err == nil {
		return c, nil
	}
	switch id := MustUint16(p[:2]); id {
	case PREIMAGE_ID:
		f := new(fulfillmentPreImage)
		if err := f.UnmarshalBinary(p); err != nil {
			return nil, err
		}
		f.weight = weight
		return f, nil
	case THRESHOLD_ID:
		f := new(fulfillmentThreshold)
		if err := f.UnmarshalBinary(p); err != nil {
			return nil, err
		}
		f.weight = weight
		return f, nil
	case ED25519_ID:
		f := new(fulfillmentEd25519)
		if err := f.UnmarshalBinary(p); err != nil {
			return nil, err
		}
		f.weight = weight
		return f, nil
	default:
		return nil, Errorf("Unexpected id=%d\n", id)
	}
}

// Condition

type Condition struct {
	id      int
	bitmask int
	hash    []byte
	size    int
	weight  int
}

func NewConditionEd25519(pub *ed25519.PublicKey, weight int) *Condition {
	return NewCondition(ED25519_ID, ED25519_BITMASK, pub.Bytes(), ED25519_SIZE, weight)
}

func NewCondition(id, bitmask int, hash []byte, size, weight int) *Condition {
	if size := len(hash); size != HASH_SIZE {
		Panicf("Expected hash with size=%d; got size=%d\n", HASH_SIZE, size)
	}
	return &Condition{
		id:      id,
		bitmask: bitmask,
		hash:    hash,
		size:    size,
		weight:  weight,
	}
}

func (c *Condition) String() string {
	b64 := Base64UrlEncode(c.hash)
	return Sprintf("cc:%x:%x:%s:%d", c.id, c.bitmask, b64, c.size)
}

func (c *Condition) FromString(uri string) error {
	if !MatchString(CONDITION_REGEX, uri) {
		return Error("Uri does not match condition regex")
	}
	parts := Split(uri, ":")
	id, err := ParseUint16(parts[1])
	if err != nil {
		return err
	}
	switch int(id) {
	case PREIMAGE_ID, THRESHOLD_ID, ED25519_ID:
	default:
		return Errorf("Unexpected id=%d\n", id)
	}
	bitmask, err := ParseUint32(parts[2])
	if err != nil {
		return err
	}
	switch int(bitmask) {
	case PREIMAGE_BITMASK, THRESHOLD_BITMASK, ED25519_BITMASK:
	default:
		return Errorf("Unexpected bitmask=%d\n", bitmask)
	}
	hash, err := Base64UrlDecode(parts[3])
	if err != nil {
		return err
	}
	size, err := ParseUint16(parts[4])
	if err != nil {
		return err
	}
	switch {
	case int(id) == ED25519_ID && int(size) != ED25519_SIZE:
		return Errorf("Expected ed25519 size=%d; got size=%d\n", ED25519_SIZE, size)
	}
	c.id = int(id)
	c.bitmask = int(bitmask)
	c.hash = hash
	c.size = int(size)
	return nil
}

func (c *Condition) Id() int { return c.id }

func (c *Condition) Bitmask() int { return c.bitmask }

func (c *Condition) Size() int { return c.size }

func (c *Condition) Len() int { p, _ := c.MarshalBinary(); return len(p) }

func (c *Condition) Weight() int { return c.weight }

func (c *Condition) Hash() []byte { return c.hash }

func (c *Condition) Condition() *Condition { return c }

// TODO:
func (c *Condition) Validate(p []byte) bool { return true }

func (c *Condition) MarshalBinary() ([]byte, error) {
	p := make([]byte, CONDITION_SIZE)
	copy(p[:2], Uint16Bytes(c.id))
	copy(p[2:6], Uint32Bytes(c.bitmask))
	copy(p[6:6+HASH_SIZE], c.hash)
	copy(p[6+HASH_SIZE:CONDITION_SIZE], Uint16Bytes(c.size))
	return p, nil
}

func (c *Condition) UnmarshalBinary(p []byte) error {
	if size := len(p); size != CONDITION_SIZE {
		return Errorf("Expected condition with size=%d; got size=%d\n", CONDITION_SIZE, size)
	}
	id := MustUint16(p[:2])
	switch id {
	case PREIMAGE_ID, THRESHOLD_ID, ED25519_ID:
	default:
		return Errorf("Unexpected id=%d\n", id)
	}
	bitmask := MustUint32(p[2:6])
	switch bitmask {
	case PREIMAGE_BITMASK, THRESHOLD_BITMASK, ED25519_BITMASK:
	default:
		return Errorf("Unexpected bitmask=%d\n", bitmask)
	}
	hash := p[6 : 6+HASH_SIZE]
	size := MustUint16(p[6+HASH_SIZE : CONDITION_SIZE])
	if id == ED25519_ID && size != ED25519_SIZE {
		return Errorf("Expected ed25519 size=%d; got size=%d\n", ED25519_SIZE, size)
	}
	c.id = id
	c.bitmask = bitmask
	c.hash = hash
	c.size = size
	return nil
}

func (c *Condition) MarshalJSON() ([]byte, error) {
	uri := c.String()
	p := MustMarshalJSON(uri)
	return p, nil
}

func (c *Condition) UnmarshalJSON(p []byte) error {
	var uri string
	if err := UnmarshalJSON(p, &uri); err != nil {
		return err
	}
	if err := c.FromString(uri); err != nil {
		return err
	}
	return nil
}
