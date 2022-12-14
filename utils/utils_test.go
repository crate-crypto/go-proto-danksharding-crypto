package utils

import (
	"bytes"
	"math"
	"math/big"
	"testing"

	curve "github.com/consensys/gnark-crypto/ecc/bls12-381"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
)

func TestReverseSlice(t *testing.T) {
	type TestCase struct {
		slice, reversedSlice []byte
	}

	var testCases = []TestCase{
		TestCase{[]byte{1, 2, 3, 4}, []byte{4, 3, 2, 1}},
		TestCase{[]byte{1, 2, 3, 4, 5}, []byte{5, 4, 3, 2, 1}},
		TestCase{[]byte{1}, []byte{1}},
		TestCase{[]byte{}, []byte{}},
	}

	for _, test := range testCases {
		got := test.slice
		expected := test.reversedSlice
		ReverseSlice(got)

		if !bytes.Equal(got, expected) {
			t.Error("expected reversed slice does not match the computed reversed slice")
		}
	}

}
func TestIsPow2(t *testing.T) {
	powInt := func(x, y uint64) uint64 {
		return uint64(math.Pow(float64(x), float64(y)))
	}

	// 0 is now a power of two
	ok := IsPowerOfTwo(0)
	if ok {
		t.Error("zero is not a power of two")
	}

	// Numbers of the form 2^x are all powers of two
	// Do this up to x=63, since we are using u64
	for i := 0; i < 63; i++ {
		pow2 := powInt(2, uint64(i))
		ok := IsPowerOfTwo(pow2)
		if !ok {
			t.Error("numbers of the form 2^x are powers of two")
		}
	}
	// Numbers of the form 2^x -1 are not powers of two
	// from x=2 until x=63
	for i := 2; i < 63; i++ {
		pow2Minus1 := powInt(2, uint64(i)) - 1
		ok := IsPowerOfTwo(pow2Minus1)
		if ok {
			t.Error("numbers of the form 2^x -1 are not powers of two from x=2")
		}
	}
}

func TestComputePowersBaseOne(t *testing.T) {
	one := fr.One()

	powers := ComputePowers(one, 10)
	for _, pow := range powers {
		if !pow.Equal(&one) {
			t.Error("powers should all be 1")
		}
	}
}

func TestComputePowersZero(t *testing.T) {
	x := fr.NewElement(1234)

	powers := ComputePowers(x, 0)
	// When given a number of 0
	// this will return an empty slice
	if len(powers) != 0 {
		t.Error("number of powers to compute was zero, but got more than 0 powers computed")
	}
}

func TestComputePowersSmoke(t *testing.T) {
	var base fr.Element
	base.SetInt64(123)

	powers := ComputePowers(base, 16)

	for index, pow := range powers {
		var expected fr.Element
		expected.Exp(base, big.NewInt(int64(index)))

		if !expected.Equal(&pow) {
			t.Error("incorrect exponentiation result")
		}
	}
}

func TestReversal(t *testing.T) {
	powInt := func(x, y int) int {
		return int(math.Pow(float64(x), float64(y)))
	}

	// We only go up to 20 because we don't want a long running test
	for i := 0; i < 20; i++ {
		size := powInt(2, i)

		scalars := randomScalars(size)
		reversed := bitReversalPermutation(scalars)

		BitReverseRoots(scalars)

		for i := 0; i < size; i++ {
			if !reversed[i].Equal(&scalars[i]) {
				t.Error("bit reversal methods are not consistent")
			}
		}

	}

}

func TestExponentiate(t *testing.T) {
	var base fr.Element
	base.SetInt64(123)
	var result fr.Element

	result.Exp(base, big.NewInt(16))
	res2 := Pow2(base, 16)

	if !res2.Equal(&result) {
		t.Fail()
	}
}

func TestBatchNormalisation(t *testing.T) {
	numPoints := 100
	g1JacGen, _, _, _ := curve.Generators()
	points := make([]curve.G1Jac, numPoints)
	points[0] = g1JacGen

	for i := 1; i < numPoints; i++ {
		points[i-1].Double(&points[i])
	}

	// Set one of the points to the point at infinity
	points[numPoints/2] = curve.G1Jac{}

	expected := make([]curve.G1Affine, numPoints)
	for i := 0; i < numPoints; i++ {
		expected[i].FromJacobian(&points[i])
	}
	got := BatchFromJacobian(points)

	for i := 0; i < numPoints; i++ {
		if !got[i].Equal(&expected[i]) {
			t.Errorf("batch normalisation produced the wrong output. Check index %d", i)
		}
	}

}

func TestArrReverse(t *testing.T) {
	arr := [32]uint8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
		21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
		31, 32,
	}
	ReverseArray(&arr)
	expected := [32]uint8{32, 31, 30, 29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}

	if !bytes.Equal(expected[:], arr[:]) {
		t.Error("bytes are not equal")
	}
}

func TestCanonicalEncoding(t *testing.T) {

	x := randReducedBigInt()
	var xPlusModulus = addModP(x)

	unreducedBytes := xPlusModulus.Bytes()

	// `SetBytes` will read the unreduced bytes and
	// return a field element. Does not matter if its canonical
	var reduced fr.Element
	reduced.SetBytes(unreducedBytes)

	// `Bytes` will return a canonical representation of the
	// field element, ie a reduced version
	reducedBytes := reduced.Bytes()

	// First we should check that the reduced version
	// is different to the unreduced version, incase one changes the
	// implementation in the future
	if bytes.Equal(unreducedBytes, reducedBytes[:]) {
		t.Error("unreduced representation of field element, is the same as the reduced representation")
	}

	// Reduce canonical should produce the same result
	scalar, isReduced := ReduceCanonical(unreducedBytes)
	if isReduced {
		t.Error("input to ReduceCanonical was unreduced bytes")
	}
	if !scalar.Equal(&reduced) {
		t.Error("incorrect field element interpretation from unreduced byte representation")
	}
}

func addModP(x big.Int) big.Int {
	modulus := fr.Modulus()

	var x_plus_modulus big.Int
	x_plus_modulus.Add(&x, modulus)

	return x_plus_modulus
}

func randReducedBigInt() big.Int {
	var randFr fr.Element
	_, _ = randFr.SetRandom()

	var randBigInt big.Int
	randFr.ToBigIntRegular(&randBigInt)

	if randBigInt.Cmp(fr.Modulus()) != -1 {
		panic("big integer is not reduced")
	}

	return randBigInt
}

func randomScalars(size int) []fr.Element {
	res := make([]fr.Element, size)
	for i := 0; i < size; i++ {
		res[i] = fr.NewElement(uint64(i))
	}
	return res
}
