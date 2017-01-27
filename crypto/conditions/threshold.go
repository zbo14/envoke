package conditions

import (
	"bytes"
	. "github.com/zballs/go_resonate/util"
	"sort"
)

// Threshold Sha256

type fulfillmentThreshold struct {
	subs      Fulfillments
	threshold int
	weight    int
}

func NewFulfillmentThreshold(subs Fulfillments, threshold, weight int) *fulfillmentThreshold {
	if len(subs) == 0 {
		panic("Must have more than 0 subs")
	}
	if threshold == 0 {
		panic("Threshold must be greater than 0")
	}
	if weight < 1 {
		panic("Weight cannot be less than 1")
	}
	return &fulfillmentThreshold{
		subs:      subs,
		threshold: threshold,
		weight:    weight,
	}
}

func (f *fulfillmentThreshold) Id() int { return THRESHOLD_ID }

func (f *fulfillmentThreshold) Weight() int { return f.weight }

func (f *fulfillmentThreshold) Bitmask() int {
	bitmask := THRESHOLD_BITMASK
	for _, sub := range f.subs {
		bitmask |= sub.Bitmask()
	}
	return bitmask
}

func (f *fulfillmentThreshold) MarshalBinary() ([]byte, error) {
	if f.threshold <= 0 {
		panic("Threshold must be greater than 0")
	}
	var i, j int
	subs := f.subs
	sort.Sort(subs)
	numSubs := subs.Len()
	j = Pow2(numSubs)
	sums := make([]int, j)
	var sets []Fulfillments
	thresholds := make([]int, j)
	for i, _ = range thresholds {
		thresholds[i] = f.threshold
	}
	for _, sub := range subs {
		j >>= 1
		with := true
		for i, _ = range sums {
			if (i+1)%j == 0 {
				with = !with
			}
			if thresholds[i] > 0 {
				if with {
					sums[i] += f.Size()
					sets[i] = append(sets[i], sub)
					thresholds[i] -= sub.Weight()
				} else if !with {
					sums[i] += f.Len()
				}
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
		set = append(set, sub.Condition())
	}
	if set.Len() != numSubs {
		//..
	}
	buf := new(bytes.Buffer)
	buf.Write(Uint16Bytes(THRESHOLD_ID))
	buf.Write(UvarintBytes(f.threshold))
	buf.Write(UvarintBytes(numSubs))
	for _, sub := range set {
		buf.Write(UvarintBytes(sub.Weight()))
		p, _ := sub.MarshalBinary()
		buf.Write(UvarintBytes(len(p)))
		buf.Write(p)
	}
	return buf.Bytes(), nil
}

func (f *fulfillmentThreshold) UnmarshalBinary(p []byte) error {
	buf := new(bytes.Buffer)
	buf.Write(p)
	bz, err := ReadN(buf, 2)
	if err != nil {
		return err
	}
	if id := MustUint16(bz); id != THRESHOLD_ID {
		return Errorf("Expected id=%d; got id=%d\n", THRESHOLD_ID, id)
	}
	threshold, err := ReadUvarint(buf)
	if err != nil {
		return err
	}
	numSubs, err := ReadUvarint(buf)
	if err != nil {
		return err
	}
	subs := make(Fulfillments, numSubs)
	for i := 0; i < numSubs; i++ {
		weight, err := ReadUvarint(buf)
		if err != nil {
			return err
		}
		n, err := ReadUvarint(buf)
		if err != nil {
			return err
		}
		p, err := ReadN(buf, n)
		if err != nil {
			return err
		}
		subs[i], err = UnmarshalFulfillment(p, weight)
		if err != nil {
			return err
		}
	}
	// should we check if subs are sorted?
	f.subs = subs
	f.threshold = threshold
	return nil
}

func (f *fulfillmentThreshold) Hash() []byte {
	subs := f.subs
	sort.Sort(subs)
	threshold := f.threshold
	hash := NewSha256()
	hash.Write(Uint32Bytes(threshold))
	hash.Write(UvarintBytes(len(subs)))
	for _, sub := range subs {
		w := sub.Weight()
		hash.Write(UvarintBytes(w))
		p, _ := sub.MarshalBinary()
		hash.Write(p)
	}
	return hash.Sum(nil)
}

func (f *fulfillmentThreshold) Size() int {
	var i, j int
	subs := f.subs
	sort.Sort(subs)
	numSubs := subs.Len()
	threshold := f.threshold
	total := threshold + UvarintSize(numSubs) + numSubs
	j = Pow2(numSubs)
	extras := make([]int, j)
	thresholds := make([]int, j)
	for i, _ = range thresholds {
		thresholds[i] = threshold
	}
	for _, sub := range subs {
		total += sub.Len()
		j >>= 1
		add := true
		extra := sub.Size() - sub.Len()
		for i, _ = range extras {
			if (i+1)%j == 0 {
				add = !add
			}
			if add && thresholds[i] > 0 {
				extras[i] += extra
				thresholds[i] -= sub.Weight()
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

func (f *fulfillmentThreshold) Len() int {
	p, _ := f.Condition().MarshalBinary()
	return len(p)
}

func (f *fulfillmentThreshold) Condition() *Condition {
	return NewCondition(THRESHOLD_ID, THRESHOLD_BITMASK, f.Hash(), f.Size(), f.Weight())
}

func (f *fulfillmentThreshold) String() string {
	p, _ := f.MarshalBinary()
	b64 := Base64UrlEncode(p)
	return Sprintf("cf:%x:%s", THRESHOLD_ID, b64)
}

func (f *fulfillmentThreshold) FromString(uri string) error {
	if !MatchString(FULFILLMENT_REGEX, uri) {
		return Error("Uri does not match fulfillment regex")
	}
	parts := Split(uri, ":")
	id, err := ParseUint16(parts[1])
	if err != nil {
		return err
	}
	if id != THRESHOLD_ID {
		return Errorf("Expected id=%d; got id=%d\n", THRESHOLD_ID, id)
	}
	bytes, err := Base64UrlDecode(parts[2])
	if err != nil {
		return err
	}
	return f.UnmarshalBinary(bytes)
}

func (f *fulfillmentThreshold) Validate(p []byte) bool { return true }
