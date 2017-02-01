package conditions

import (
	"bytes"
	"github.com/zballs/envoke/crypto/ed25519"
	. "github.com/zballs/envoke/util"
	"sort"
)

// Sha256 Pre-Image

type fulfillmentPreImage struct {
	*fulfillment
}

func NewFulfillmentPreImage(preimage []byte, weight int) *fulfillmentPreImage {
	f := new(fulfillmentPreImage)
	f.fulfillment = NewFulfillment(PREIMAGE_ID, preimage, weight)
	f.Init()
	return f
}

func (f *fulfillmentPreImage) Init() {
	f.Bitmask()
	f.Hash()
	f.Size()
}

func (f *fulfillmentPreImage) Bitmask() int {
	if f.bitmask > 0 {
		return f.bitmask
	}
	f.bitmask = PREIMAGE_BITMASK
	return f.bitmask
}

func (f *fulfillmentPreImage) Hash() []byte {
	if f.hash != nil {
		return f.hash
	}
	f.hash = Checksum256(f.payload)
	return f.hash
}

func (f *fulfillmentPreImage) Validate(p []byte) bool { return true }

// Ed25519

type fulfillmentEd25519 struct {
	*fulfillment
}

func NewFulfillmentEd25519(msg []byte, priv *ed25519.PrivateKey, weight int) *fulfillmentEd25519 {
	f := new(fulfillmentEd25519)
	pub := priv.Public()
	sig := priv.Sign(msg)
	payload := append(pub.Bytes(), sig.Bytes()...)
	f.fulfillment = NewFulfillment(ED25519_ID, payload, weight)
	f.Init()
	return f
}

func (f *fulfillmentEd25519) Init() {
	f.Bitmask()
	f.Hash()
	f.Size()
}

func (f *fulfillmentEd25519) Bitmask() int {
	if f.bitmask > 0 {
		return f.bitmask
	}
	f.bitmask = ED25519_BITMASK
	return f.bitmask
}

func (f *fulfillmentEd25519) Hash() []byte {
	if f.hash != nil {
		return f.hash
	}
	f.hash = f.payload[:ed25519.PUBKEY_SIZE]
	return f.hash
}

func (f *fulfillmentEd25519) MarshalJSON() ([]byte, error) {
	if f == nil {
		return nil, nil
	}
	json := MustMarshalJSON(struct {
		Bitmask   int                `json:"bitmask"`
		PubKey    *ed25519.PublicKey `json:"public_key"`
		Signature string             `json:"signature,omitempty"`
		Type      string             `json:"type"`
		TypeId    int                `json:"type_id"`
	}{
		Bitmask:   f.Bitmask(),
		PubKey:    f.PubKey(),
		Signature: f.String(),
		Type:      FULFILLMENT_TYPE,
		TypeId:    f.Id(),
	})
	return json, nil
}

func (f *fulfillmentEd25519) PubKey() *ed25519.PublicKey {
	p := f.payload[:ed25519.PUBKEY_SIZE]
	pub, err := ed25519.NewPublicKey(p)
	Check(err)
	return pub
}

func (f *fulfillmentEd25519) Validate(p []byte) bool {
	pub, _ := ed25519.NewPublicKey(f.payload[:ed25519.PUBKEY_SIZE])
	sig, _ := ed25519.NewSignature(f.payload[ed25519.PUBKEY_SIZE:])
	return pub.Verify(p, sig)
}

// Sha256 Threshold

type fulfillmentThreshold struct {
	bitmask int
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
	f.fulfillment = NewFulfillment(THRESHOLD_ID, payload, weight)
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
	f.Bitmask()
	f.Hash()
	f.Size()
}

func (f *fulfillmentThreshold) Bitmask() int {
	if f.bitmask > 0 {
		return f.bitmask
	}
	f.bitmask = ThresholdBitmask(f.subs)
	return f.bitmask
}

func (f *fulfillmentThreshold) Hash() []byte {
	if f.hash != nil {
		return f.hash
	}
	f.hash = ThresholdHash(f.subs, f.threshold)
	return f.hash
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
	j = Pow2(numSubs)
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
			if sums[i] < sum {
				set = sets[i]
				sum = sums[i]
			}
		}
	}
FOR_LOOP:
	for _, sub := range subs {
		for _, s := range set {
			if sub == s {
				continue FOR_LOOP
			}
		}
		sub.Init()
		set = append(set, sub.Condition())
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

// Hash the subconditions
func ThresholdHash(subs Fulfillments, threshold int) []byte {
	hash := NewSha256()
	hash.Write(Uint32Bytes(threshold))
	hash.Write(UvarintBytes(len(subs)))
	for _, sub := range subs {
		weight := sub.Weight()
		hash.Write(UvarintBytes(weight))
		c := sub.Condition()
		p, _ := c.MarshalBinary()
		hash.Write(p)
	}
	return hash.Sum(nil)
}

func ThresholdSize(subs Fulfillments, threshold int) int {
	var i, j int
	numSubs := subs.Len()
	total := 4 + UvarintSize(numSubs) + numSubs
	j = Pow2(numSubs)
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
	subs := f.subs
	sort.Sort(subs)
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
	if valid < threshold {
		return false
	}
	return true
}
