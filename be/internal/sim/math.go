package sim

import "github.com/OJIEG65/car-go-ai/be/internal/types"

// Lerp performs linear interpolation between A and B at parameter t.
func Lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

// GetIntersection finds the intersection point of segments AB and CD.
// Returns nil if the segments do not intersect.
func GetIntersection(a, b, c, d types.Point) *types.Intersection {
	tTop := (d.X-c.X)*(a.Y-c.Y) - (d.Y-c.Y)*(a.X-c.X)
	uTop := (c.Y-a.Y)*(a.X-b.X) - (c.X-a.X)*(a.Y-b.Y)
	bottom := (d.Y-c.Y)*(b.X-a.X) - (d.X-c.X)*(b.Y-a.Y)

	if bottom != 0 {
		t := tTop / bottom
		u := uTop / bottom
		if t >= 0 && t <= 1 && u >= 0 && u <= 1 {
			return &types.Intersection{
				X:      Lerp(a.X, b.X, t),
				Y:      Lerp(a.Y, b.Y, t),
				Offset: t,
			}
		}
	}
	return nil
}

// PolysIntersect tests if two polygons share any edge intersection.
func PolysIntersect(poly1, poly2 types.Polygon) bool {
	for i := 0; i < len(poly1); i++ {
		for j := 0; j < len(poly2); j++ {
			touch := GetIntersection(
				poly1[i],
				poly1[(i+1)%len(poly1)],
				poly2[j],
				poly2[(j+1)%len(poly2)],
			)
			if touch != nil {
				return true
			}
		}
	}
	return false
}
