package sim

import (
	"math"
	"testing"

	"github.com/OJIEG65/car-go-ai/be/internal/types"
)

func TestLerp(t *testing.T) {
	tests := []struct {
		a, b, param, want float64
	}{
		{0, 10, 0, 0},
		{0, 10, 1, 10},
		{0, 10, 0.5, 5},
		{-1, 1, 0.5, 0},
		{5, 5, 0.7, 5},
	}
	for _, tt := range tests {
		got := Lerp(tt.a, tt.b, tt.param)
		if math.Abs(got-tt.want) > 1e-9 {
			t.Errorf("Lerp(%v, %v, %v) = %v, want %v", tt.a, tt.b, tt.param, got, tt.want)
		}
	}
}

func TestGetIntersection(t *testing.T) {
	// Two perpendicular segments crossing at (5, 5)
	a := types.Point{X: 0, Y: 5}
	b := types.Point{X: 10, Y: 5}
	c := types.Point{X: 5, Y: 0}
	d := types.Point{X: 5, Y: 10}

	hit := GetIntersection(a, b, c, d)
	if hit == nil {
		t.Fatal("expected intersection, got nil")
	}
	if math.Abs(hit.X-5) > 1e-9 || math.Abs(hit.Y-5) > 1e-9 {
		t.Errorf("intersection at (%v, %v), want (5, 5)", hit.X, hit.Y)
	}
	if math.Abs(hit.Offset-0.5) > 1e-9 {
		t.Errorf("offset = %v, want 0.5", hit.Offset)
	}
}

func TestGetIntersectionNoHit(t *testing.T) {
	// Parallel segments
	a := types.Point{X: 0, Y: 0}
	b := types.Point{X: 10, Y: 0}
	c := types.Point{X: 0, Y: 5}
	d := types.Point{X: 10, Y: 5}

	hit := GetIntersection(a, b, c, d)
	if hit != nil {
		t.Errorf("expected nil for parallel segments, got %v", hit)
	}
}

func TestPolysIntersect(t *testing.T) {
	// Two overlapping squares
	poly1 := types.Polygon{
		{X: 0, Y: 0}, {X: 10, Y: 0}, {X: 10, Y: 10}, {X: 0, Y: 10},
	}
	poly2 := types.Polygon{
		{X: 5, Y: 5}, {X: 15, Y: 5}, {X: 15, Y: 15}, {X: 5, Y: 15},
	}

	if !PolysIntersect(poly1, poly2) {
		t.Error("expected polygons to intersect")
	}

	// Non-overlapping squares
	poly3 := types.Polygon{
		{X: 20, Y: 20}, {X: 30, Y: 20}, {X: 30, Y: 30}, {X: 20, Y: 30},
	}
	if PolysIntersect(poly1, poly3) {
		t.Error("expected polygons not to intersect")
	}
}
