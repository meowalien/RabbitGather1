package math

import "math"

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
