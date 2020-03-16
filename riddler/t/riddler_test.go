package t

import (
	"math"
	"testing"
)

func TestLog2Rounding(t *testing.T) {
	// assertion: rounding the log2 of a power of ten to get the power of two
	// is equivalent to the picking the lesser percentage error relative to the
	// tens power
	var log2of10 = math.Log2(10)
	for tensExp := 1; tensExp < 1000; tensExp++ {
		var (
			n            = math.Pow10(tensExp)
			twosExpLower = (log2of10 * math.Floor(float64(tensExp)))
			twosExpRound = (log2of10 * math.Round(float64(tensExp)))
			twosExpUpper = (log2of10 * math.Ceil(float64(tensExp)))
			nLower       = math.Exp2(twosExpLower)
			nRound       = math.Exp2(twosExpRound)
			nUpper       = math.Exp2(twosExpUpper)
			devLower     = math.Abs((twosExpLower - n) / n)
			devRound     = math.Abs((twosExpUpper - n) / n)
			devUpper     = math.Abs((twosExpUpper - n) / n)
			want         = math.Min(devLower, devUpper)
			got          = devRound
		)
		if math.IsInf(n, 0) {
			break
		}
		if want != got {
			t.Error(tensExp, n, nLower, nRound, nUpper)
		}
	}
}
