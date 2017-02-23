package dsp

import (
	. "github.com/zbo14/envoke/common"
	"testing"
)

var path = "/Users/zach/Desktop/music/Allegro from Duet in C Major.mp3"

func TestDsp(t *testing.T) {
	// FFT
	x := []float64{1.0, 2.0, 3.2, 4.5, 6.7, 7.01, 8.9, 10}
	z := FftReal(x)
	InvDecimationFreq(z)
	FftDecimationTime(z)
	t.Log(InvReal(z))
	// Convolution
	file, err := OpenFile(path)
	if err != nil {
		t.Fatal(err)
	}
	x, err = ReadTimeDomain(file)
	if err != nil {
		t.Fatal(err)
	}
	out := FftConvolution(0.21, 398, 625, x)
	file = MustOpenWriteFile("test.json")
	if err = WriteJSON(file, out); err != nil {
		t.Fatal(err)
	}
	// t.Log(out)
}
