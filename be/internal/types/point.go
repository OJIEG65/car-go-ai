package types

// Point represents a 2D coordinate.
type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Segment is a line segment defined by two points.
type Segment [2]Point

// Intersection holds a ray-segment hit result.
type Intersection struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Offset float64 `json:"offset"`
}

// Polygon is a slice of points forming a closed shape.
type Polygon []Point
