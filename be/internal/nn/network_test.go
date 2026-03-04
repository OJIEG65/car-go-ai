package nn

import (
	"testing"
)

func TestNewLevel(t *testing.T) {
	l := NewLevel(5, 4)
	if len(l.Inputs) != 5 {
		t.Errorf("inputs length = %d, want 5", len(l.Inputs))
	}
	if len(l.Outputs) != 4 {
		t.Errorf("outputs length = %d, want 4", len(l.Outputs))
	}
	if len(l.Biases) != 4 {
		t.Errorf("biases length = %d, want 4", len(l.Biases))
	}
	if len(l.Weights) != 5 {
		t.Errorf("weights rows = %d, want 5", len(l.Weights))
	}
	for i, w := range l.Weights {
		if len(w) != 4 {
			t.Errorf("weights[%d] length = %d, want 4", i, len(w))
		}
	}
	// Check that weights are in [-1, 1]
	for i, row := range l.Weights {
		for j, w := range row {
			if w < -1 || w > 1 {
				t.Errorf("weight[%d][%d] = %v, out of [-1, 1]", i, j, w)
			}
		}
	}
}

func TestFeedForwardStepActivation(t *testing.T) {
	l := NewLevel(2, 2)
	// Set known weights and biases
	l.Weights[0][0] = 1.0
	l.Weights[0][1] = -1.0
	l.Weights[1][0] = -1.0
	l.Weights[1][1] = 1.0
	l.Biases[0] = 0.0
	l.Biases[1] = 0.0

	outputs := l.FeedForward([]float64{1, 0})
	// output[0] = 1*1 + 0*(-1) = 1 > 0 → 1
	// output[1] = 1*(-1) + 0*1 = -1 > 0 → 0
	if outputs[0] != 1 {
		t.Errorf("output[0] = %v, want 1", outputs[0])
	}
	if outputs[1] != 0 {
		t.Errorf("output[1] = %v, want 0", outputs[1])
	}
}

func TestNewNetwork(t *testing.T) {
	n := NewNetwork([]int{5, 6, 8, 4})
	if len(n.Levels) != 3 {
		t.Fatalf("levels count = %d, want 3", len(n.Levels))
	}
	// Check dimensions
	if len(n.Levels[0].Inputs) != 5 || len(n.Levels[0].Outputs) != 6 {
		t.Error("level 0 dimensions wrong")
	}
	if len(n.Levels[1].Inputs) != 6 || len(n.Levels[1].Outputs) != 8 {
		t.Error("level 1 dimensions wrong")
	}
	if len(n.Levels[2].Inputs) != 8 || len(n.Levels[2].Outputs) != 4 {
		t.Error("level 2 dimensions wrong")
	}
}

func TestNetworkFeedForward(t *testing.T) {
	n := NewNetwork([]int{5, 6, 8, 4})
	inputs := []float64{0.5, 0.3, 0.0, 0.8, 0.1}
	outputs := n.FeedForward(inputs)

	if len(outputs) != 4 {
		t.Fatalf("output length = %d, want 4", len(outputs))
	}
	// All outputs should be 0 or 1 (step activation)
	for i, o := range outputs {
		if o != 0 && o != 1 {
			t.Errorf("output[%d] = %v, want 0 or 1", i, o)
		}
	}
}
