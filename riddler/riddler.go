/*

This program is an over-optimized solution to a generalization of the problem
posed by this blog post:

https://fivethirtyeight.com/features/can-you-get-another-haircut-already/

	a kilobyte is in fact 1,024... thanks to the happy coincidence that there’s
	a power of two that’s very close to 1,000. Working through the numbers, 210
	is a mere 2.4 percent more than 103...

	After 2^10, what’s the next (whole number) power of 2 that comes closer to
	a power of 10? (To be clear, “closer” doesn’t refer to the absolute
	difference — it means your power of 2 should differ from a power of 10 by
	less than 2.4 percent.)

Figuring out how close the nearest power of two is to a given power of ten
takes some algebra:

	10^tensExponent = 2^twosExponent
	log2(10^tensExponent) = log2(2^twosExponent)
	tensExponent*log2(10) = twosExponent

But twosExponent is irrational; we want the nearest integer power of two, so
let's round:

	twosExponent = round(tensExponent*log2(10))

We'll define a measure of "closeness", matching the problem description as a
percentage relative to the power of ten:

	error = abs(twosPower-tensPower)/tensPower

The partial answer computed by this program is

1, 3, 28, 59, 146, 643, 4004, 8651, 12655, 21306, 76573, 97879, ...

which curiously appears as OEIS A046104 and A116984, both of which concern
continuing fractions related to base-2 logarithms of 5. Since 2 and 10 appear
in this problem, there's probably a reason for the correlation.

*/

package main

import (
	"fmt"
	"math"
	"math/big"
	"math/bits"
	"os"
	"runtime/pprof"
	"time"
)

func main() {
	// constants
	const resolution = 1 << 20 // see below
	var (
		res      = big.NewInt(resolution)
		log2of10 = math.Log2(10)
		one      = big.NewInt(1)
		ten      = big.NewInt(10)
	)

	// loop state
	var (
		tensPower = new(big.Int)
		twosPower = new(big.Int)
		diff      = new(big.Int)
		estimator = new(big.Int)
		bestEst   = new(big.Int)
		bestError = big.NewRat(1, 1)
	)

	// consider each power of ten in a loop
	tensPower.SetUint64(1)
	for tensExp := 1; ; tensExp++ {
		tensPower.Mul(tensPower, ten)

		// see TestLog2Rounding for some justification for rounding
		twosExp := int(math.Round(float64(tensExp) * log2of10))

		// we'll make a big.Int called twosPower = 2^twosExp.
		// originally this read:
		//
		//   twosPower.Exp(big.NewInt(2), big.NewInt(int64(twosExp)), nil)
		//
		// but that Exp was found to dominate runtime after taming Rat's
		// performance (see below).
		//
		// creating a big.Int with a power of two shouldn't be expensive; it's
		// binary on the inside so we just need to find the right bit to set.
		// Int.SetBits is the tool for this: since a a big.Int is mostly a
		// little-endian word slice (see TestBigIntPowerOfTwoBits), let's cheat
		// by tacking a high quadword onto some padding zeros.
		zeros, mod := twosExp/bits.UintSize, twosExp%bits.UintSize
		bits := append(make([]big.Word, zeros, zeros+1), 1<<mod)
		twosPower.SetBits(bits)

		// now compute the error:
		//
		//   error = abs(twosPower-tensPower)/tensPower
		//
		// the linear part is cheap and easy:
		//
		//   diff = abs(twosPower - tensPower)
		diff.Sub(twosPower, tensPower).Abs(diff)

		// but the division is expensive.
		//
		// just creating a Rat via SetFrac was the first bottleneck this
		// program encountered, so that had to be abandoned in the fast path.
		// but we don't really need its infinite precision here: we're just
		// trying to see if one ratio is smaller than our record so far, and we
		// can trade some precision for speed if we're careful. by falling back
		// to Rat to test plausible candidates we get the best of both worlds.
		//
		// instead of computing Rat.Float64 and tracking it across iterations,
		// we'll use Int.Div to approximate the ratio. but that has two issues:
		//
		// 1) it discards its remainder, losing precision
		// 2) it returns an integer so our sub-unity error needs to be inverted
		//
		// to minimize the amount of precision lost by truncation, we'll
		// multiply by a constant first, at the expense of costlier bigint
		// math.
		//
		// TODO: tune the constant via benchmarking.
		//
		// we define the "estimator" as an integer that increases inversely
		// with error and is scaled by a constant called resolution. in order
		// to conservatively underestimate the true error we'll round up (since
		// it's inverted):
		//
		//   estimator =    ceil(resolution/error)
		//             =    ceil(resolution/(diff/tensPower))
		//             =    ceil(resolution*tensPower/diff)
		//             = 1+floor(resolution*tensPower/diff)
		//           --> Add(1, Div(Mul(res, tensPower), diff))
		estimator.Mul(res, tensPower).Div(estimator, diff).Add(estimator, one)

		if estimator.Cmp(bestEst) >= 0 {
			bestEst.Set(estimator)
			// in case we did underestimate and this isn't actually a new
			// record, double check with full precision
			actualError := new(big.Rat).SetFrac(diff, tensPower)
			if actualError.Cmp(bestError) < 0 {
				bestError.Set(actualError)
				e, _ := bestError.Float64()
				fmt.Printf("10^%d ~ 2^%d (%.g%%)\n", tensExp, twosExp, 100*e)
			}
		}
	}
}

func init() {
	// analyze with `go get github.com/google/pprof && pprof -web riddler.pprof`
	fh, err := os.Create("riddler.pprof")
	if err == nil {
		pprof.StartCPUProfile(fh)
		go func() {
			time.Sleep(60 * time.Second)
			pprof.StopCPUProfile()
			fh.Close()
		}()
	}
}
