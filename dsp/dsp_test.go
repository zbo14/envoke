package dsp

import (
	. "github.com/zbo14/envoke/common"
	"testing"
)

var path = "/Users/zach/Desktop/music/Allegro from Duet in C Major.mp3"

func TestDsp(t *testing.T) {
	// FFT
	// x := []float64{1.0, 2.0, 3.2, 4.5, 6.7, 7.01, 8.9, 10}
	// z := FftReal(x)
	// InvDecimationFreq(z)
	// FftDecimationTime(a)
	// t.Log(InvReal(z))
	file := MustOpenFile(path)
	x := MustReadTimeDomain(file)
	s, freqs, err := FftSpectrogram(44100, 1024, 0.5, BlackmanWindow, x)
	if err != nil {
		t.Fatal(err)
	}
	file = MustOpenWriteFile("test.json")
	t.Log(s)
	t.Log(freqs)
	// Println(out)
	// Convolution
	// out := FftConvolution(0.21, 398, 625, x)
	// t.Log(out)
}
