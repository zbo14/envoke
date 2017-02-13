package chroma

import (
	"github.com/go-fingerprint/fingerprint"
	"github.com/go-fingerprint/gochroma"
	. "github.com/zbo14/envoke/common"
	"io"
)

const (
	CHANNELS    = 2
	MAX_SECONDS = 120
	SAMPLE_RATE = 44100
)

var THRESHOLD float64 = 0.95

// base64 url-safe
var Encode = Base64UrlEncode
var Decode = Base64UrlDecode

func NewFingerprint(r io.Reader) (string, error) {
	raw, err := NewFingerprintRaw(r)
	if err != nil {
		return "", err
	}
	p := Int32SliceBytes(raw)
	return Encode(p), nil
}

func NewFingerprintRaw(r io.Reader) ([]int32, error) {
	fpcalc := gochroma.New(gochroma.AlgorithmDefault)
	defer fpcalc.Close()
	raw, err := fpcalc.RawFingerprint(
		fingerprint.RawInfo{
			Channels:   CHANNELS,
			MaxSeconds: MAX_SECONDS,
			Rate:       SAMPLE_RATE,
			Src:        r,
		})
	if err != nil {
		return nil, err
	}
	return raw, nil
}

func NewFingerprintCompressed(r io.Reader) (string, error) {
	fpcalc := gochroma.New(gochroma.AlgorithmDefault)
	defer fpcalc.Close()
	compressed, err := fpcalc.Fingerprint(
		fingerprint.RawInfo{
			Channels:   CHANNELS,
			MaxSeconds: MAX_SECONDS,
			Rate:       SAMPLE_RATE,
			Src:        r,
		})
	if err != nil {
		return "", err
	}
	return compressed, nil
}

func FingerprintToRaw(fprint string) ([]int32, error) {
	p, err := Decode(fprint)
	if err != nil {
		return nil, err
	}
	return Int32Slice(p)
}

func FingerprintFromRaw(raw []int32) string {
	p := Int32SliceBytes(raw)
	return Encode(p)
}

func CompareFingerprints(fprint1, fprint2 string) (bool, error) {
	p1, err := Decode(fprint1)
	if err != nil {
		return false, err
	}
	raw1, err := Int32Slice(p1)
	if err != nil {
		return false, err
	}
	p2, err := Decode(fprint2)
	if err != nil {
		return false, err
	}
	raw2, err := Int32Slice(p2)
	if err != nil {
		return false, err
	}
	return CompareFingerprintsRaw(raw1, raw2)
}

func CompareFingerprintsRaw(raw1, raw2 []int32) (bool, error) {
	s, err := fingerprint.Compare(raw1, raw2)
	if err != nil {
		return false, err
	}
	return s >= THRESHOLD, nil
}
