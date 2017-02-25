package dsp

import (
	"io"
	"math"

	. "github.com/zbo14/envoke/common"
)

// Window

func ApplyWindow(win func(int) []float64, x []float64) {
	h := win(len(x))
	for i := range x {
		x[i] *= h[i]
	}
}

func BlackmanWindow(l int) []float64 {
	return Window(Blackman, l)
}

func HammingWindow(l int) []float64 {
	return Window(Hamming, l)
}

func Window(fn func(int, int) float64, l int) []float64 {
	h := make([]float64, l)
	if l == 1 {
		h[0] = 1
		return h
	}
	for i := 0; i < l; i++ {
		h[i] = fn(i, l-1)
	}
	return h
}

func Blackman(i, l int) float64 {
	return 0.42 - 0.5*math.Cos(2*float64(i)*math.Pi/float64(l)) + 0.08*math.Cos(4*float64(i)*math.Pi/float64(l))
}

func Hamming(i, l int) float64 {
	return 0.54 - 0.46*math.Cos(2*float64(i)*math.Pi/float64(l))
}

// Signal

func MustReadTimeDomain(r io.Reader) []float64 {
	x, err := ReadTimeDomain(r, 1024)
	Check(err)
	return x
}

func ReadTimeDomain(r io.Reader, sz int) ([]float64, error) {
	x := make([]float64, sz)
	for i := 0; ; i++ {
		if i == len(x) {
			x = append(x, make([]float64, sz)...)
		}
		i16, err := ReadInt16(r)
		if err != nil {
			// TODO: get rid of ErrUnexpectedEOF
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				return x[:i], nil
			}
			return nil, err
		}
		x[i] = float64(i16)
	}
}

// Segment

func PadSegment(seg int, x []float64) ([][]float64, int) {
	if mod := len(x) % seg; mod != 0 {
		pad := seg - mod
		return Segment(seg, Pad(pad, x)), pad
	}
	return Segment(seg, x), 0
}

func PadSegmentOverlap(olap float64, seg int, x []float64) (s [][]float64, pad int) {
	// TODO: optimize
	var err error
	for {
		s, err = SegmentOverlap(olap, seg, x)
		if err == nil {
			return s, pad
		}
		x = Pad(1, x)
		pad++
	}
}

func Pad(n int, x []float64) []float64 {
	return append(x, make([]float64, n)...)
}

func PadTo(n int, x []float64) []float64 {
	if n <= len(x) {
		return x
	}
	y := make([]float64, n)
	copy(y, x)
	return y
}

func PadToPow2(x []float64) []float64 {
	return PadTo(Pow2Ceil(len(x)), x)
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

func SegmentOverlap(olap float64, seg int, x []float64) ([][]float64, error) {
	n := len(x)
	l, o := n/seg, 0
	for ; l < n; l++ {
		o = int(olap * float64(l))
		if l*seg-o*(seg-1) == n {
			break
		}
	}
	if l == n {
		return nil, ErrorAppend(ErrInvalidSize, "cannot divide signal into segments")
	}
	s := make([][]float64, seg)
	for i, j := 0, 0; i < seg; i, j = i+1, j+l-o {
		s[i] = x[j : j+l]
	}
	return s, nil
}

func LengthOverlap(l int, olap float64, x []float64) ([][]float64, error) {
	n := len(x)
	o := int(olap * float64(l))
	var seg int
	if n == l {
		seg = 1
	} else if n > l {
		seg = (n + o) / (l - o)
	} else {
		return nil, ErrorAppend(ErrInvalidSize, "cannot divide signal into segments")
	}
	s := make([][]float64, seg)
	for i, j := 0, 0; i < seg; i, j = i+1, j+l-o {
		if i == seg-1 {
			s[i] = PadTo(l, x[j:])
		} else {
			s[i] = x[j : j+l]
		}
	}
	return s, nil
}

/*
// Filter

func CustomFilter(l int, win func(int, int) float64, x []float64) []float64 {
	n := len(x)
	y := make([]float64, n)
	l2 := l >> 1
	for i := 0; i < n; i++ {
		y[i] = x[(i+l2)%n]
	}
	copy(x, y)
	for i := 0; i <= l; i++ {
		x[i] *= win(i, l)
	}
	return x[:l+1]
}

func BandPassFilter(cfl, cfh float64, l int, sinc func(float64, int) []float64) []float64 {
	hl := LowPassFilter(cfl, l, sinc)
	hh := HighPassFilter(cfh, l, sinc)
	h := make([]float64, l+1)
	l2 := l >> 1
	for i := range h {
		h[i] = -hl[i] - hh[i]
		if i == l2 {
			h[i]++
		}
	}
	return h
}

func LowPassFilter(cf float64, l int, sinc func(float64, int) []float64) []float64 {
	h := sinc(cf, l)
	sum := float64(0)
	for i := range h {
		sum += h[i]
	}
	for i := range h {
		h[i] /= sum
	}
	return h
}

func HighPassFilter(cf float64, l int, sinc func(float64, int) []float64) []float64 {
	h := sinc(cf, l)
	sum := float64(0)
	for i := range h {
		sum += h[i]
	}
	l2 := l >> 1
	for i := range h {
		h[i] /= -sum
		if i == l2 {
			h[i]++
		}
	}
	return h
}

// Sinc

func BlackmanSinc(cf float64, l int) []float64 {
	return SincKernel(cf, Blackman, l)
}

func HammingSinc(cf float64, l int) []float64 {
	return SincKernel(cf, Hamming, l)
}

func SincKernel(cf float64, fn func(int, int) float64, l int) []float64 {
	if cf <= 0 || cf > 0.5 {
		panic("Cutoff frequency must be between 0 and 0.5")
	}
	if l&1 != 0 {
		panic("m must be an even number")
	}
	h := make([]float64, l+1)
	l2 := l >> 1
	for i := range h {
		if i == l2 {
			h[i] = 2 * math.Pi * cf
			continue
		}
		h[i] = math.Sin(2*math.Pi*cf*float64(i-l2)) * fn(i, l) / float64(i-l2)
	}
	return h
}
*/
