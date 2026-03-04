package nn

import "math/rand"

// Level represents a single layer in the neural network.
type Level struct {
	Inputs  []float64   `json:"inputs"`
	Outputs []float64   `json:"outputs"`
	Biases  []float64   `json:"biases"`
	Weights [][]float64 `json:"weights"`
}

// NewLevel creates a randomly initialized level with the given dimensions.
func NewLevel(inputCount, outputCount int) *Level {
	l := &Level{
		Inputs:  make([]float64, inputCount),
		Outputs: make([]float64, outputCount),
		Biases:  make([]float64, outputCount),
		Weights: make([][]float64, inputCount),
	}
	for i := range l.Weights {
		l.Weights[i] = make([]float64, outputCount)
	}
	l.randomize()
	return l
}

func (l *Level) randomize() {
	for i := range l.Weights {
		for j := range l.Weights[i] {
			l.Weights[i][j] = rand.Float64()*2 - 1
		}
	}
	for i := range l.Biases {
		l.Biases[i] = rand.Float64()*2 - 1
	}
}

// FeedForward computes step-activation outputs for the given inputs.
// Uses binary step function: output = 1 if weighted sum > bias, else 0.
func (l *Level) FeedForward(givenInputs []float64) []float64 {
	copy(l.Inputs, givenInputs)

	for i := range l.Outputs {
		sum := 0.0
		for j := range l.Inputs {
			sum += l.Inputs[j] * l.Weights[j][i]
		}
		if sum > l.Biases[i] {
			l.Outputs[i] = 1
		} else {
			l.Outputs[i] = 0
		}
	}
	return l.Outputs
}
