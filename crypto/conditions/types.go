package conditions

import (
	"bytes"
	. "github.com/zbo14/envoke/common"
	"github.com/zbo14/envoke/crypto/crypto"
	"github.com/zbo14/envoke/crypto/ed25519"
	"github.com/zbo14/envoke/crypto/rsa"
	"sort"
)

// SHA256 Pre-Image

type fulfillmentPreImage struct {
	*fulfillment
}

func NewFulfillmentPreImage(preimage []byte, weight int) *fulfillmentPreImage {
	f := new(fulfillmentPreImage)
	f.fulfillment = NewFulfillment(PREIMAGE_ID, f, preimage, weight)
	f.Init()
	return f
}

func (f *fulfillmentPreImage) Init() {
	f.bitmask = PREIMAGE_BITMASK
	f.hash = Checksum256(f.payload)
	f.size = len(f.payload)
}

// SHA256 Prefix

type fulfillmentPrefix struct {
	*fulfillment
	prefix []byte
	sub    Fulfillment
}

func NewFulfillmentPrefix(prefix []byte, sub Fulfillment, weight int) *fulfillmentPrefix {
	if sub.IsCondition() {
		panic("Expected non-condition fulfillment")
	}
	f := new(fulfillmentPrefix)
	p, _ := sub.MarshalBinary()
	payload := append(VarBytes(prefix), p...)
	f.fulfillment = NewFulfillment(PREFIX_ID, f, payload, weight)
	f.prefix = prefix
	f.sub = sub
	f.Init()
	return f
}

func (f *fulfillmentPrefix) Init() {
	if f.prefix == nil && f.sub == nil {
		buf := new(bytes.Buffer)
		buf.Write(f.payload)
		f.prefix = MustReadVarBytes(buf)
		var err error
		f.sub, err = UnmarshalBinary(buf.Bytes(), f.weight)
		Check(err)
		if f.sub.IsCondition() {
			panic("Expected non-condition fulfillment")
		}
	}
	if f.prefix != nil && f.sub != nil {
		f.bitmask = PREFIX_BITMASK
		p, _ := GetCondition(f.sub).MarshalBinary()
		f.hash = Checksum256(append(f.prefix, p...))
		f.size = len(f.payload)
		return
	}
	panic("Prefix and subfulfillment must both be set")
}

func (f *fulfillmentPrefix) Validate(p []byte) bool {
	if !f.fulfillment.Validate(nil) {
		return false
	}
	return f.sub.Validate(append(f.prefix, p...))
}

// ED25519

type fulfillmentEd25519 struct {
	*fulfillment
	pub *ed25519.PublicKey
	sig *ed25519.Signature
}

func NewFulfillmentEd25519(pub *ed25519.PublicKey, sig *ed25519.Signature, weight int) *fulfillmentEd25519 {
	f := new(fulfillmentEd25519)
	payload := append(pub.Bytes(), sig.Bytes()...)
	f.fulfillment = NewFulfillment(ED25519_ID, f, payload, weight)
	f.pub = pub
	f.sig = sig
	f.Init()
	return f
}

func (f *fulfillmentEd25519) Init() {
	if f.pub.Bytes() == nil {
		f.pub = new(ed25519.PublicKey)
		err := f.pub.FromBytes(f.payload[:ed25519.PUBKEY_SIZE])
		Check(err)
	}
	if f.sig.Bytes() == nil {
		f.sig = new(ed25519.Signature)
		f.sig.FromBytes(f.payload[ed25519.PUBKEY_SIZE:])
		// ignore err for now
	}
	f.bitmask = ED25519_BITMASK
	f.hash = f.pub.Bytes()
	f.size = ED25519_SIZE
}

func (f *fulfillmentEd25519) PublicKey() crypto.PublicKey {
	if f.pub.Bytes() == nil {
		return nil
	}
	return f.pub
}

func (f *fulfillmentEd25519) Signature() crypto.Signature {
	if f.sig.Bytes() == nil {
		return nil
	}
	return f.sig
}

func (f *fulfillmentEd25519) Validate(p []byte) bool {
	if !f.fulfillment.Validate(nil) {
		return false
	}
	return f.pub.Verify(p, f.sig)
}

// SHA256 RSA

type fulfillmentRSA struct {
	*fulfillment
	pub *rsa.PublicKey
	sig *rsa.Signature
}

func NewFulfillmentRSA(pub *rsa.PublicKey, sig *rsa.Signature, weight int) *fulfillmentRSA {
	f := new(fulfillmentRSA)
	payload := append(pub.Bytes(), sig.Bytes()...)
	f.fulfillment = NewFulfillment(RSA_ID, f, payload, weight)
	f.pub = pub
	f.sig = sig
	f.Init()
	return f
}

func (f *fulfillmentRSA) Init() {
	if f.pub.Bytes() == nil {
		f.pub = new(rsa.PublicKey)
		err := f.pub.FromBytes(f.payload[:rsa.KEY_SIZE])
		Check(err)
	}
	if f.sig.Bytes() == nil {
		f.sig = new(rsa.Signature)
		err := f.sig.FromBytes(f.payload[rsa.KEY_SIZE:])
		Check(err)
	}
	f.bitmask = RSA_BITMASK
	f.hash = Checksum256(f.pub.Bytes())
	f.size = RSA_SIZE
}

func (f *fulfillmentRSA) PublicKey() crypto.PublicKey {
	if f.pub.Bytes() == nil {
		return nil
	}
	return f.pub
}

func (f *fulfillmentRSA) Signature() crypto.Signature {
	if f.sig.Bytes() == nil {
		return nil
	}
	return f.sig
}

func (f *fulfillmentRSA) Validate(p []byte) bool {
	if !f.fulfillment.Validate(nil) {
		return false
	}
	return f.pub.Verify(p, f.sig)
}

// SHA256 Threshold

type fulfillmentThreshold struct {
	*fulfillment
	subs      Fulfillments
	threshold int
}

func NewFulfillmentThreshold(subs Fulfillments, threshold, weight int) *fulfillmentThreshold {
	if len(subs) == 0 {
		panic("Must have more than 0 subs")
	}
	if threshold <= 0 {
		panic("Threshold must be greater than 0")
	}
	sort.Sort(subs)
	payload := ThresholdPayload(subs, threshold)
	f := new(fulfillmentThreshold)
	f.fulfillment = NewFulfillment(THRESHOLD_ID, f, payload, weight)
	f.subs = subs
	f.threshold = threshold
	f.Init()
	return f
}

func (f *fulfillmentThreshold) Init() {
	if f.subs == nil && f.threshold == 0 {
		f.ThresholdSubs()
	}
	if f.subs != nil && f.threshold > 0 {
		//..
	} else {
		Panicf("Cannot have %d subs, threshold=%d\n", len(f.subs), f.threshold)
	}
	f.bitmask = ThresholdBitmask(f.subs)
	f.hash = ThresholdHash(f.subs, f.threshold)
	f.size = ThresholdSize(f.subs, f.threshold)
}

func ThresholdBitmask(subs Fulfillments) int {
	bitmask := THRESHOLD_BITMASK
	for _, sub := range subs {
		bitmask |= sub.Bitmask()
	}
	return bitmask
}

func ThresholdPayload(subs Fulfillments, threshold int) []byte {
	var i, j int
	numSubs := subs.Len()
	j = Exp2(numSubs)
	sums := make([]int, j)
	sets := make([]Fulfillments, j)
	thresholds := make([]int, j)
	for i, _ = range thresholds {
		thresholds[i] = threshold
	}
	for _, sub := range subs {
		j >>= 1
		with := true
		for i, _ = range sums {
			if thresholds[i] > 0 {
				if with {
					sums[i] += sub.Size()
					sets[i] = append(sets[i], sub)
					thresholds[i] -= sub.Weight()
				} else if !with {
					sums[i] += CONDITION_SIZE
				}
			}
			if (i+1)%j == 0 {
				with = !with
			}
		}
	}
	sum := 0
	var set Fulfillments
	for i, _ = range sets {
		if thresholds[i] <= 0 {
			if sums[i] < sum || sum == 0 {
				set = sets[i]
				sum = sums[i]
			}
		}
	}
OUTER:
	for _, sub := range subs {
		for _, s := range set {
			if sub == s {
				continue OUTER
			}
		}
		sub.Init()
		set = append(set, GetCondition(sub))
	}
	if set.Len() != numSubs {
		//..
	}
	buf := new(bytes.Buffer)
	buf.Write(UvarintBytes(threshold))
	buf.Write(UvarintBytes(numSubs))
	for _, sub := range set {
		buf.Write(UvarintBytes(sub.Weight()))
		p, _ := sub.MarshalBinary()
		MustWriteVarBytes(p, buf)
	}
	return buf.Bytes()
}

func (f *fulfillmentThreshold) ThresholdSubs() {
	if f.subs != nil && f.threshold > 0 {
		return
	}
	if f.subs == nil && f.threshold == 0 {
		var err error
		f.subs, f.threshold, err = ThresholdSubs(f.payload)
		Check(err)
		return
	}
	Panicf("Cannot have %d subs, threshold=%d\n", len(f.subs), f.threshold)
}

func ThresholdSubs(p []byte) (Fulfillments, int, error) {
	if p == nil {
		//..
	}
	buf := new(bytes.Buffer)
	buf.Write(p)
	threshold, err := ReadUvarint(buf)
	if err != nil {
		return nil, 0, err
	}
	numSubs, err := ReadUvarint(buf)
	if err != nil {
		return nil, 0, err
	}
	subs := make(Fulfillments, numSubs)
	for i := 0; i < numSubs; i++ {
		weight, err := ReadUvarint(buf)
		if err != nil {
			return nil, 0, err
		}
		p, err := ReadVarBytes(buf)
		if err != nil {
			return nil, 0, err
		}
		subs[i], err = UnmarshalBinary(p, weight)
		if err != nil {
			return nil, 0, err
		}
	}
	return subs, threshold, nil
}

// Sort subconditions then hash them..
func ThresholdHash(subs Fulfillments, threshold int) []byte {
	numSubs := len(subs)
	conds := make(Fulfillments, numSubs)
	for i, sub := range subs {
		sub.Init()
		conds[i] = GetCondition(sub)
	}
	sort.Sort(conds)
	hash := NewSha256()
	hash.Write(Uint32Bytes(threshold))
	hash.Write(UvarintBytes(numSubs))
	for _, c := range conds {
		hash.Write(UvarintBytes(c.Weight()))
		p, err := c.MarshalBinary()
		Check(err)
		hash.Write(p)
	}
	return hash.Sum(nil)
}

func ThresholdSize(subs Fulfillments, threshold int) int {
	var i, j int
	numSubs := subs.Len()
	total := 4 + UvarintSize(numSubs) + numSubs
	j = Exp2(numSubs)
	extras := make([]int, j)
	thresholds := make([]int, j)
	for i, _ = range thresholds {
		thresholds[i] = threshold
	}
	for _, sub := range subs {
		total += CONDITION_SIZE
		if weight := sub.Weight(); weight > 1 {
			total += UvarintSize(weight)
		}
		j >>= 1
		add := true
		extra := sub.Size() - CONDITION_SIZE
		for i, _ = range extras {
			if add && thresholds[i] > 0 {
				extras[i] += extra
				thresholds[i] -= sub.Weight()
			}
			if (i+1)%j == 0 {
				add = !add
			}
		}
	}
	extra := 0
	for i, _ = range extras {
		if thresholds[i] <= 0 {
			if extras[i] > extra {
				extra = extras[i]
			}
		}
	}
	if extra == 0 {
		panic("Insufficient subconditions/weights to meet threshold")
	}
	total += extra
	return total
}

func (f *fulfillmentThreshold) Validate(p []byte) bool {
	if !f.fulfillment.Validate(nil) {
		return false
	}
	subs := f.subs
	threshold := f.threshold
	min, total := 0, 0
	var subf Fulfillments
	for _, sub := range subs {
		if !sub.IsCondition() {
			subf = append(subf, sub)
			weight := sub.Weight()
			if weight < min || min == 0 {
				min = weight
			}
			total += min
		}
	}
	if total < threshold {
		return false
	}
	valid := 0
	buf := new(bytes.Buffer)
	buf.Write(p)
	for _, f := range subf {
		p, err := ReadVarBytes(buf)
		if err != nil {
			return false
		}
		if f.Validate(p) {
			valid += f.Weight()
		}
	}
	return valid >= threshold
}

// SHA256 Timeout
type fulfillmentTimeout struct {
	expires int64
	*fulfillment
}

func NewFulfillmentTimeout(expires int64, weight int) *fulfillmentTimeout {
	f := new(fulfillmentTimeout)
	payload := TimestampBytes(expires)
	f.fulfillment = NewFulfillment(TIMEOUT_ID, f, payload, weight)
	f.expires = expires
	f.Init()
	return f
}

func (f *fulfillmentTimeout) Init() {
	if f.expires == 0 {
		f.expires = TimestampFromBytes(f.payload)
	}
	f.bitmask = TIMEOUT_BITMASK
	f.hash = Checksum256(f.payload)
	f.size = len(f.payload)
}

func (f *fulfillmentTimeout) Validate(p []byte) bool {
	if !f.fulfillment.Validate(nil) {
		return false
	}
	timestamp := TimestampFromBytes(p)
	return timestamp <= f.expires
}
