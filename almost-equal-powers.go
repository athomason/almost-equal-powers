package main

import (
	"flag"
	"fmt"
	"math"
	"math/big"
)

func main() {
	var (
		minBase, maxBase, minPower, maxPower int
		tolerance                            float64
		primeBases, verbose, debug           bool
	)
	flag.IntVar(&minBase, "min-base", 2, "minimum base")
	flag.IntVar(&maxBase, "max-base", 100, "maximum base")
	flag.IntVar(&minPower, "min-power", 1, "minimum power")
	flag.IntVar(&maxPower, "max-power", 200, "maximum power")
	flag.Float64Var(&tolerance, "tolerance", 0.00001, "tolerance for search")
	flag.BoolVar(&primeBases, "prime-bases", false, "consider only prime bases")
	flag.BoolVar(&verbose, "verbose", false, "show all calculations")
	flag.BoolVar(&debug, "debug", false, "show debugging info")
	flag.Parse()

	lowTolerance, highTolerance := big.NewRat(0, 1), big.NewRat(0, 1)
	lowTolerance.SetFloat64(1 + tolerance)
	highTolerance.Inv(lowTolerance)

	bases := make([]*big.Int, 0)
	for i := minBase; i <= maxBase; i++ {
		base := big.NewInt(int64(i))
		if !primeBases || base.ProbablyPrime(16) {
			bases = append(bases, base)
		}
	}
	if debug {
		fmt.Printf("bases=%v\n", bases)
	}
	powers := make([]*big.Int, 0, maxPower-minPower+1)
	for i := minPower; i <= maxPower; i++ {
		powers = append(powers, big.NewInt(int64(i)))
	}

	type tuple struct {
		base, power *big.Int
	}
	hits, misses := 0, 0
	cache := make(map[tuple]*big.Int)
	pow := func(base, power *big.Int) *big.Int {
		t := tuple{base, power}
		if value, exists := cache[t]; exists {
			hits++
			return value
		}
		misses++
		value := big.NewInt(0).Exp(base, power, nil)
		cache[t] = value
		return value
	}

	oneInt := big.NewInt(1)

	for _, base1 := range bases {
		log1 := math.Log(float64(base1.Int64()))
		for _, base2 := range bases {
			// don't duplicate work
			if base1.Cmp(base2) >= 0 {
				continue
			}

			// only consider relative primes
			var gcd big.Int
			gcd.GCD(nil, nil, base1, base2)
			if gcd.Cmp(oneInt) != 0 {
				continue
			}

			for _, power1 := range powers {
				big1 := pow(base1, power1)
				log2 := math.Log(float64(base2.Int64()))

				power2 := float64(power1.Int64()) * log1 / log2
				power2Low := big.NewInt(int64(math.Floor(power2)))
				power2High := big.NewInt(int64(math.Ceil(power2)))

				big2Low := pow(base2, power2Low)
				big2High := pow(base2, power2High)

				if verbose {
					fmt.Printf("%v**%v (%v) < %v**%v (%v) < %v**%v (%v)\n",
						base2, power2Low, big2Low,
						base1, power1, big1,
						base2, power2High, big2High,
					)
				}

				ratioLow := big.NewRat(0, 1).SetFrac(big1, big2Low)
				ratioHigh := big.NewRat(0, 1).SetFrac(big1, big2High)

				if ratioLow.Cmp(lowTolerance) < 0 {
					rat, _ := ratioLow.Float64()
					fmt.Printf("%v**%v=%v ~~ %v**%v=%v (%.6f%%)\n",
						base1, power1, big1,
						base2, power2Low, big2Low,
						100*(rat-1),
					)
				}

				if ratioHigh.Cmp(highTolerance) > 0 {
					rat, _ := ratioHigh.Float64()
					fmt.Printf("%v**%v=%v ~~ %v**%v=%v (%.6f%%)\n",
						base1, power1, big1,
						base2, power2High, big2High,
						100*(rat-1),
					)
				}
			}
		}
	}

	if debug {
		fmt.Printf("Cache: %d/%d=%.1f%%\n", hits, misses, 100*float64(hits)/(float64(hits+misses)))
	}
}
