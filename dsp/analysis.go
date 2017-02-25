package dsp

import (
	"encoding/binary"
	"math"
	"math/cmplx"

	. "github.com/zbo14/envoke/common"
)

var (
	FAN_VALUE            = 15
	NEIGHBORHOOD_SIZE    = 20
	OVERLAP_RATIO        = 0.5
	SAMPLING_RATE        = 44100
	SIMILARITY_THRESHOLD = 0.1
	WINDOW               = HammingWindow
	WINDOW_SIZE          = 4096
)

func DefaultCompareConstellations(c1, c2 map[string]struct{}) bool {
	return CompareConstellations(c1, c2, SIMILARITY_THRESHOLD)
}

func CompareConstellations(c1, c2 map[string]struct{}, threshold float64) bool {
	shared := 0
	for fprint := range c1 {
		if _, ok := c2[fprint]; ok {
			shared++
		}
	}
	ratio := 2 * float64(shared) / float64(len(c1)+len(c2))
	Println(ratio)
	return ratio >= threshold
}

func DefaultConstellation(peakFreqs [][]float64) map[string]struct{} {
	return Constellation(FAN_VALUE, peakFreqs)
}

func Constellation(fan int, peakFreqs [][]float64) map[string]struct{} {
	fprints := make(map[string]struct{})
	n := len(peakFreqs)
	for i, freqs := range peakFreqs {
		m := len(freqs)
	OUTER:
		for j := range freqs {
			rem := fan
			for k := j + 1; k < m; k++ {
				fprints[Fingerprint(freqs[j], freqs[k], 0)] = struct{}{}
				if rem--; rem == 0 {
					continue OUTER
				}
			}
			for a := i + 1; a < n; a++ {
				for b := 0; b < len(peakFreqs[a]); b++ {
					fprints[Fingerprint(freqs[j], peakFreqs[a][b], a-i)] = struct{}{}
					if rem--; rem == 0 {
						continue OUTER
					}
				}
			}
		}
	}
	return fprints
}

func Fingerprint(freq1, freq2 float64, tdelta int) string {
	p := make([]byte, 16)
	binary.BigEndian.PutUint64(p[0:8], math.Float64bits(freq1))
	binary.BigEndian.PutUint64(p[8:], math.Float64bits(freq2))
	p = append(p, Uint64Bytes(tdelta)...)
	return BytesToHex(Checksum256(p))
}

func DefaultPeakFrequencies(freqs []float64, sgram [][]float64) [][]float64 {
	return PeakFrequencies(freqs, NEIGHBORHOOD_SIZE, sgram)
}

func PeakFrequencies(freqs []float64, nbr int, sgram [][]float64) [][]float64 {
	n := len(sgram)
	peakFreqs := make([][]float64, n)
	for i, x := range sgram {
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
			peakFreqs[i] = append(peakFreqs[i], freqs[j])
		}
	}
	return peakFreqs
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
		for j := 0; j < l2p1; j++ {
			sgram[i][j] = math.Log10(math.Sqrt(real(z[j] * cmplx.Conj(z[j]))))
		}
		sgram[i] = sgram[i][:l2p1]
	}
	freqs := make([]float64, l2p1)
	for i := range freqs {
		freqs[i] = float64(i) * float64(fs) / float64(l)
	}
	return freqs, sgram, nil
}

/*
func SimilarityMeasure(dists1, dists2 []float64) (diff float64) {
	l1 := len(dists1)
	l2 := len(dists2)
	if l1 == l2 {
		for i, dist := range dists1 {
			diff += math.Pow(dist-dists2[i], 2)
		}
		return diff
	} else if l1 < l2 {
		diff = math.Inf(1)
		for i := 0; i+l1 <= l2; i++ {
			d := SimilarityMeasure(dists1, dists2[i:i+l1])
			if d < diff {
				diff = d
			}
		}
	} else {
		diff = math.Inf(1)
		for i := 0; i+l2 < l1; i++ {
			d := SimilarityMeasure(dists1[i:i+l2], dists2)
			if d < diff {
				diff = d
			}
		}
	}
	return diff
}
*/
