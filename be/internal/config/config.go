package config

// SimulationConfig holds all tunable simulation parameters.
type SimulationConfig struct {
	CarCount      int     `json:"carCount"`
	LaneCount     int     `json:"laneCount"`
	RoadWidth     float64 `json:"roadWidth"`
	RoadCenterX   float64 `json:"roadCenterX"`
	TickRate      int     `json:"tickRate"`
	MutationRate  float64 `json:"mutationRate"`
	MaxSpeed      float64 `json:"maxSpeed"`
	Acceleration  float64 `json:"acceleration"`
	Friction      float64 `json:"friction"`
	TurnRate      float64 `json:"turnRate"`
	SensorCount   int     `json:"sensorCount"`
	SensorLength  float64 `json:"sensorLength"`
	SensorSpread  float64 `json:"sensorSpread"`
	CarWidth      float64 `json:"carWidth"`
	CarHeight     float64 `json:"carHeight"`
	NeuronCounts  []int   `json:"neuronCounts"`
}

// DefaultConfig returns sensible defaults matching the original JS simulation.
func DefaultConfig() SimulationConfig {
	return SimulationConfig{
		CarCount:     100,
		LaneCount:    3,
		RoadWidth:    200,
		RoadCenterX:  100,
		TickRate:     60,
		MutationRate: 0.1,
		MaxSpeed:     3,
		Acceleration: 0.2,
		Friction:     0.05,
		TurnRate:     0.03,
		SensorCount:  5,
		SensorLength: 150,
		SensorSpread: 1.5707963267948966, // math.Pi / 2
		CarWidth:     30,
		CarHeight:    50,
		NeuronCounts: []int{5, 6, 8, 4}, // sensor inputs → hidden → hidden → controls
	}
}
