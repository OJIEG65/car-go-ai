package nn

// NeuralNetwork is a feed-forward network with step activation.
type NeuralNetwork struct {
	Levels []*Level `json:"levels"`
}

// NewNetwork creates a neural network with the given layer sizes.
// e.g. [5, 6, 8, 4] → 3 levels: 5→6, 6→8, 8→4
func NewNetwork(neuronCounts []int) *NeuralNetwork {
	n := &NeuralNetwork{
		Levels: make([]*Level, len(neuronCounts)-1),
	}
	for i := 0; i < len(neuronCounts)-1; i++ {
		n.Levels[i] = NewLevel(neuronCounts[i], neuronCounts[i+1])
	}
	return n
}

// FeedForward propagates inputs through all levels and returns the final outputs.
func (n *NeuralNetwork) FeedForward(givenInputs []float64) []float64 {
	outputs := n.Levels[0].FeedForward(givenInputs)
	for i := 1; i < len(n.Levels); i++ {
		outputs = n.Levels[i].FeedForward(outputs)
	}
	return outputs
}
