/*
https://fivethirtyeight.com/features/can-you-get-another-haircut-already/

After 2^10, what’s the next (whole number) power of 2 that comes closer to a
power of 10? (To be clear, “closer” doesn’t refer to the absolute difference —
it means your power of 2 should differ from a power of 10 by less than 2.4
percent.)
*/

package main

import (
	"fmt"
	"math"
	"math/big"
	"os"
	"runtime/pprof"
	"time"
)

func main() {
	fh, err := os.Create("riddler.pprof")
	if err == nil {
		pprof.StartCPUProfile(fh)
		go func() {
			time.Sleep(30 * time.Second)
			pprof.StopCPUProfile()
			fh.Close()
		}()
	}
	var (
		two, ten         = big.NewInt(2), big.NewInt(10)
		pTen, pTwo, diff big.Int
		rat              big.Rat
		log2of10         = math.Log2(10)
		bestErr          = 1.
	)
	pTen.Set(ten)
	for tensExp := 1; ; tensExp++ {
		// compare the power of ten to the powers of two below and above it
		twosExp := float64(tensExp) * log2of10
		lowerTwosExp := int(math.Floor(twosExp))
		for i, twosExp := range [...]int{lowerTwosExp, lowerTwosExp + 1} {
			if i == 0 {
				pTwo.Exp(two, big.NewInt(int64(twosExp)), nil)
			} else {
				pTwo.Mul(&pTwo, two)
			}
			diff.Sub(&pTwo, &pTen).Abs(&diff)
			rat.SetFrac(&diff, &pTen)
			if e, _ := rat.Float64(); e < bestErr {
				bestErr = e
				fmt.Printf("10**%d ~ 2**%d (%.2g%%)\n", tensExp, twosExp, 100*e)
			}
		}
		pTen.Mul(&pTen, ten)
	}
}
