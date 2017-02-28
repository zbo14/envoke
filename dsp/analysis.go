package dsp

import (
	"bytes"
	"encoding/binary"
	"math"
	"math/cmplx"

	. "github.com/zbo14/envoke/common"
)

var (
	EPSILON           = 0.01
	FAN_VALUE         = 15
	NEIGHBORHOOD_SIZE = 20
	OVERLAP_RATIO     = 0.5
	SAMPLING_RATE     = 44100
	TIME_DELTA        = 200
	WINDOW            = HammingWindow
	WINDOW_SIZE       = 4096
)

func DefaultCompareDistances(dists1, dists2 []float64) float64 {
	return CompareDistances(dists1, dists2, EPSILON)
}

func CompareDistances(dists1, dists2 []float64, eps float64) (score float64) {
	n, m := len(dists1), len(dists2)
	if n == 0 || m == 0 {
		return 0
	}
	if n < m {
		score = float64(n)
	} else {
		score = float64(m)
	}
	lcs1 := LcsDistances(dists1, dists2, eps)
	lcs2 := LcsDistances(dists2, dists1, eps)
	if lcs1 > lcs2 {
		score = float64(lcs1) / score
	} else {
		score = float64(lcs2) / score
	}
	return
}

func LcsDistances(dists1, dists2 []float64, eps float64) int {
	n, m := len(dists1), len(dists2)
	if n == 0 || m == 0 {
		return 0
	}
	matrix := make([][]int, n+1)
	for i := range matrix {
		matrix[i] = make([]int, m+1)
	}
	for i, d1 := range dists1 {
		for j, d2 := range dists2 {
			if math.Abs(d1-d2)/math.Min(d1, d2) < eps {
				if matrix[i][j] > matrix[i+1][j] && matrix[i][j] > matrix[i][j+1] {
					matrix[i+1][j+1] = matrix[i][j] + 1
				} else if matrix[i+1][j] > matrix[i][j] && matrix[i+1][j] > matrix[i][j+1] {
					matrix[i+1][j+1] = matrix[i+1][j] + 1
				} else {
					matrix[i+1][j+1] = matrix[i][j+1] + 1
				}
			} else if matrix[i+1][j] > matrix[i][j+1] {
				matrix[i+1][j+1] = matrix[i+1][j]
			} else {
				matrix[i+1][j+1] = matrix[i][j+1]
			}
		}
	}
	return matrix[n][m]
}

func DefaultFindDistances(peaks [][]float64) []float64 {
	return FindDistances(FAN_VALUE, peaks, TIME_DELTA)
}

func FindDistances(fan int, peaks [][]float64, tdelta int) (dists []float64) {
	n := len(peaks)
	for i := range peaks {
		m := len(peaks[i])
	OUTER:
		for j := 0; j < m; j++ {
			rem := fan
			for k := j + 1; k < m; k++ {
				dists = append(dists, math.Abs(peaks[i][j]-peaks[i][k]))
				if rem--; rem == 0 {
					continue OUTER
				}
			}
			for a := i + 1; a < n && (tdelta <= 0 || a-i <= tdelta); a++ {
				for b := 0; b < len(peaks[a]); b++ {
					dists = append(dists, math.Abs(peaks[i][j]-peaks[a][b]))
					if rem--; rem == 0 {
						continue OUTER
					}
				}
			}
		}
	}
	return
}

func CompareConstellations(c1, c2 [][]byte) (score float64) {
	n, m := len(c1), len(c2)
	if n == 0 || m == 0 {
		return 0
	}
	if n < m {
		score = float64(n)
	} else {
		score = float64(m)
	}
	lcs1 := LcsConstellations(c1, c2)
	lcs2 := LcsConstellations(c2, c1)
	if lcs1 > lcs2 {
		score = float64(lcs1) / score
	} else {
		score = float64(lcs2) / score
	}
	return
}

func LcsConstellations(c1, c2 [][]byte) int {
	n, m := len(c1), len(c2)
	if n == 0 || m == 0 {
		return 0
	}
	matrix := make([][]int, n+1)
	for i := range matrix {
		matrix[i] = make([]int, m+1)
	}
	for i := range c1 {
		for j := range c2 {
			if bytes.Equal(c1[i], c2[j]) {
				if matrix[i][j] > matrix[i+1][j] && matrix[i][j] > matrix[i][j+1] {
					matrix[i+1][j+1] = matrix[i][j] + 1
				} else if matrix[i+1][j] > matrix[i][j] && matrix[i+1][j] > matrix[i][j+1] {
					matrix[i+1][j+1] = matrix[i+1][j] + 1
				} else {
					matrix[i+1][j+1] = matrix[i][j+1] + 1
				}
			} else if matrix[i+1][j] > matrix[i][j+1] {
				matrix[i+1][j+1] = matrix[i+1][j]
			} else {
				matrix[i+1][j+1] = matrix[i][j+1]
			}
		}
	}
	return matrix[n][m]
}

func DefaultConstellation(peaks [][]float64) [][]byte {
	return Constellation(FAN_VALUE, peaks, TIME_DELTA)
}

func Constellation(fan int, peaks [][]float64, tdelta int) [][]byte {
	var fprints [][]byte
	n := len(peaks)
	for i := range peaks {
		m := len(peaks[i])
	OUTER:
		for j := 0; j < m; j++ {
			rem := fan
			for k := j + 1; k < m; k++ {
				fprints = append(fprints, Fingerprint(peaks[i][j], peaks[i][k], 0))
				if rem--; rem == 0 {
					continue OUTER
				}
			}
			for a := i + 1; a < n && (tdelta <= 0 || a-i <= tdelta); a++ {
				for b := 0; b < len(peaks[a]); b++ {
					fprints = append(fprints, Fingerprint(peaks[i][j], peaks[a][b], a-i))
					if rem--; rem == 0 {
						continue OUTER
					}
				}
			}
		}
	}
	return fprints
}

func Fingerprint(peak1, peak2 float64, tdelta int) []byte {
	p := make([]byte, 16)
	binary.BigEndian.PutUint64(p[0:8], math.Float64bits(peak1))
	binary.BigEndian.PutUint64(p[8:], math.Float64bits(peak2))
	p = append(p, Uint64Bytes(tdelta)...)
	return Checksum256(p)
}

func DefaultFindPeaks(freqs []float64, sgram [][]float64) [][]float64 {
	return FindPeaks(freqs, NEIGHBORHOOD_SIZE, sgram)
}

func SequencePeaks(peaks [][]float64) (seq []float64) {
	for i := range peaks {
		seq = append(seq, peaks[i]...)
	}
	return
}

func FindPeaks(freqs []float64, nbr int, sgram [][]float64) [][]float64 {
	n := len(sgram)
	peaks := make([][]float64, n)
	for i, x := range sgram {
		m := len(x)
	OUTER:
		for j := range x {
			for k := 0; k <= nbr; k++ {
				if k > 0 {
					if j-k >= 0 {
						if x[j] <= x[j-k] {
							continue OUTER
						}
					}
					if j+k < m {
						if x[j] <= x[j+k] {
							continue OUTER
						}
					}
				}
				for b := 1; b <= nbr-k; b++ {
					if i-b >= 0 {
						if j-k >= 0 {
							if x[j] <= sgram[i-b][j-k] {
								continue OUTER
							}
						}
						if j+k < m {
							if x[j] <= sgram[i-b][j+k] {
								continue OUTER
							}
						}
					}
					if i+b < n {
						if j-k >= 0 {
							if x[j] <= sgram[i+b][j-k] {
								continue OUTER
							}
						}
						if j+k < m {
							if x[j] <= sgram[i+b][j+k] {
								continue OUTER
							}
						}
					}
				}
			}
			peaks[i] = append(peaks[i], freqs[j])
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
		for j := 0; j < l2p1; j++ {
			sgram[i][j] = math.Log10(real(z[j] * cmplx.Conj(z[j])))
		}
		sgram[i] = sgram[i][:l2p1]
	}
	freqs := make([]float64, l2p1)
	for i := range freqs {
		freqs[i] = float64(i) * float64(fs) / float64(l)
	}
	return freqs, sgram, nil
}
