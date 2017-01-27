package conditions

import (
	. "github.com/zballs/go_resonate/util"
)

// Sha256 Pre-Image

type fulfillmentPreImage struct {
	preimage []byte
	weight   int
}

func NewFulfillmentPreImage(preimage []byte, weight int) *fulfillmentPreImage {
	if size := len(preimage); size > MAX_PAYLOAD_SIZE {
		panic("PreImage exceeds max size")
	}
	if weight < 1 {
		panic("Weight cannot be less than 1")
	}
	return &fulfillmentPreImage{
		preimage: preimage,
		weight:   weight,
	}
}

func (f *fulfillmentPreImage) Id() int { return PREIMAGE_ID }

func (f *fulfillmentPreImage) Bitmask() int { return PREIMAGE_BITMASK }

func (f *fulfillmentPreImage) Size() int { return 2 + len(f.preimage) }

func (f *fulfillmentPreImage) Len() int {
	bin, _ := f.Condition().MarshalBinary()
	return len(bin)
}

func (f *fulfillmentPreImage) Weight() int { return f.weight }

func (f *fulfillmentPreImage) String() string {
	b64 := Base64UrlEncode(f.preimage)
	return Sprintf("cf:%x:%s", PREIMAGE_ID, b64)
}

func (f *fulfillmentPreImage) FromString(uri string) error {
	if !MatchString(FULFILLMENT_REGEX, uri) {
		return Error("Uri does not match fulfillment regex")
	}
	parts := Split(uri, ":")
	id, err := ParseUint16(parts[1])
	if err != nil {
		return err
	}
	if id != PREIMAGE_ID {
		return Errorf("Expected id=%d; got id=%d\n", PREIMAGE_ID, id)
	}
	p, err := Base64UrlDecode(parts[2])
	if err != nil {
		return err
	}
	if size := len(p); size > MAX_PAYLOAD_SIZE {
		return Error("Exceeds max payload size")
	}
	f.preimage = p
	return nil
}

func (f *fulfillmentPreImage) Hash() []byte {
	return Checksum256(f.preimage)
}

func (f *fulfillmentPreImage) Condition() *Condition {
	return NewCondition(PREIMAGE_ID, PREIMAGE_BITMASK, f.Hash(), f.Size(), f.Weight())
}

func (f *fulfillmentPreImage) Validate(msg []byte) bool { return true }

func (f *fulfillmentPreImage) MarshalBinary() ([]byte, error) {
	return append(Uint16Bytes(PREIMAGE_ID), f.preimage...), nil

}

func (f *fulfillmentPreImage) UnmarshalBinary(p []byte) error {
	if id := MustUint16(p[:2]); id != PREIMAGE_ID {
		return Errorf("Expected id=%d; got id=%d\n", PREIMAGE_ID, id)
	}
	if size := len(p[2:]); size > MAX_PAYLOAD_SIZE {
		return Error("Exceeds max payload size")
	}
	f.preimage = p
	return nil
}
