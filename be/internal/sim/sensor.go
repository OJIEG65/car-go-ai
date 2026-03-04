package sim

import (
	"math"

	"github.com/OJIEG65/car-go-ai/be/internal/types"
)

// SensorReading stores a single ray's hit result plus the ray endpoints.
type SensorReading struct {
	Hit *types.Intersection
	Ray types.Segment
}

// Sensor casts rays from a car to detect obstacles.
type Sensor struct {
	RayCount  int
	RayLength float64
	RaySpread float64
	Rays      []types.Segment
	Readings  []*types.Intersection
}

// NewSensor creates a sensor matching the JS defaults.
func NewSensor(rayCount int, rayLength, raySpread float64) *Sensor {
	return &Sensor{
		RayCount:  rayCount,
		RayLength: rayLength,
		RaySpread: raySpread,
		Rays:      make([]types.Segment, rayCount),
		Readings:  make([]*types.Intersection, rayCount),
	}
}

// Update recasts all rays from the car's position and finds readings.
func (s *Sensor) Update(carX, carY, carAngle float64, roadBorders []types.Segment, traffic []*Car) {
	s.castRays(carX, carY, carAngle)
	for i := 0; i < s.RayCount; i++ {
		s.Readings[i] = s.getReading(s.Rays[i], roadBorders, traffic)
	}
}

func (s *Sensor) castRays(carX, carY, carAngle float64) {
	for i := 0; i < s.RayCount; i++ {
		var t float64
		if s.RayCount == 1 {
			t = 0.5
		} else {
			t = float64(i) / float64(s.RayCount-1)
		}
		rayAngle := Lerp(s.RaySpread/2, -s.RaySpread/2, t) + carAngle

		start := types.Point{X: carX, Y: carY}
		end := types.Point{
			X: carX - math.Sin(rayAngle)*s.RayLength,
			Y: carY - math.Cos(rayAngle)*s.RayLength,
		}
		s.Rays[i] = types.Segment{start, end}
	}
}

func (s *Sensor) getReading(ray types.Segment, roadBorders []types.Segment, traffic []*Car) *types.Intersection {
	var best *types.Intersection

	for i := range roadBorders {
		touch := GetIntersection(ray[0], ray[1], roadBorders[i][0], roadBorders[i][1])
		if touch != nil && (best == nil || touch.Offset < best.Offset) {
			best = touch
		}
	}

	for _, car := range traffic {
		if car.Polygon == nil {
			continue
		}
		poly := car.Polygon
		for j := 0; j < len(poly); j++ {
			touch := GetIntersection(
				ray[0], ray[1],
				poly[j], poly[(j+1)%len(poly)],
			)
			if touch != nil && (best == nil || touch.Offset < best.Offset) {
				best = touch
			}
		}
	}

	return best
}

// Offsets returns normalized sensor values: 0 = nothing detected, 1 = obstacle touching car.
func (s *Sensor) Offsets() []float64 {
	offsets := make([]float64, s.RayCount)
	for i, r := range s.Readings {
		if r == nil {
			offsets[i] = 0
		} else {
			offsets[i] = 1 - r.Offset
		}
	}
	return offsets
}
