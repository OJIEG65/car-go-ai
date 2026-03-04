package sim

import (
	"math"
	"testing"
)

func TestNewCarDummy(t *testing.T) {
	c := NewCar(100, 0, 30, 50, CarTypeDummy, 2, 0.2, 0.05, 0.03, nil, 0, 0, 0)
	if c.Sensor != nil {
		t.Error("dummy car should not have sensor")
	}
	if c.Brain != nil {
		t.Error("dummy car should not have brain")
	}
	if !c.Controls.Forward {
		t.Error("dummy car should have forward=true")
	}
}

func TestNewCarAI(t *testing.T) {
	c := NewCar(100, 100, 30, 50, CarTypeAI, 3, 0.2, 0.05, 0.03, []int{5, 6, 8, 4}, 5, 150, math.Pi/2)
	if c.Sensor == nil {
		t.Error("AI car should have sensor")
	}
	if c.Brain == nil {
		t.Error("AI car should have brain")
	}
	if !c.UseBrain {
		t.Error("AI car should use brain")
	}
}

func TestCarPolygon(t *testing.T) {
	c := NewCar(100, 100, 30, 50, CarTypeDummy, 2, 0.2, 0.05, 0.03, nil, 0, 0, 0)
	if len(c.Polygon) != 4 {
		t.Fatalf("polygon has %d points, want 4", len(c.Polygon))
	}
	// Polygon should be centered roughly around (100, 100)
	avgX, avgY := 0.0, 0.0
	for _, p := range c.Polygon {
		avgX += p.X
		avgY += p.Y
	}
	avgX /= 4
	avgY /= 4
	if math.Abs(avgX-100) > 1e-9 || math.Abs(avgY-100) > 1e-9 {
		t.Errorf("polygon center at (%.2f, %.2f), want (100, 100)", avgX, avgY)
	}
}

func TestCarMoveForward(t *testing.T) {
	c := NewCar(100, 100, 30, 50, CarTypeDummy, 2, 0.2, 0.05, 0.03, nil, 0, 0, 0)
	startY := c.Y
	// Simulate a few ticks of forward movement (no borders, no traffic)
	for i := 0; i < 10; i++ {
		c.move()
	}
	if c.Y >= startY {
		t.Errorf("car should move forward (decrease Y), got Y=%v from startY=%v", c.Y, startY)
	}
}

func TestCarDamageOnBorder(t *testing.T) {
	road := NewRoad(100, 200, 3)
	borders := road.BorderSegments()

	// Place car right at the left border
	c := NewCar(road.Left-5, 0, 30, 50, CarTypeDummy, 2, 0.2, 0.05, 0.03, nil, 0, 0, 0)
	c.Polygon = c.createPolygon()
	damaged := c.assessDamage(borders, nil)
	if !damaged {
		t.Error("car at border edge should be damaged")
	}
}
