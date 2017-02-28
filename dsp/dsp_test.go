package dsp

import (
	. "github.com/zbo14/envoke/common"
	"testing"
)

var (
	path1 = "/Users/zach/Desktop/music/hey_jude_1.mp3"
	path2 = "/Users/zach/Desktop/music/hey_jude_2.mp3"
)

func TestDsp(t *testing.T) {
	/*
		// FFT
		x := []float64{1.0, 2.0, 3.2, 4.5, 6.7, 7.01, 8.9, 10}
		z := FftReal(x)
		InvDecimationFreq(z)
		FftDecimationTime(a)
		t.Log(InvReal(z))

		// PEAKS
		freqs := []float64{100, 200, 300, 400, 500, 600, 700, 800}
		sgram := [][]float64{
			[]float64{1, 4, 2, 4, 10, 1, 2, 4},
			[]float64{1, 2, 3, 4, 5, 6, 7, 8},
			[]float64{3, 5, 6, 2, 3, 4, 3, 3},
			[]float64{3, 1, 6, 2, 2, 12, 3, 3},
		}
		t.Log(FindPeaks(freqs, 3, sgram))
	*/

	file := MustOpenFile(path1)
	x := MustReadTimeDomain(file)
	freqs, sgram, err := DefaultFftSpectrogram(x)
	if err != nil {
		t.Fatal(err)
	}
	peaks := DefaultFindPeaks(freqs, sgram)
	dists1 := DefaultFindDistances(peaks)

	file = MustOpenFile(path2)
	x = MustReadTimeDomain(file)
	freqs, sgram, err = DefaultFftSpectrogram(x)
	if err != nil {
		t.Fatal(err)
	}
	peaks = DefaultFindPeaks(freqs, sgram)
	dists2 := DefaultFindDistances(peaks)

	t.Log(DefaultCompareDistances(dists1, dists2))
}
