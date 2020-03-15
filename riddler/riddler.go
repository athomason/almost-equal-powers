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
	const multiplier = 1e12
	var (
		two, ten         = big.NewInt(2), big.NewInt(10)
		pTen, pTwo, diff big.Int
		invErr, bestErr  big.Int
		log2of10         = math.Log2(10)
		mult             = big.NewInt(multiplier)
	)
	pTen.Set(ten)
	for tensExp := 1; ; tensExp++ {
		// compare the power of ten to the powers of two below and above it
		twosExp := float64(tensExp) * log2of10
		lowerTwosExp := int(math.Floor(twosExp))
		for i, twosExp := range [...]int{lowerTwosExp, lowerTwosExp + 1} {
			// pTwo = 2**twosExp
			if i == 0 {
				pTwo.Exp(two, big.NewInt(int64(twosExp)), nil)
			} else {
				pTwo.Mul(&pTwo, two)
			}
			diff.Sub(&pTwo, &pTen).Abs(&diff) // diff = abs(pTwo-pTen)
			// err = diff/pTen => invErr = pTen/diff
			invErr.Mul(&pTen, mult)
			invErr.Div(&invErr, &diff)
			if invErr.Cmp(&bestErr) > 0 {
				bestErr = invErr
				//fmt.Printf("%v-%v=%v ie=%v\n", pTwo.String(), pTen.String(),
				//	diff.String(), invErr.String())
				e := multiplier / float64(invErr.Int64())
				fmt.Printf("10**%d ~ 2**%d (%.2g%%)\n", tensExp, twosExp, 100*e)
			}
		}
		pTen.Mul(&pTen, ten)
	}
}
