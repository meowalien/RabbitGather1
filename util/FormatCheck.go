package util

// Is_WGS84_2D check if the input comply with WGS84_2D format
func Is_WGS84_2D(Long float64, Lat float64) bool {
	return (Lat > -90 || Lat < 90) && (Long > -180 || Long < 180)
}
