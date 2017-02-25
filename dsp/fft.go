package dsp

import (
	"math"

	. "github.com/zbo14/envoke/common"
)

func FftReal(x []float64) []complex128 {
	n := len(x)
	if !Pow2(n) {
		panic("Size must be power of 2")
	}
	n2, n4 := n>>1, n>>2
	z := make([]complex128, n)
	for i := 0; i < n2; i++ {
		z[i] = complex(x[i<<1], x[i<<1+1])
	}
	FftDecimationTime(z[:n2])
	for i := 1; i < n4; i++ {
		rl1, im1 := real(z[i]), imag(z[i])
		rl2, im2 := real(z[n2-i]), imag(z[n2-i])
		z[i+n2] = complex((im1+im2)/2, -(rl1-rl2)/2)
		z[n-i] = complex(real(z[i+n2]), -imag(z[i+n2]))
		z[i] = complex((rl1+rl2)/2, (im1-im2)/2)
		z[n2-i] = complex(real(z[i]), -imag(z[i]))
	}
	z[n*3/4] = complex(imag(z[n4]), 0)
	z[n2] = complex(imag(z[0]), 0)
	z[n4] = complex(real(z[n4]), 0)
	z[0] = complex(real(z[0]), 0)
	sin, cos := math.Sincos(math.Pi / float64(n2))
	s := complex(cos, -sin)
	u := complex(1, 0)
	for i := 0; i < n2; i++ {
		c := z[i+n2] * u
		z[i], z[i+n2] = z[i]+c, z[i]-c
		u *= s
	}
	return z
}

func InvReal(z []complex128) []float64 {
	n := len(z)
	if !Pow2(n) {
		panic("Size must be power of 2")
	}
	n2 := n >> 1
	for i := 1; i < n2; i++ {
		z[i+n2] = complex(real(z[n2-i]), -imag(z[n2-i]))
	}
	x := make([]float64, n)
	for i := range x {
		x[i] = real(z[i]) + imag(z[i])
	}
	z = FftReal(x)
	for i := range x {
		x[i] = (real(z[i]) + imag(z[i])) / float64(n)
	}
	return x
}

func FftDecimationTime(z []complex128) {
	n := len(z)
	if !Pow2(n) {
		panic("Size must be power of 2")
	}
	log2 := Log2Floor(n)
	zz := make([]complex128, n)
	copy(zz, z)
	for i := 0; i < n; i++ {
		a := 0
		for j := 0; j < log2; j++ {
			a |= ((i >> uint(log2-j-1)) << uint(j)) & (1 << uint(j))
		}
		z[i] = zz[a%n]
	}
	for a := 2; a <= n; a <<= 1 {
		b := a >> 1
		sin, cos := math.Sincos(math.Pi / float64(b))
		s := complex(cos, -sin)
		u := complex(1, 0)
		for i := 0; i < b; i++ {
			for j := i; j < n; j += a {
				c := z[j+b] * u
				z[j], z[j+b] = z[j]+c, z[j]-c
			}
			u *= s
		}
	}
}

func InvDecimationTime(z []complex128) {
	n := len(z)
	if !Pow2(n) {
		panic("Size must be power of 2")
	}
	for i := range z {
		z[i] = complex(real(z[i]), -imag(z[i]))
	}
	FftDecimationTime(z)
	for i := range z {
		z[i] = complex(real(z[i])/float64(n), -imag(z[i])/float64(n))
	}
}

func FftDecimationFreq(z []complex128) {
	n := len(z)
	if !Pow2(n) {
		panic("Size must be power of 2")
	}
	for a := n; a >= 2; a >>= 1 {
		b := a >> 1
		sin, cos := math.Sincos(math.Pi / float64(b))
		s := complex(cos, -sin)
		u := complex(1, 0)
		for i := 0; i < b; i++ {
			for j := i; j < n; j += a {
				z[j], z[j+b] = z[j]+z[j+b], u*(z[j]-z[j+b])
			}
			u *= s
		}
	}
	for i, j := 1, 1; i <= n; i++ {
		if i < j {
			z[i-1], z[j-1] = z[j-1], z[i-1]
		}
		k := n >> 1
		for k < j && k > 0 {
			j -= k
			k >>= 1
		}
		j += k
	}
}

func InvDecimationFreq(z []complex128) {
	n := len(z)
	if !Pow2(n) {
		panic("Size must be power of 2")
	}
	for i := range z {
		z[i] = complex(real(z[i]), -imag(z[i]))
	}
	FftDecimationFreq(z)
	for i := range z {
		z[i] = complex(real(z[i])/float64(n), -imag(z[i])/float64(n))
	}
}

/*
func FftConvolution(cf float64, m, seg int, x []float64) []float64 {
	h := HighPassFilter(cf, m, BlackmanSinc)
	h = Pad(seg, h)
	hz := FftReal(h)
	h2 := len(hz) >> 1
	s, pad := PadSegment(seg, x)
	olap := make([]float64, m+1)
	out := make([]float64, len(x)+pad)
	for i, y := range s {
		y = Pad(m+1, y)
		z := FftReal(y)
		for j := 0; j <= h2; j++ {
			z[j] *= hz[j]
		}
		y = InvReal(z)
		for j := 0; j < seg; j++ {
			if math.IsNaN(y[j]) {
				y[j] = 0
			}
			if j <= m {
				y[j] += olap[j]
				if math.IsNaN(y[j+seg]) {
					y[j+seg] = 0
				}
				olap[j] = y[j+seg]
			}
		}
		copy(out[i*seg:(i+1)*seg], y[:seg])
	}
	return append(out, olap...)
}
*/
