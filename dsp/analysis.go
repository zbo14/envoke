package dsp

import (
	"encoding/binary"
	"math"
	"math/cmplx"

	. "github.com/zbo14/envoke/common"
)

var (
	NEIGHBORHOOD_SIZE    = 20
	OVERLAP_RATIO        = 0.5
	SAMPLING_RATE        = 44100
	SIMILARITY_THRESHOLD = 0.3
	WINDOW               = HammingWindow
	WINDOW_SIZE          = 4096
)

func Fingerprint(dists []float64) []byte {
	p := make([]byte, 8*len(dists))
	for i, dist := range dists {
		binary.BigEndian.PutUint64(p[i*8:(i+1)*8], math.Float64bits(dist))
	}
	return Checksum256(p)
}

func SimilarityMeasure(dists1, dists2 []float64) (diff float64) {
	var i int
	l1 := len(dists1)
	l2 := len(dists2)
	if l1 == l2 {
		for i = range dists1 {
			diff += math.Pow(dists1[i]-dists2[i], 2)
		}
		return diff / float64(l1)
	} else if l1 < l2 {
		for i = 0; i+l1 < l2; i++ {
			diff += SimilarityMeasure(dists1, dists2[i:i+l1])
		}
	} else {
		for i = 0; i+l2 < l1; i++ {
			diff += SimilarityMeasure(dists1[i:i+l2], dists2)
		}
	}
	return diff
}

func PeakDistances(freqs []float64, peaks [][]bool) (dists []float64) {
	freqMin, freqMax := freqs[0], freqs[len(freqs)-1]
	freqRange := freqMax - freqMin
	prevFreq := freqMin
	prevTime := 0
	n := len(peaks)
	for i := range peaks {
		for j := range peaks[i] {
			if peaks[i][j] {
				dists = append(dists, math.Pow((freqs[j]-prevFreq)/freqRange, 2)+math.Pow(float64(i-prevTime)/float64(n-1), 2))
				prevFreq, prevTime = freqs[j], i
			}
		}
	}
	return
}

func DefaultFindPeaks(sgram [][]float64) [][]bool {
	return FindPeaks(NEIGHBORHOOD_SIZE, sgram)
}

func FindPeaks(nbr int, sgram [][]float64) [][]bool {
	n := len(sgram)
	peaks := make([][]bool, n)
	for i, x := range sgram {
	OUTER:
		for j := range x {
			m := len(x)
			peaks[i] = make([]bool, m)
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
						if x[j] < sgram[i-b][k] {
							continue OUTER
						}
					}
					if i+b < n {
						if x[j] < sgram[i+b][k] {
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

func DefaultFftSpectrogram(x []float64) ([]float64, [][]float64, error) {
	return FftSpectrogram(SAMPLING_RATE, WINDOW_SIZE, OVERLAP_RATIO, WINDOW, x)
}

func FftSpectrogram(fs, l int, olap float64, win func(int) []float64, x []float64) ([]float64, [][]float64, error) {
	l2p1 := l>>1 + 1
	sgram, err := LengthOverlap(l, olap, x)
	if err != nil {
		return nil, nil, err
	}
	for i := range sgram {
		ApplyWindow(win, sgram[i])
		z := FftReal(sgram[i])
		min, max := math.Inf(1), math.Inf(-1)
		for j := 0; j < l2p1; j++ {
			sgram[i][j] = math.Sqrt(real(z[j] * cmplx.Conj(z[j])))
			if sgram[i][j] < min {
				min = sgram[i][j]
			}
			if sgram[i][j] > max {
				max = sgram[i][j]
			}
		}
		for j := 0; j < l2p1; j++ {
			sgram[i][j] -= min
			sgram[i][j] /= (max - min)
		}
		sgram[i] = sgram[i][:l2p1]
	}
	freqs := make([]float64, l2p1)
	for i := range freqs {
		freqs[i] = float64(i) * float64(fs) / float64(l)
	}
	return freqs, sgram, nil
}
