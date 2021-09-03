package util

// Point2D represents a coordinate point in 2D
type Point2D struct {
	X float64 `json:"x"  form:"x"  binding:"required"`
	Y float64 `json:"y" form:"y"  binding:"required"`
}
