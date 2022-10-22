package utils

import (
	"math"

	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
)

func ComputePowers(x fr.Element, n int) []fr.Element {
	powers := make([]fr.Element, n)
	powers[0].SetOne()
	for i := 1; i < n; i++ {
		powers[i].Mul(&powers[i-1], &x)
	}
	return powers
}

func IsPowerOfTwo(value uint64) bool {
	return value > 0 && (value&(value-1) == 0)
}

// Raise an element to the power of two
// exp is of the form 2^y
// TODO: benchmark this versus using a .Exp which does not make the
// TODO assumption of power of 2. I expect this to be about 30% faster
// TODO: though since we are talking about nanoseconds, we can probably
// TODO use the non optimised version, since its only a small percentage
func Pow2(x fr.Element, exp uint64) *fr.Element {
	if !IsPowerOfTwo(exp) {
		// This can only happen if we stop using roots of unity whose domain is a power of 2
		panic("The domain should always be a power of two, for our use-case")
	}
	pow := int(math.Log2(float64(exp)))
	if pow == 0 {
		one := fr.One()
		return &one
	}

	var result = x
	for i := 0; i < pow; i++ {
		result.Square(&result)
	}
	return &result
}