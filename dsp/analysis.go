package dsp

import (
	"math"
	"math/cmplx"
)

func FindPeaks(nbr int, sgram [][]float64) [][]bool {
	n := len(sgram)
	peaks := make([][]bool, n)
	for i, x := range s {
	OUTER:
		for j := range x {
			m := len(x)
			for k := 0; k <= nbr; k++ {
				if k != 0 {
					if j-k >= 0 {
						if x[j] < x[j-k] {
							continue OUTER
						}
					}
					if j+k < m {
						if x[j] < x[j+k] {
							continue OUTER
						}
					}
				}
				for b := 1; b <= nbr-k; b++ {
					if i-b >= 0 {
						if x[j] < s[i-b][k] {
							continue OUTER
						}
					}
					if i+b < n {
						if x[j] < s[i+b][k] {
							continue OUTER
						}
					}
				}
			}
			peaks[i][j] = true
		}
	}
	return peaks
}

func FftSpectrogram(fs float64, l int, olap float64, win func(int) []float64, x []float64) ([]float64, [][]float64, error) {
	l2p1 := l>>1 + 1
	sgram, err := LengthOverlap(l, olap, x)
	if err != nil {
		return nil, nil, err
	}
	for i := range sgram {
		ApplyWindow(win, s[i])
		z := FftReal(s[i])
		for j := 0; j < l2p1; j++ {
			sgram[i][j] = math.Log10(math.Sqrt(real(z[j] * cmplx.Conj(z[j]))))
		}
		sgram[i] = sgram[i][:l2p1]
	}
	freqs := make([]float64, l2p1)
	for i := range freqs {
		freqs[i] = float64(i) * fs / float64(l)
	}
	return freqs, sgram, nil
}
