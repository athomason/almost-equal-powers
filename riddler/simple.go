/*
Problem: https://fivethirtyeight.com/features/can-you-get-another-haircut-already/

> After 2^10, what’s the next (whole number) power of 2 that comes closer to a
> power of 10? (To be clear, “closer” doesn’t refer to the absolute difference —
> it means your power of 2 should differ from a power of 10 by less than 2.4
> percent.)

Solution: for each power of ten, we'll compute the difference ratio to its closer power of two. If this ratio goes down, we have the next "closer" value.

Observation: the powers of ten turn out to have an interesting pattern: 1, 3, 28, 59, 146, 643, 4004, 8651, 12655, 21306, 76573, 97879, ... which appears as OEIS A046104 and A116984

*/

package main

import (
	"fmt"
	"math"
	"math/big"
)

var values = [...]int{
	1, 3, 28, 59, 146, 643, 4004, 8651, 12655, 21306, 76573, 97879, 1838395,
	1936274, 13456039, 15392313, 44240665, 59632978, 103873643, 475127550,
	579001193, 24793177656, 149338067129, 174131244785, 845863046269,
	1865857337323, 6443435058238, 8309292395561, 23062019849360,
	146681411491721, 169743431341081, 655911705514964, 2793390253400937,
	3449301958915901, 30387805924728145, 33837107883644046, 165736237459304329,
	199573345342948375, 564882928145201079, 1329339201633350533,
}

func main() {
	var (
		two      = big.NewInt(2)
		ten      = big.NewInt(10)
		pTen     = new(big.Int)
		pTwo     = new(big.Int)
		diff     = new(big.Int)
		rat      = new(big.Rat)
		log2of10 = math.Log2(10)
		bestErr  = 1.
	)
	for _, tensExp := range values {
		pTen.Exp(ten, big.NewInt(int64(tensExp)), nil)
		twosExp := int(math.Round(float64(tensExp) * log2of10))
		pTwo.Exp(two, big.NewInt(int64(twosExp)), nil)
		diff.Sub(pTwo, pTen).Abs(diff)
		rat.SetFrac(diff, pTwo)
		if e, _ := rat.Float64(); e < bestErr {
			bestErr = e
			fmt.Printf("10**%d ~ 2**%d (%.2g%%)\n", tensExp, twosExp, 100*e)
		}
	}
}
