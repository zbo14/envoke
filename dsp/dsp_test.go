package dsp

import (
	. "github.com/zbo14/envoke/common"
	"testing"
)

var (
	path1 = "/Users/zach/Desktop/music/rhapsody_1.mp3"
	path2 = "/Users/zach/Desktop/music/rhapsody_2.mp3"
)

func TestDsp(t *testing.T) {
	// FFT
	// x := []float64{1.0, 2.0, 3.2, 4.5, 6.7, 7.01, 8.9, 10}
	// z := FftReal(x)
	// InvDecimationFreq(z)
	// FftDecimationTime(a)
	// t.Log(InvReal(z))
	file := MustOpenFile(path1)
	x := MustReadTimeDomain(file)
	freqs, sgram, err := DefaultFftSpectrogram(x)
	if err != nil {
		t.Fatal(err)
	}
	peakFreqs := DefaultPeakFrequencies(freqs, sgram)
	constellation1 := DefaultConstellation(peakFreqs)
	file = MustOpenFile(path2)
	x = MustReadTimeDomain(file)
	freqs, sgram, err = DefaultFftSpectrogram(x)
	if err != nil {
		t.Fatal(err)
	}
	peakFreqs = DefaultPeakFrequencies(freqs, sgram)
	constellation2 := DefaultConstellation(peakFreqs)
	t.Log(DefaultCompareConstellations(constellation1, constellation2))
}
