package common

import (
	"math"
	"math/big"
)

func Pow2(x int) bool {
	return x != 0 && (x&(x-1)) == 0
}

func Exp2(x int) int {
	return 1 << uint(x)
}

func Pow2Ceil(x int) int {
	return Exp2(Log2Ceil(x))
}

func Pow2Floor(x int) int {
	return Exp2(Log2Floor(x))
}

func Log2Floor(x int) int {
	i, log := x, 0
	for {
		if i >>= 1; i == 0 {
			return log
		}
		log++
	}
}

func Log2Ceil(x int) int {
	log := Log2Floor(x)
	if Pow2(x) {
		return log
	}
	return log + 1
}

func EvenSquare(n int) bool {
	sqrt := math.Sqrt(float64(n))
	return float64(int(sqrt)) != sqrt
}

func BigIntFromBytes(p []byte) *big.Int {
	x := new(big.Int)
	return x.SetBytes(p)
}
