package util

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

// haversin(θ) function
func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

// Distance function returns the distance (in meters) between two points of
//     a given longitude and latitude relatively accurately (using a spherical
//     approximation of the Earth) through the Haversin Distance Formula for
//     great arc distance on a sphere with accuracy for small distances
//
// point coordinates are supplied in degrees and converted into rad. in the func
//
// distance returned is METERS!!!!!!
// http://en.wikipedia.org/wiki/Haversine_formula
func Distance(lat1, lon1, lat2, lon2 float64) float64 {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2, r float64
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
}
