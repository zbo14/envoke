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
	FULFILLMENT_TYPE       = "fulfillment"

	// Types
	PREIMAGE_ID      = 0
	PREIMAGE_BITMASK = 0x03

	THRESHOLD_ID      = 2
	THRESHOLD_BITMASK = 0x09

	RSA_ID      = 3
	RSA_BITMASK = 0x11
	RSA_SIZE    = rsa.KEY_SIZE + rsa.SIGNATURE_SIZE

	ED25519_ID      = 4
	ED25519_BITMASK = 0x20
	ED25519_SIZE    = ed25519.PUBKEY_SIZE + ed25519.SIGNATURE_SIZE
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
	Validate([]byte) bool
	Weight() int
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
	return NewCondition(f.Bitmask(), f.Hash(), f.Id(), f.Size(), f.Weight())
}

// For bigchain txs..
// Get JSON-serializable details from non-condition fulfillment
func GetDetails(f Fulfillment) interface{} {
	if f == nil {
		return nil
	}
	if f.IsCondition() {
		panic("Cannot get details from condition")
	}
	id := f.Id()
	bitmask := f.Bitmask()
	pub := f.PublicKey()
	uri := GetCondition(f).String()
	switch {
	case
		id == PREIMAGE_ID && bitmask == PREIMAGE_BITMASK,
		id == ED25519_ID && bitmask == ED25519_BITMASK,
		id == RSA_ID && bitmask == RSA_BITMASK,
		id == THRESHOLD_ID && bitmask >= THRESHOLD_BITMASK:
		// Ok.. should we allow nil pubkey?
	default:
		Panicf("Unexpected id=%d, bitmask=%d\n", id, bitmask)
	}
	return struct {
		Details struct {
			Bitmask   int              `json:"bitmask"`
			PubKey    crypto.PublicKey `json:"public_key"`
			Signature crypto.Signature `json:"signature"`
			Type      string           `json:"type"`
			TypeId    int              `json:"type_id"`
		} `json:"details"`
		URI string `json:"uri"`
	}{
		Details: struct {
			Bitmask   int              `json:"bitmask"`
			PubKey    crypto.PublicKey `json:"public_key"`
			Signature crypto.Signature `json:"signature"`
			Type      string           `json:"type"`
			TypeId    int              `json:"type_id"`
		}{
			Bitmask:   bitmask,
			PubKey:    pub,
			Signature: nil,
			Type:      FULFILLMENT_TYPE,
			TypeId:    id,
		},
		URI: uri,
	}
}

func FulfillmentURI(p []byte) (string, error) {
	if size := len(p); size <= 2 {
		return "", Errorf("Expected fulfillment size > 2; got size=%d\n", size)
	}
	id := MustUint16(p[:2])
	payload64 := Base64UrlEncode(p[2:])
	return Sprintf("cf:%x:%s", id, payload64), nil
}

func ConditionURI(p []byte) (string, error) {
	if size := len(p); size != CONDITION_SIZE {
		return "", Errorf("Expected condition size=%d; got size=%d\n", CONDITION_SIZE, size)
	}
	id := MustUint16(p[:2])
	bitmask := MustUint32(p[2:6])
	hash64 := Base64UrlEncode(p[6 : 6+HASH_LENGTH])
	size := MustUint16(p[6+HASH_LENGTH : CONDITION_SIZE])
	return Sprintf("cc:%x:%x:%s:%d", id, bitmask, hash64, size), nil
}

func UnmarshalBinary(p []byte, weight int) (Fulfillment, error) {
	if uri, err := ConditionURI(p); err == nil {
		if MatchString(CONDITION_REGEX, uri) {
			c := NilCondition()
			c.weight = weight
			if err := c.UnmarshalBinary(p); err == nil {
				return c, nil
			} else {
				return nil, err
			}
		}
	}
	if uri, err := FulfillmentURI(p); err == nil {
		if MatchString(FULFILLMENT_REGEX, uri) {
			f := new(fulfillment)
			if err := f.UnmarshalBinary(p); err != nil {
				return nil, err
			}
			if weight < 1 {
				return nil, Error("Weight cannot be less than 1")
			}
			f.weight = weight
			switch f.id {
			case PREIMAGE_ID:
				return &fulfillmentPreImage{f}, nil
			case ED25519_ID:
				return &fulfillmentEd25519{f}, nil
			case RSA_ID:
				return &fulfillmentRSA{f}, nil
			case THRESHOLD_ID:
				subs, threshold, err := ThresholdSubs(f.payload)
				if err != nil {
					return nil, err
				}
				return &fulfillmentThreshold{
					fulfillment: f,
					subs:        subs,
					threshold:   threshold,
				}, nil
			default:
				return nil, Errorf("Unexpected id=%d\n", f.id)
			}
		}
	}
	return nil, Error("Could not match bytes to regex")
}

func UnmarshalURI(uri string) (f Fulfillment, err error) {
	if MatchString(CONDITION_REGEX, uri) {
		// Try to parse condition
		c := NilCondition()
		parts := Split(uri, ":")
		c.id, err = ParseUint16(parts[1])
		if err != nil {
			return nil, err
		}
		c.bitmask, err = ParseUint32(parts[2])
		if err != nil {
			return nil, err
		}
		switch {
		case
			c.id == PREIMAGE_ID && c.bitmask == PREIMAGE_BITMASK,
			c.id == ED25519_ID && c.bitmask == ED25519_BITMASK,
			c.id == RSA_ID && c.bitmask == RSA_BITMASK,
			c.id == THRESHOLD_ID && c.bitmask >= THRESHOLD_BITMASK:
			// Ok..
		default:
			return nil, Errorf("Unexpected id=%d, bitmask=%d\n", c.id, c.bitmask)
		}
		c.hash, err = Base64UrlDecode(parts[3])
		if err != nil {
			return nil, err
		}
		c.size, err = ParseUint16(parts[4])
		if err != nil {
			return nil, err
		}
		// TODO: check ed25519, RSA size
		if length := len(c.hash); length != HASH_LENGTH {
			return nil, Errorf("Expected hash with length=%d; got length=%d\n", HASH_LENGTH, length)
		}
		if c.size > MAX_PAYLOAD_SIZE {
			return nil, Error("Exceeded max payload size")
		}
		return c, nil
	}
	if MatchString(FULFILLMENT_REGEX, uri) {
		// Try to parse non-condition fulfillment
		_f := new(fulfillment)
		parts := Split(uri, ":")
		_f.id, err = ParseUint16(parts[1])
		if err != nil {
			return nil, err
		}
		_f.payload, err = Base64UrlDecode(parts[2])
		if err != nil {
			return nil, err
		}
		if size := len(_f.payload); size > MAX_PAYLOAD_SIZE {
			return nil, Error("Exceeds max payload size")
		}
		//TODO: check ed25519, RSA size
		switch _f.id {
		case PREIMAGE_ID:
			f = &fulfillmentPreImage{_f}
		case ED25519_ID:
			f = &fulfillmentEd25519{_f}
		case RSA_ID:
			f = &fulfillmentRSA{_f}
		case THRESHOLD_ID:
			f = &fulfillmentThreshold{
				fulfillment: _f,
			}
		default:
			return nil, Errorf("Unexpected id=%d\n", _f.id)
		}
		f.Init()
		return f, nil
	}
	return nil, Error("Could not match URI with regex")
}

type fulfillment struct {
	bitmask int
	hash    []byte
	id      int
	payload []byte
	size    int
	weight  int
}

func NewFulfillment(id int, payload []byte, weight int) *fulfillment {
	switch id {
	case PREIMAGE_ID, ED25519_ID, RSA_ID, THRESHOLD_ID:
		// Ok..
	default:
		Panicf("Unexpected id=%d\n", id)
	}
	if size := len(payload); size > MAX_PAYLOAD_SIZE {
		panic("Exceeded max payload size")
	}
	if weight < 1 {
		panic("Weight cannot be less than 1")
	}
	return &fulfillment{
		id:      id,
		payload: payload,
		weight:  weight,
	}
}

func (f *fulfillment) Bitmask() int { return f.bitmask }

func (f *fulfillment) FromString(uri string) error {
	if !MatchString(FULFILLMENT_REGEX, uri) {
		return Error("URI does not match fulfillment regex")
	}
	parts := Split(uri, ":")
	id, err := ParseUint16(parts[1])
	if err != nil {
		return err
	}
	switch id {
	case PREIMAGE_ID, ED25519_ID, RSA_ID, THRESHOLD_ID:
	default:
		return Errorf("Unexpected id=%d\n", id)
	}
	p, err := Base64UrlDecode(parts[2])
	if err != nil {
		return err
	}
	if size := len(p); size > MAX_PAYLOAD_SIZE {
		return Error("Exceeds max payload size")
	}
	//TODO: check ed25519, RSA size
	f.id = id
	f.payload = p
	return nil
}

func (f *fulfillment) Hash() []byte { return f.hash }

func (f *fulfillment) Id() int { return f.id }

func (c *Condition) Init() { /* no op */ }

func (f *fulfillment) IsCondition() bool { return false }

func (f *fulfillment) MarshalBinary() ([]byte, error) {
	return append(Uint16Bytes(f.id), f.payload...), nil
}

func (f *fulfillment) MarshalJSON() ([]byte, error) {
	if f == nil {
		return nil, nil
	}
	uri := f.String()
	json := MustMarshalJSON(uri)
	return json, nil
}

func (f *fulfillment) PublicKey() crypto.PublicKey { return nil }

func (f *fulfillment) Signature() crypto.Signature { return nil }

func (f *fulfillment) Size() int { return f.size }

func (f *fulfillment) String() string {
	payload64 := Base64UrlEncode(f.payload)
	return Sprintf("cf:%x:%s", f.id, payload64)
}

func (f *fulfillment) UnmarshalBinary(p []byte) error {
	id := MustUint16(p[:2])
	switch id {
	case PREIMAGE_ID, THRESHOLD_ID, RSA_ID, ED25519_ID:
	default:
		return Errorf("Unexpected id=%d\n", id)
	}
	if size := len(p[2:]); size > MAX_PAYLOAD_SIZE {
		return Error("Exceeds max payload size")
	}
	f.id = id
	f.payload = p[2:]
	return nil
}

func (f *fulfillment) Weight() int {
	return f.weight
}

// Condition

type Condition struct {
	*fulfillment
}

func NilCondition() *Condition {
	return &Condition{&fulfillment{}}
}

func NewConditionEd25519(pub *ed25519.PublicKey, weight int) *Condition {
	return NewCondition(ED25519_BITMASK, pub.Bytes(), ED25519_ID, ED25519_SIZE, weight)
}

func NewCondition(bitmask int, hash []byte, id int, size, weight int) *Condition {
	switch {
	case
		id == PREIMAGE_ID && bitmask == PREIMAGE_BITMASK,
		id == ED25519_ID && bitmask == ED25519_BITMASK,
		id == RSA_ID && bitmask == RSA_BITMASK,
		id == THRESHOLD_ID && bitmask >= THRESHOLD_BITMASK:
		// Ok..
	default:
		Panicf("Unexpected id=%d, bitmask=%d\n", id, bitmask)
	}
	if length := len(hash); length != HASH_LENGTH {
		Panicf("Expected hash length=%d; got length=%d\n", HASH_LENGTH, length)
	}
	if size > MAX_PAYLOAD_SIZE {
		panic("Exceeded max payload size")
	}
	if weight < 1 {
		panic("Weight cannot be less than 1")
	}
	return &Condition{
		&fulfillment{
			bitmask: bitmask,
			hash:    hash,
			id:      id,
			size:    size,
			weight:  weight,
		},
	}
}

func (c *Condition) FromString(uri string) error {
	if !MatchString(CONDITION_REGEX, uri) {
		return Error("URI does not match condition regex")
	}
	parts := Split(uri, ":")
	id, err := ParseUint16(parts[1])
	if err != nil {
		return err
	}
	bitmask, err := ParseUint32(parts[2])
	if err != nil {
		return err
	}
	switch {
	case
		id == PREIMAGE_ID && bitmask == PREIMAGE_BITMASK,
		id == ED25519_ID && bitmask == ED25519_BITMASK,
		id == RSA_ID && bitmask == RSA_BITMASK,
		id == THRESHOLD_ID && bitmask >= THRESHOLD_BITMASK:
		// Ok..
	default:
		return Errorf("Unexpected id=%d, bitmask=%d\n", id, bitmask)
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
	// TODO: check ed25519, RSA size
	case len(hash) != HASH_LENGTH:
		return Errorf("Expected hash with size=%d; got size=%d\n", HASH_LENGTH, len(hash))
	case size > MAX_PAYLOAD_SIZE:
		return Error("Exceeded max payload size")
	}
	c.bitmask = bitmask
	c.id = id
	c.hash = hash
	c.size = size
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

func (c *Condition) MarshalJSON() ([]byte, error) {
	if c == nil {
		return nil, nil
	}
	uri := c.String()
	json := MustMarshalJSON(uri)
	return json, nil

}

func (c *Condition) String() string {
	payload64 := Base64UrlEncode(c.hash)
	return Sprintf("cc:%x:%x:%s:%d", c.id, c.bitmask, payload64, c.size)
}

func (c *Condition) UnmarshalBinary(p []byte) error {
	if size := len(p); size != CONDITION_SIZE {
		return Errorf("Expected condition with size=%d; got size=%d\n", CONDITION_SIZE, size)
	}
	id := MustUint16(p[:2])
	bitmask := MustUint32(p[2:6])
	switch {
	case
		id == PREIMAGE_ID && bitmask == PREIMAGE_BITMASK,
		id == ED25519_ID && bitmask == ED25519_BITMASK,
		id == RSA_ID && bitmask == RSA_BITMASK,
		id == THRESHOLD_ID && bitmask >= THRESHOLD_BITMASK:
		// Ok..
	default:
		return Errorf("Unexpected id=%d, bitmask=%d\n", id, bitmask)
	}
	hash := p[6 : 6+HASH_LENGTH]
	size := MustUint16(p[6+HASH_LENGTH : CONDITION_SIZE])
	switch {
	// TODO: check ed25519, RSA size
	case
		len(hash) != HASH_LENGTH:
		return Errorf("Expected hash with size=%d; got size=%d\n", HASH_LENGTH, len(hash))
	case
		size > MAX_PAYLOAD_SIZE:
		return Error("Exceeded max payload size")
	}
	c.bitmask = bitmask
	c.hash = hash
	c.id = id
	c.size = size
	return nil
}

func (c *Condition) Validate(p []byte) bool {
	switch {
	case
		c.id == PREIMAGE_ID && c.bitmask == PREIMAGE_BITMASK,
		c.id == ED25519_ID && c.bitmask == ED25519_BITMASK,
		c.id == RSA_ID && c.bitmask == RSA_BITMASK,
		c.id == THRESHOLD_ID && c.bitmask >= THRESHOLD_BITMASK:
		// Ok..
	default:
		return false
	}
	switch {
	case
		len(c.hash) != HASH_LENGTH,
		c.size > MAX_PAYLOAD_SIZE,
		// TODO: check ed25519, RSA size
		c.weight < 1:
		return false
	}
	return true
}
