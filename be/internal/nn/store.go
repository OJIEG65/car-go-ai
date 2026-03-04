package nn

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Store handles saving and loading brains to/from JSON files.
type Store struct {
	Dir string
}

// NewStore creates a store at the given directory, creating it if needed.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create brain dir: %w", err)
	}
	return &Store{Dir: dir}, nil
}

// Save writes a brain to a named JSON file.
func (s *Store) Save(name string, brain *NeuralNetwork) error {
	path := filepath.Join(s.Dir, name+".json")
	data, err := json.MarshalIndent(brain, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal brain: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write brain file: %w", err)
	}
	return nil
}

// Load reads a brain from a named JSON file.
func (s *Store) Load(name string) (*NeuralNetwork, error) {
	path := filepath.Join(s.Dir, name+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read brain file: %w", err)
	}
	var brain NeuralNetwork
	if err := json.Unmarshal(data, &brain); err != nil {
		return nil, fmt.Errorf("unmarshal brain: %w", err)
	}
	return &brain, nil
}

// List returns all saved brain names.
func (s *Store) List() ([]string, error) {
	entries, err := os.ReadDir(s.Dir)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			names = append(names, e.Name()[:len(e.Name())-5])
		}
	}
	return names, nil
}

// Delete removes a saved brain file.
func (s *Store) Delete(name string) error {
	path := filepath.Join(s.Dir, name+".json")
	return os.Remove(path)
}
