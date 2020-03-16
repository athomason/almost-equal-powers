package t

import (
	"math/big"
	"math/bits"
	"reflect"
	"testing"
)

func TestBigIntPowerOfTwoBits(t *testing.T) {
	var (
		two  = big.NewInt(2)
		pTwo = new(big.Int)
	)
	for twosExp := 1; twosExp < 100; twosExp++ {
		pTwo.Exp(two, big.NewInt(int64(twosExp)), nil)
		want := pTwo.Bits()

		// big.Int has a sign bit and a little endian slice of words: see
		// math/big/nat.go. we can build a power of two by setting the correct
		// bit in a word after the right number of zeros.
		zeros, mod := twosExp/bits.UintSize, twosExp%bits.UintSize
		got := append(make([]big.Word, zeros), 1<<mod)

		if !reflect.DeepEqual(want, got) {
			t.Errorf("%v != %v", want, got)
		}

		// round trip
		rev := new(big.Int)
		rev.SetBits(got)
		if rev.Cmp(pTwo) != 0 {
			t.Errorf("%v != %v", rev, pTwo)
		}
	}
}
