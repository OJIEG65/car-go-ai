package sim

import "github.com/OJIEG65/car-go-ai/be/internal/types"

const roadInfinity = 1000000

// Road represents the driving surface with lane borders.
type Road struct {
	X         float64
	Width     float64 // narrowed width (90% of original)
	LaneCount int
	Left      float64
	Right     float64
	Top       float64
	Bottom    float64
	Borders   [2]types.Segment
}

// NewRoad creates a road centered at x with the given total width and lane count.
// Width is narrowed to 90% to provide visual margins, matching the JS implementation.
func NewRoad(x, width float64, laneCount int) *Road {
	narrowed := width * 0.9
	left := x - narrowed/2
	right := x + narrowed/2

	r := &Road{
		X:         x,
		Width:     narrowed,
		LaneCount: laneCount,
		Left:      left,
		Right:     right,
		Top:       roadInfinity,
		Bottom:    -roadInfinity,
	}

	r.Borders[0] = types.Segment{
		{X: left, Y: r.Top},
		{X: left, Y: r.Bottom},
	}
	r.Borders[1] = types.Segment{
		{X: right, Y: r.Top},
		{X: right, Y: r.Bottom},
	}
	return r
}

// GetLaneCenter returns the horizontal center of the given lane index.
func (r *Road) GetLaneCenter(laneIndex int) float64 {
	laneWidth := r.Width / float64(r.LaneCount)
	idx := laneIndex
	if idx > r.LaneCount-1 {
		idx = r.LaneCount - 1
	}
	return r.Left + laneWidth/2 + float64(idx)*laneWidth
}

// BorderSegments returns the borders as a slice for sensor intersection tests.
func (r *Road) BorderSegments() []types.Segment {
	return r.Borders[:]
}
