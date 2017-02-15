package chroma

import (
	"bytes"
	"io"
	"os/exec"

	"github.com/go-fingerprint/fingerprint"
	. "github.com/zbo14/envoke/common"
	// acoustid "github.com/acoustid/go-acoustid/chromaprint"
	// "github.com/go-fingerprint/gochroma/chromaprint"
)

var THRESHOLD float64 = 0.95

func CompareFingerprints(fprint1, fprint2 string) (bool, error) {
	bytes1, err := Base64UrlDecode(fprint1)
	if err != nil {
		return false, err
	}
	raw1, err := Int32Slice(bytes1)
	if err != nil {
		return false, err
	}
	bytes2, err := Base64UrlDecode(fprint2)
	if err != nil {
		return false, err
	}
	raw2, err := Int32Slice(bytes2)
	if err != nil {
		return false, err
	}
	return CompareFingerprintsRaw(raw1, raw2)
}

func CompareFingerprintsRaw(raw1, raw2 []int32) (bool, error) {
	score, err := fingerprint.Compare(raw1, raw2)
	if err != nil {
		return false, err
	}
	Println(score)
	return score >= THRESHOLD, nil
}

func NewFingerprint(duration int, path string) (string, error) {
	raw, err := NewFingerprintRaw(duration, path)
	if err != nil {
		return "", err
	}
	p := Int32SliceBytes(raw)
	return Base64UrlEncode(p), nil
}

func NewFingerprintRaw(duration int, path string) (raw []int32, err error) {
	cmd := exec.Command("fpcalc", "-length", Itoa(duration), "-raw", path)
	stdout := new(bytes.Buffer)
	cmd.Stdout = stdout
	if err = cmd.Run(); err != nil {
		return nil, err
	}
	subs := SubmatchStr(`FINGERPRINT=(.*)`, string(stdout.Bytes()))
	strs := SplitStr(subs[1], ",")
	raw = make([]int32, len(strs))
	for i, s := range strs {
		raw[i], err = ParseInt32(s, 10)
	}
	return raw, nil
}

/*
const (
	CHANNELS    = 2
	MAX_SECONDS = 10
	SAMPLE_RATE = 44100
)

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

func hammingDist(raw1, raw2 []int32) int {
	dist := 0
	for i, x := range raw1 {
		y := x ^ raw2[i]
		for y > 0 {
			y >>= 1
			if y&1 == 1 {
				dist++
			}
		}
	}
	return dist
}

func SlicesRaw(raw []int32, size int) [][]int32 {
	slices := make([][]int32, len(raw)-size+1)
	for i := 0; i+size <= len(raw); i++ {
		slices[i] = raw[i : i+size]
	}
	return slices
}

len1 := len(raw1)
	len2 := len(raw2)
	if len1 > len2 {
		slices := SlicesRaw(raw1, len2)
		for _, slice := range slices {
			dist := hammingDist(slice, raw2)
			score := 1 - float64(dist)/float64(len2*32)
			Println(score)
			if score >= THRESHOLD {
				return true, nil
			}
		}
	} else {
		slices := SlicesRaw(raw2, len1)
		for _, slice := range slices {
			dist := hammingDist(raw1, slice)
			score := 1 - float64(dist)/float64(len1*32)
			Println(score)
			if score >= THRESHOLD {
				return true, nil
			}
		}
	}
	return false, nil
}
*/
