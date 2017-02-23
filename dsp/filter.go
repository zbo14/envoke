package dsp

import (
	"math"
)

func CustomFilter(m int, win func(int, int) float64, x []float64) []float64 {
	n := len(x)
	y := make([]float64, n)
	m2 := m >> 1
	for i := 0; i < n; i++ {
		y[i] = x[(i+m2)%n]
	}
	copy(x, y)
	for i := 0; i <= m; i++ {
		x[i] *= win(i, m)
	}
	return x[:m+1]
}

func BandPassFilter(cfl, cfh float64, m int, sinc func(float64, int) []float64) []float64 {
	hl := LowPassFilter(cfl, m, sinc)
	hh := HighPassFilter(cfh, m, sinc)
	h := make([]float64, m+1)
	m2 := m >> 1
	for i := range h {
		h[i] = -hl[i] - hh[i]
		if i == m2 {
			h[i]++
		}
	}
	return h
}

func LowPassFilter(cf float64, m int, sinc func(float64, int) []float64) []float64 {
	h := sinc(cf, m)
	sum := float64(0)
	for i := range h {
		sum += h[i]
	}
	for i := range h {
		h[i] /= sum
	}
	return h
}

func HighPassFilter(cf float64, m int, sinc func(float64, int) []float64) []float64 {
	h := sinc(cf, m)
	sum := float64(0)
	for i := range h {
		sum += h[i]
	}
	m2 := m >> 1
	for i := range h {
		h[i] /= -sum
		if i == m2 {
			h[i]++
		}
	}
	return h
}

func BlackmanSinc(cf float64, m int) []float64 {
	return SincKernel(cf, m, BlackmanWindow)
}

func HammingSinc(cf float64, m int) []float64 {
	return SincKernel(cf, m, HammingWindow)
}

func SincKernel(cf float64, m int, win func(int, int) float64) []float64 {
	if cf <= 0 || cf > 0.5 {
		panic("Cutoff frequency must be between 0 and 0.5")
	}
	if m&1 != 0 {
		panic("m must be an even number")
	}
	h := make([]float64, m+1)
	m2 := m >> 1
	for i := range h {
		if i == m2 {
			h[i] = 2 * math.Pi * cf
			continue
		}
		h[i] = math.Sin(2*math.Pi*cf*float64(i-m2)) * win(i, m) / float64(i-m2)
	}
	return h
}

func BlackmanWindow(i, m int) float64 {
	return 0.42 - 0.5*math.Cos(2*math.Pi*float64(i/m)) + 0.08*math.Cos(4*math.Pi*float64(i/m))
}

func HammingWindow(i, m int) float64 {
	return 0.54 - 0.46*math.Cos(2*math.Pi*float64(i/m))
}
