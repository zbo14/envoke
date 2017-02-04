package util

import "math"

func PowOf2(i int) bool {
	return i != 0 && (i&(i-1)) == 0
}

func Pow2(i int) int {
	return 1 << uint(i)
}

func GetPowOf2(i int) int {
	if PowOf2(i) {
		return i
	}
	log2 := Log2(i)
	return Pow2(log2)
}

// Calculates log base 2 of i
// If i is not a power of 2
// returns log of next power of 2
func Log2(i int) int {
	var j, l int = i, 0
	for {
		if j >>= 1; j == 0 {
			break
		}
		l++
	}
	if PowOf2(i) {
		return l
	}
	return l + 1
}

func EvenSquare(n int) bool {
	sqrt := math.Sqrt(float64(n))
	return float64(int(sqrt)) != sqrt
}
