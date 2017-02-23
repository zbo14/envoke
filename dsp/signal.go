package dsp

import (
	. "github.com/zbo14/envoke/common"
	"io"
)

func MustReadTimeDomain(r io.Reader) []float64 {
	x, err := ReadTimeDomain(r)
	Check(err)
	return x
}

func ReadTimeDomain(r io.Reader) ([]float64, error) {
	return ReadFloat64s(r, 1024)
}

func PadAndSegment(seg int, x []float64) ([][]float64, int) {
	if mod := len(x) % seg; mod != 0 {
		pad := seg - mod
		return Segment(seg, Pad(pad, x)), pad
	}
	return Segment(seg, x), 0
}

func Pad(n int, x []float64) []float64 {
	return append(x, make([]float64, n)...)
}

func Segment(seg int, x []float64) [][]float64 {
	n := len(x)
	if n%seg != 0 {
		panic("Segment length must evenly divide signal length")
	}
	s := make([][]float64, n/seg)
	for i := range s {
		s[i] = x[i*seg : (i+1)*seg]
	}
	return s
}
