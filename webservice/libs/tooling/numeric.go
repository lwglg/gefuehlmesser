package tooling

import (
	"fmt"
	"math"
	"strconv"
)

// Both functions adapted from: https://stackoverflow.com/a/77009270/29040254

// Precision controls the number of digits (excluding the exponent)
// Precision of -1 uses the smallest number of digits
func TruncFloat(f float64, precision int) (float64, error) {
	floatBits := 64

	if math.IsNaN(f) || math.IsInf(f, 1) || math.IsInf(f, -1) {
		return 0, fmt.Errorf("bad float val %f", f)
	}

	fTruncStr := strconv.FormatFloat(f, 'f', precision+1, floatBits)
	fTruncStr = fTruncStr[:len(fTruncStr)-1]
	fTrunc, err := strconv.ParseFloat(fTruncStr, floatBits)
	if err != nil {
		return 0, err
	}

	return fTrunc, nil
}

// Precision controls the number of digits (excluding the exponent)
// Precision of -1 uses the smallest number of digits
func RoundFloat(f float64, precision int) (float64, error) {
	floatBits := 64

	if math.IsNaN(f) || math.IsInf(f, 1) || math.IsInf(f, -1) {
		return 0, fmt.Errorf("bad float val %f", f)
	}

	fRoundedStr := strconv.FormatFloat(f, 'f', precision, floatBits)
	fRounded, err := strconv.ParseFloat(fRoundedStr, floatBits)
	if err != nil {
		return 0, err
	}

	return fRounded, nil
}

func RoundFloatFast(f float64, precision int) (float64, error) {
	mul := math.Pow10(precision)
	if mul == 0 {
		return 0, nil
	}

	product := f * mul
	var roundingErr error
	if product > float64(math.MaxInt64) {
		roundingErr = fmt.Errorf("unsafe round: float64=%+v, places=%d", f, precision)
	}

	return math.Round(product) / mul, roundingErr
}

func IsPrime(n int64) bool {
	if n <= 1 || n%2 == 0 {
		return false
	}

	if n == 2 {
		return true
	}

	sqrt := int(math.Sqrt(float64(n)))
	for i := 3; i <= sqrt; i += 2 {
		if n%int64(i) == 0 {
			return false
		}
	}
	return true
}
