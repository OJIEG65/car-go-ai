package sim

import (
	"math"
	"testing"

	"github.com/OJIEG65/car-go-ai/be/internal/types"
)

func TestSensorCastRays(t *testing.T) {
	s := NewSensor(5, 150, math.Pi/2)
	s.castRays(100, 100, 0)

	if len(s.Rays) != 5 {
		t.Fatalf("expected 5 rays, got %d", len(s.Rays))
	}

	// All rays should start at (100, 100)
	for i, ray := range s.Rays {
		if math.Abs(ray[0].X-100) > 1e-9 || math.Abs(ray[0].Y-100) > 1e-9 {
			t.Errorf("ray[%d] start at (%v, %v), want (100, 100)", i, ray[0].X, ray[0].Y)
		}
		// Ray length should be ~150
		dx := ray[1].X - ray[0].X
		dy := ray[1].Y - ray[0].Y
		length := math.Hypot(dx, dy)
		if math.Abs(length-150) > 1e-6 {
			t.Errorf("ray[%d] length = %v, want 150", i, length)
		}
	}

	// Middle ray (index 2) should point straight forward (negative Y direction)
	middleRay := s.Rays[2]
	dx := middleRay[1].X - middleRay[0].X
	dy := middleRay[1].Y - middleRay[0].Y
	if math.Abs(dx) > 1e-6 {
		t.Errorf("middle ray dx = %v, want ~0", dx)
	}
	if dy >= 0 {
		t.Error("middle ray should point in negative Y direction")
	}
}

func TestSensorReadingHitsBorder(t *testing.T) {
	s := NewSensor(5, 150, math.Pi/2)

	// Place a border very close in front
	border := types.Segment{
		{X: 0, Y: 50},
		{X: 200, Y: 50},
	}
	borders := []types.Segment{border}

	s.Update(100, 100, 0, borders, nil)

	// The forward-facing ray should hit the border
	// Ray goes from (100, 100) forward (negative Y), border at Y=50
	hitFound := false
	for _, r := range s.Readings {
		if r != nil {
			hitFound = true
			break
		}
	}
	if !hitFound {
		t.Error("expected at least one sensor reading to hit the border")
	}
}

func TestSensorOffsets(t *testing.T) {
	s := NewSensor(3, 100, math.Pi/2)
	// No hits → all offsets should be 0
	s.Readings = make([]*types.Intersection, 3)
	offsets := s.Offsets()
	for i, o := range offsets {
		if o != 0 {
			t.Errorf("offset[%d] = %v, want 0 for nil reading", i, o)
		}
	}

	// Hit at offset 0.3 → sensor value = 0.7
	s.Readings[1] = &types.Intersection{Offset: 0.3}
	offsets = s.Offsets()
	if math.Abs(offsets[1]-0.7) > 1e-9 {
		t.Errorf("offset[1] = %v, want 0.7", offsets[1])
	}
}
