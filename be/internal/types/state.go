package types

// CarState is the serializable state of a single car sent to the frontend.
type CarState struct {
	X       float64 `json:"x"`
	Y       float64 `json:"y"`
	Angle   float64 `json:"angle"`
	Damaged bool    `json:"damaged"`
	Polygon Polygon `json:"polygon"`
	Speed   float64 `json:"speed"`
}

// SensorState holds one car's sensor data for rendering.
type SensorState struct {
	Rays     [][2]Point     `json:"rays"`
	Readings []*Intersection `json:"readings"`
}

// WorldState is the per-tick payload sent to the frontend.
type WorldState struct {
	Cars       []CarState   `json:"cars"`
	Traffic    []CarState   `json:"traffic"`
	BestIndex  int          `json:"bestIndex"`
	BestSensor *SensorState `json:"bestSensor,omitempty"`
	Tick       int          `json:"tick"`
	Generation int          `json:"generation"`
}
