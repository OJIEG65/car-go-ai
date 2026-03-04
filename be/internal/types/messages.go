package types

import "encoding/json"

// WSMessage is the envelope for all WebSocket messages.
type WSMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

// InitPayload is sent once on WebSocket connect.
type InitPayload struct {
	Road   json.RawMessage `json:"road"`
	Config json.RawMessage `json:"config"`
}

// GenerationEndPayload is sent when all cars are damaged.
type GenerationEndPayload struct {
	Generation   int     `json:"generation"`
	BestFitness  float64 `json:"bestFitness"`
	AvgFitness   float64 `json:"avgFitness"`
	BestDistance float64 `json:"bestDistance"`
}

// BrainPayload carries a serialized neural network.
type BrainPayload struct {
	Name string          `json:"name,omitempty"`
	Data json.RawMessage `json:"data"`
}

// ErrorPayload carries error info.
type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
