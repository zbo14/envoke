package dsp

import (
	"bytes"
	"encoding/binary"
	"math"
	"math/cmplx"

	. "github.com/zbo14/envoke/common"
)

var (
	FAN_VALUE         = 15
	TIME_DELTA        = 200
	NEIGHBORHOOD_SIZE = 20
	OVERLAP_RATIO     = 0.5
	SAMPLING_RATE     = 44100
	WINDOW            = HammingWindow
	WINDOW_SIZE       = 4096
)

func CompareConstellations(c1, c2 [][]byte, eps float64) float64 {

	n, m := len(c1), len(c2)
	if n == 0 || m == 0 {
		return 0
	}
	go func() {
		lcs := 0
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if bytes.Equal(c1[i], c2[j]) {
					a, b := i, j
					prev, seq := b, 1
					for {
						if a++; a == n {
							break
						}
						for {
							if b++; b == m {
								b = prev
								break
							}
							if bytes.Equal(c1[a], c2[b]) {
								prev = b
								seq++
								break
							}
						}
					}
					if seq > lcs {
						lcs = seq
					}
				}
			}
		}
		ratio1, ratio2 := float64(lcs)/float64(n), float64(lcs)/float64(m)
		if math.Abs(ratio1-ratio2) < eps {
			ch <- math.Max(ratio1, ratio2)
		} else {
			ch <- 0
		}
	}()
	go func() {
		lcs := 0
		for i := 0; i < m; i++ {
			for j := 0; j < n; j++ {
				if bytes.Equal(c2[i], c1[j]) {
					a, b := i, j
					prev, seq := b, 1
					for {
						if a++; a == n {
							break
						}
						for {
							if b++; b == m {
								b = prev
								break
							}
							if bytes.Equal(c2[a], c1[b]) {
								prev = b
								seq++
								break
							}
						}
					}
					if seq > lcs {
						lcs = seq
					}
				}
			}
		}
		ratio1, ratio2 := float64(lcs)/float64(n), float64(lcs)/float64(m)
		if math.Abs(ratio1-ratio2) < eps {
			ch <- math.Max(ratio1, ratio2)
		} else {
			ch <- 0
		}
	}()
	return math.Max(<-ch, <-ch)
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

func CompareDistances(dists1, dists2 []float64, eps float64) float64 {
	ch := make(chan float64, 2)
	n, m := len(dists1), len(dists2)
	if n == 0 || m == 0 {
		return 0
	}
	go func() {
		lcs := 0
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if math.Abs(dists1[i]-dists2[j]) < eps {
					a, b := i, j
					prev, seq := b, 1
					for {
						if a++; a == n {
							break
						}
						for {
							if b++; b == m {
								b = prev
								break
							}
							if math.Abs(dists1[a]-dists2[b]) < eps {
								prev = b
								seq++
								break
							}
						}
					}
					if seq > lcs {
						lcs = seq
					}
				}
			}
		}
		ratio1, ratio2 := float64(lcs)/float64(n), float64(lcs)/float64(m)
		if math.Abs(ratio1-ratio2) < eps {
			ch <- math.Max(ratio1, ratio2)
		} else {
			ch <- 0
		}
	}()
	go func() {
		lcs := 0
		for i := 0; i < m; i++ {
			for j := 0; j < n; j++ {
				if math.Abs(dists2[i]-dists1[j]) < eps {
					a, b := i, j
					prev, seq := b, 1
					for {
						if a++; a == n {
							break
						}
						for {
							if b++; b == m {
								b = prev
								break
							}
							if math.Abs(dists2[a]-dists1[b]) < eps {
								prev = b
								seq++
								break
							}
						}
					}
					if seq > lcs {
						lcs = seq
					}
				}
			}
		}
		ratio1, ratio2 := float64(lcs)/float64(n), float64(lcs)/float64(m)
		if math.Abs(ratio1-ratio2) < eps {
			ch <- math.Max(ratio1, ratio2)
		} else {
			ch <- 0
		}
	}()
	return math.Max(<-ch, <-ch)
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

func DefaultFindPeaks(freqs []float64, sgram [][]float64) [][]float64 {
	return FindPeaks(freqs, NEIGHBORHOOD_SIZE, sgram)
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
			peaks[i] = append(peaks[i], x[j])
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
