package conditions

import (
	"bytes"
	"github.com/zballs/go_resonate/crypto/ed25519"
	. "github.com/zballs/go_resonate/util"
	"sort"
)

// crypto-conditions

const (
	// Params
	HASH_SIZE            = 32
	CONDITION_SIZE       = 10 + HASH_SIZE
	MAX_FULLFILMENT_SIZE = 0xfff
	SUPPORTED_BITMASK    = 0x3f

	// Regex
	CONDITION_REGEX        = `^cc:([1-9a-f][0-9a-f]{0,3}|0):[1-9a-f][0-9a-f]{0,15}:[a-zA-Z0-9_-]{0,86}:([1-9][0-9]{0,17}|0)$`
	CONDITION_REGEX_STRICT = `^cc:([1-9a-f][0-9a-f]{0,3}|0):[1-9a-f][0-9a-f]{0,7}:[a-zA-Z0-9_-]{0,86}:([1-9][0-9]{0,17}|0)$`
	FULFILLMENT_REGEX      = `^cf:([1-9a-f][0-9a-f]{0,3}|0):[a-zA-Z0-9_-]*$`

	// Types
	PREIMAGE_ID               uint16 = 0
	PREIMAGE_BITMASK          uint32 = 0x03
	PREIMAGE_FULFILLMENT_SIZE uint16 = MAX_FULLFILMENT_SIZE

	THRESHOLD_ID               uint16 = 2
	THRESHOLD_BITMASK          uint32 = 0x09
	THRESHOLD_FULFILLMENT_SIZE uint16 = MAX_FULLFILMENT_SIZE

	ED25519_ID               uint16 = 4
	ED25519_BITMASK          uint32 = 0x20
	ED25519_FULFILLMENT_SIZE uint16 = ed25519.PUBKEY_SIZE + ed25519.SIGNATURE_SIZE
)

// Fulfillment

type Fulfillment interface {
	Bitmask() uint32
	String() string
	FromString(string) error
	Hash() []byte
	Condition() *Condition
	Size() uint16
	Validate([]byte) bool
}

// Sha256 Preimage

type fulfillmentPreImage struct {
	preimage []byte
}

// PRNG
func NewFulfillmentPreImage(preimage []byte) *fulfillmentPreImage {
	if size := len(preimage); size > MAX_FULLFILMENT_SIZE {
		panic("PreImage is too big")
	}
	return &fulfillmentPreImage{preimage}
}

func (f *fulfillmentPreImage) Bitmask() uint32 { return PREIMAGE_BITMASK }

func (f *fulfillmentPreImage) Size() uint16 { return PREIMAGE_FULFILLMENT_SIZE }

func (f *fulfillmentPreImage) String() string {
	b64 := Base64UrlEncode(f.preimage)
	return Sprintf("cf:%x:%s", PREIMAGE_ID, b64)
}

func (f *fulfillmentPreImage) FromString(uri string) error {
	if !MatchString(FULFILLMENT_REGEX, uri) {
		return Error("Uri does not match fulfillment regex")
	}
	parts := Split(uri, ":")
	typeId, err := ParseUint16(parts[1])
	if err != nil {
		return err
	}
	if typeId != PREIMAGE_ID {
		return Errorf("Expected type_id=%d; got type_id=%d\n", PREIMAGE_ID, typeId)
	}
	bytes, err := Base64UrlDecode(parts[3])
	if err != nil {
		return err
	}
	if size := len(bytes); uint16(size) > PREIMAGE_FULFILLMENT_SIZE {
		return Errorf("Expected fulfillment_size <= %d; got fulfillment_size=%d\n", PREIMAGE_FULFILLMENT_SIZE, size)
	}
	f.preimage = bytes
	return nil
}

func (f *fulfillmentPreImage) Hash() []byte {
	return Checksum256(f.preimage)
}

func (f *fulfillmentPreImage) Condition() *Condition {
	return NewCondition(PREIMAGE_ID, PREIMAGE_BITMASK, f.Hash(), PREIMAGE_FULFILLMENT_SIZE)
}

func (f *fulfillmentPreImage) MarshalBinary() ([]byte, error) {
	return f.Condition().MarshalBinary()
}

func (f *fulfillmentPreImage) Validate(msg []byte) bool { return true }

// Threshold Sha256

type Sub interface {
	Bitmask() uint32
	Hash() []byte
	MarshalBinary() ([]byte, error)
}

type Subs []Sub

func (s Subs) Len() int {
	return len(s)
}

func (s Subs) Less(i, j int) bool {
	return bytes.Compare(s[i].Hash(), s[j].Hash()) == -1
}

func (s Subs) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type fulfillmentThreshold struct {
	submap    map[int]Subs
	threshold uint32
}

func NewFulfillmentThreshold(subs Subs, threshold uint32, weights []int) *fulfillmentThreshold {
	if len(subs) != len(weights) {
		panic("Number of subs must equal number of weights")
	}
	submap := make(map[int]Subs)
	for i, w := range weights {
		submap[w] = append(submap[w], subs[i])
	}
	return &fulfillmentThreshold{
		submap:    submap,
		threshold: threshold,
	}
}

func (f *fulfillmentThreshold) Bitmask() uint32 {
	bitmask := uint32(THRESHOLD_BITMASK)
	for _, subs := range f.submap {
		for _, sub := range subs {
			bitmask |= sub.Bitmask()
		}
	}
	return bitmask
}

func (f *fulfillmentThreshold) Hash() []byte {
	buf := new(bytes.Buffer)
	buf.Write(Uint32Bytes(f.threshold))
	weights := make([]int, len(f.submap))
	i, j := 0, 0
	for w, subs := range f.submap {
		weights[i] = w
		i++
		j += len(subs)
	}
	buf.Write(UvarintBytes(uint64(j)))
	sort.Ints(weights)
	for _, w := range weights {
		subs := f.submap[w]
		sort.Sort(subs)
		wbz := UvarintBytes(uint64(w))
		for _, sub := range subs {
			buf.Write(wbz)
			bin, _ := sub.MarshalBinary()
			buf.Write(bin)
		}
	}
	return buf.Bytes()
}

// Ed25519

type fulfillmentEd25519 struct {
	pub *ed25519.PublicKey
	sig *ed25519.Signature
}

func NewFulfillmentEd25519(msg []byte, priv *ed25519.PrivateKey) *fulfillmentEd25519 {
	pub := priv.Public()
	sig := priv.Sign(msg)
	return &fulfillmentEd25519{
		pub: pub,
		sig: sig,
	}
}

func (f *fulfillmentEd25519) Bitmask() uint32 { return ED25519_BITMASK }

func (f *fulfillmentEd25519) Size() uint16 { return ED25519_FULFILLMENT_SIZE }

func (f *fulfillmentEd25519) String() string {
	bytes := append(f.pub.Bytes(), f.sig.Bytes()...)
	b64 := Base64UrlEncode(bytes)
	return Sprintf("cf:%x:%s", ED25519_ID, b64)
}

func (f *fulfillmentEd25519) FromString(uri string) error {
	if !MatchString(FULFILLMENT_REGEX, uri) {
		return Error("Uri does not match fulfillment regex")
	}
	parts := Split(uri, ":")
	typeId, err := ParseUint16(parts[1])
	if err != nil {
		return err
	}
	if typeId != ED25519_ID {
		return Errorf("Expected type_id=%d; got type_id=%d\n", ED25519_ID, typeId)
	}
	bytes, err := Base64UrlDecode(parts[3])
	if err != nil {
		return err
	}
	if size := len(bytes); uint16(size) != ED25519_FULFILLMENT_SIZE {
		return Errorf("Expected fulfillment_size=%d; got fulfillment_size=%d\n", ED25519_FULFILLMENT_SIZE, size)
	}
	f.pub, _ = ed25519.NewPublicKey(bytes[:ed25519.PUBKEY_SIZE])
	f.sig, _ = ed25519.NewSignature(bytes[ed25519.PUBKEY_SIZE:ED25519_FULFILLMENT_SIZE])
	return nil
}

func (f *fulfillmentEd25519) Hash() []byte {
	return f.pub.Bytes()
}

func (f *fulfillmentEd25519) Condition() *Condition {
	return NewCondition(ED25519_ID, ED25519_BITMASK, f.Hash(), ED25519_FULFILLMENT_SIZE)
}

func (f *fulfillmentEd25519) Validate(msg []byte) bool {
	return f.pub.Verify(msg, f.sig)
}

// Condition

type Condition struct {
	typeId          uint16
	bitmask         uint32
	hash            []byte
	fulfillmentSize uint16
}

func NewConditionEd25519(pub *ed25519.PublicKey) *Condition {
	return NewCondition(ED25519_ID, ED25519_BITMASK, pub.Bytes(), ED25519_FULFILLMENT_SIZE)
}

func NewCondition(typeId uint16, bitmask uint32, hash []byte, fulfillmentSize uint16) *Condition {
	if size := len(hash); size != HASH_SIZE {
		Panicf("Expected hash with size=%d; got size=%d\n", HASH_SIZE, size)
	}
	return &Condition{
		typeId:          typeId,
		bitmask:         bitmask,
		hash:            hash,
		fulfillmentSize: fulfillmentSize,
	}
}

func (c *Condition) String() string {
	b64 := Base64UrlEncode(c.hash)
	return Sprintf("cc:%x:%x:%s:%d", c.typeId, c.bitmask, b64, c.fulfillmentSize)
}

func (c *Condition) FromString(uri string) error {
	if !MatchString(CONDITION_REGEX, uri) {
		return Error("Uri does not match condition regex")
	}
	parts := Split(uri, ":")
	typeId, err := ParseUint16(parts[1])
	if err != nil {
		return err
	}
	bitmask, err := ParseUint32(parts[2])
	if err != nil {
		return err
	}
	if bitmask&^SUPPORTED_BITMASK > 0 {
		return Error("Unsupported bitmask")
	}
	hash, err := Base64UrlDecode(parts[3])
	if err != nil {
		return err
	}
	fulfillmentSize, err := ParseUint16(parts[4])
	if err != nil {
		return err
	}
	c.typeId = typeId
	c.bitmask = bitmask
	c.hash = hash
	c.fulfillmentSize = fulfillmentSize
	return nil
}

func (c *Condition) TypeId() uint16 { return c.typeId }

func (c *Condition) Bitmask() uint32 { return c.bitmask }

func (c *Condition) Size() uint16 { return c.fulfillmentSize }

func (c *Condition) Hash() []byte { return c.hash }

func (c *Condition) Validate() bool {
	switch c.typeId {
	case ED25519_ID, PREIMAGE_ID, THRESHOLD_ID:
		return true
	default:
		return false
	}
}

func (c *Condition) MarshalBinary() ([]byte, error) {
	bytes := make([]byte, CONDITION_SIZE)
	copy(bytes[:2], Uint16Bytes(c.typeId))
	copy(bytes[2:6], Uint32Bytes(c.bitmask))
	copy(bytes[6:6+HASH_SIZE], c.hash)
	copy(bytes[6+HASH_SIZE:CONDITION_SIZE], Uint16Bytes(c.fulfillmentSize))
	return bytes, nil
}

func (c *Condition) UnmarshalBinary(bytes []byte) error {
	if size := len(bytes); size != CONDITION_SIZE {
		return Errorf("Expected condition with size=%d; got size=%d\n", CONDITION_SIZE, size)
	}
	c.typeId = MustUint16(bytes[:2])
	c.bitmask = MustUint32(bytes[2:6])
	c.hash = bytes[6 : 6+HASH_SIZE]
	c.fulfillmentSize = MustUint16(bytes[6+HASH_SIZE : CONDITION_SIZE])
	return nil
}

func (c *Condition) MarshalJSON() ([]byte, error) {
	uri := c.String()
	bytes := MustMarshalJSON(uri)
	return bytes, nil
}

func (c *Condition) UnmarshalJSON(bytes []byte) error {
	var uri string
	if err := UnmarshalJSON(bytes, &uri); err != nil {
		return err
	}
	if err := c.FromString(uri); err != nil {
		return err
	}
	return nil
}
