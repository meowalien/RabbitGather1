package math

import (
	"errors"
	"math"
)

func CutIntMax(target, max int64) int64 {
	return CutIntBetween(target, 1, max)
}
func CutIntBetween(target, min, max int64) int64 {
	if min < 1 {
		panic("min must >= 1")
	}
	return target / int64(math.Pow(10, float64(min-1))) % int64(math.Pow(10, float64(max-min+1)))
}

func Round(x float64) int {
	return int(math.Floor(x + 0.5))
}

func IntLength(a int) int {
	count := 0
	for a != 0 {
		a /= 10
		count++
	}
	return count
}

var ErrOverflow = errors.New("integer overflow")

func Add32(left, right int32) (int32, error) {
	if right > 0 {
		if left > math.MaxInt32-right {
			return 0, ErrOverflow
		}
	} else {
		if left < math.MinInt32-right {
			return 0, ErrOverflow
		}
	}
	return left + right, nil
}


// return true if success (not Overflow)
func Add64(left, right int64) (int64, bool) {
	if right > 0 {
		if left > math.MaxInt64-right {
			return 0, false
		}
	} else {
		if left < math.MinInt64-right {
			return 0, false
		}
	}
	return left + right, true
}