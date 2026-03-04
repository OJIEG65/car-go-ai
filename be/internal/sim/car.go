package sim

import (
	"math"

	"github.com/OJIEG65/car-go-ai/be/internal/nn"
	"github.com/OJIEG65/car-go-ai/be/internal/types"
)

// CarType distinguishes AI-controlled cars from dummy traffic.
type CarType int

const (
	CarTypeAI    CarType = iota
	CarTypeDummy
)

// Controls holds the current input state for a car.
type Controls struct {
	Forward bool
	Left    bool
	Right   bool
	Reverse bool
}

// Car represents a vehicle in the simulation.
type Car struct {
	X            float64
	Y            float64
	Width        float64
	Height       float64
	Speed        float64
	Acceleration float64
	MaxSpeed     float64
	Friction     float64
	TurnRate     float64
	Angle        float64
	Damaged      bool
	Polygon      types.Polygon
	Controls     Controls
	Sensor       *Sensor
	Brain        *nn.NeuralNetwork
	UseBrain     bool
	CarType      CarType
	StartY       float64
	TicksAlive   int
}

// NewCar creates a car with the given parameters.
func NewCar(x, y, width, height float64, carType CarType, maxSpeed, accel, friction, turnRate float64, neuronCounts []int, sensorCount int, sensorLength, sensorSpread float64) *Car {
	c := &Car{
		X:            x,
		Y:            y,
		Width:        width,
		Height:       height,
		Speed:        0,
		Acceleration: accel,
		MaxSpeed:     maxSpeed,
		Friction:     friction,
		TurnRate:     turnRate,
		Angle:        0,
		Damaged:      false,
		UseBrain:     carType == CarTypeAI,
		CarType:      carType,
		StartY:       y,
	}

	if carType != CarTypeDummy {
		c.Sensor = NewSensor(sensorCount, sensorLength, sensorSpread)
		c.Brain = nn.NewNetwork(neuronCounts)
	}

	if carType == CarTypeDummy {
		c.Controls.Forward = true
	}

	// Initial polygon
	c.Polygon = c.createPolygon()
	return c
}

// Update advances the car one tick.
func (c *Car) Update(roadBorders []types.Segment, traffic []*Car) {
	if !c.Damaged {
		c.move()
		c.Polygon = c.createPolygon()
		c.Damaged = c.assessDamage(roadBorders, traffic)
	}
	if c.Sensor != nil {
		c.Sensor.Update(c.X, c.Y, c.Angle, roadBorders, traffic)
		offsets := c.Sensor.Offsets()
		outputs := c.Brain.FeedForward(offsets)
		if c.UseBrain {
			c.Controls.Forward = outputs[0] > 0.5
			c.Controls.Left = outputs[1] > 0.5
			c.Controls.Right = outputs[2] > 0.5
			c.Controls.Reverse = outputs[3] > 0.5
		}
	}
	if !c.Damaged {
		c.TicksAlive++
	}
}

func (c *Car) move() {
	if c.Controls.Forward {
		c.Speed += c.Acceleration
	}
	if c.Controls.Reverse {
		c.Speed -= c.Acceleration
	}

	if c.Speed > c.MaxSpeed {
		c.Speed = c.MaxSpeed
	}
	if c.Speed < -c.MaxSpeed/2 {
		c.Speed = -c.MaxSpeed / 2
	}

	if c.Speed > 0 {
		c.Speed -= c.Friction
	}
	if c.Speed < 0 {
		c.Speed += c.Friction
	}

	if c.Speed != 0 {
		flip := 1.0
		if c.Speed < 0 {
			flip = -1.0
		}

		if c.Controls.Left {
			c.Angle += c.TurnRate * flip
		}
		if c.Controls.Right {
			c.Angle -= c.TurnRate * flip
		}
	}

	c.X -= math.Sin(c.Angle) * c.Speed
	c.Y -= math.Cos(c.Angle) * c.Speed
}

func (c *Car) createPolygon() types.Polygon {
	points := make(types.Polygon, 4)
	rad := math.Hypot(c.Width, c.Height) / 2
	alpha := math.Atan2(c.Width, c.Height)

	points[0] = types.Point{
		X: c.X - math.Sin(c.Angle-alpha)*rad,
		Y: c.Y - math.Cos(c.Angle-alpha)*rad,
	}
	points[1] = types.Point{
		X: c.X - math.Sin(c.Angle+alpha)*rad,
		Y: c.Y - math.Cos(c.Angle+alpha)*rad,
	}
	points[2] = types.Point{
		X: c.X - math.Sin(math.Pi+c.Angle-alpha)*rad,
		Y: c.Y - math.Cos(math.Pi+c.Angle-alpha)*rad,
	}
	points[3] = types.Point{
		X: c.X - math.Sin(math.Pi+c.Angle+alpha)*rad,
		Y: c.Y - math.Cos(math.Pi+c.Angle+alpha)*rad,
	}
	return points
}

func (c *Car) assessDamage(roadBorders []types.Segment, traffic []*Car) bool {
	for i := range roadBorders {
		borderPoly := types.Polygon{roadBorders[i][0], roadBorders[i][1]}
		if PolysIntersect(c.Polygon, borderPoly) {
			return true
		}
	}
	for _, other := range traffic {
		if other == c || other.Polygon == nil {
			continue
		}
		if PolysIntersect(c.Polygon, other.Polygon) {
			return true
		}
	}
	return false
}

// Distance returns how far the car has traveled from its start (lower Y = further).
func (c *Car) Distance() float64 {
	return c.StartY - c.Y
}
