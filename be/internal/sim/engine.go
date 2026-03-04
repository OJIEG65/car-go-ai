package sim

import (
	"log"
	"sync"
	"time"

	"github.com/OJIEG65/car-go-ai/be/internal/config"
	"github.com/OJIEG65/car-go-ai/be/internal/nn"
	"github.com/OJIEG65/car-go-ai/be/internal/types"
)

// EngineState represents the simulation's current state for subscribers.
type EngineState struct {
	Cars       []*Car
	Traffic    []*Car
	BestCar    *Car
	Tick       int
	Generation int
}

// Subscriber receives state updates from the engine.
type Subscriber struct {
	Ch   chan EngineState
	done chan struct{}
}

// Engine runs the simulation tick loop.
type Engine struct {
	Config     config.SimulationConfig
	Road       *Road
	Cars       []*Car
	Traffic    []*Car
	BestCar    *Car
	Tick       int
	Generation int

	mu          sync.RWMutex
	subscribers map[*Subscriber]struct{}
	running     bool
	paused      bool
	stopCh      chan struct{}
	pauseCh     chan struct{}
	resumeCh    chan struct{}
}

// NewEngine creates a simulation engine from config.
func NewEngine(cfg config.SimulationConfig) *Engine {
	e := &Engine{
		Config:      cfg,
		subscribers: make(map[*Subscriber]struct{}),
		stopCh:      make(chan struct{}),
		pauseCh:     make(chan struct{}),
		resumeCh:    make(chan struct{}),
	}
	e.Reset()
	return e
}

// Reset reinitializes all cars and traffic for a new generation.
func (e *Engine) Reset() {
	e.Road = NewRoad(e.Config.RoadCenterX, e.Config.RoadWidth, e.Config.LaneCount)

	e.Cars = make([]*Car, e.Config.CarCount)
	spawnX := e.Road.GetLaneCenter(1)
	for i := range e.Cars {
		e.Cars[i] = NewCar(
			spawnX, 100,
			e.Config.CarWidth, e.Config.CarHeight,
			CarTypeAI,
			e.Config.MaxSpeed,
			e.Config.Acceleration,
			e.Config.Friction,
			e.Config.TurnRate,
			e.Config.NeuronCounts,
			e.Config.SensorCount,
			e.Config.SensorLength,
			e.Config.SensorSpread,
		)
	}

	e.Traffic = SpawnTraffic(e.Road, DefaultTraffic())
	e.Tick = 0
	e.updateBestCar()
}

// Subscribe adds a state listener. The returned Subscriber's Ch receives updates each tick.
func (e *Engine) Subscribe() *Subscriber {
	s := &Subscriber{
		Ch:   make(chan EngineState, 3),
		done: make(chan struct{}),
	}
	e.mu.Lock()
	e.subscribers[s] = struct{}{}
	e.mu.Unlock()
	return s
}

// Unsubscribe removes a listener.
func (e *Engine) Unsubscribe(s *Subscriber) {
	e.mu.Lock()
	delete(e.subscribers, s)
	e.mu.Unlock()
	close(s.done)
}

// Start begins the tick loop in a goroutine.
func (e *Engine) Start() {
	e.mu.Lock()
	if e.running {
		e.mu.Unlock()
		return
	}
	e.running = true
	e.mu.Unlock()

	go e.tickLoop()
	log.Printf("Engine started: %d cars, tick rate %d/s, generation %d",
		e.Config.CarCount, e.Config.TickRate, e.Generation)
}

// Stop halts the tick loop.
func (e *Engine) Stop() {
	e.mu.Lock()
	if !e.running {
		e.mu.Unlock()
		return
	}
	e.running = false
	e.mu.Unlock()
	close(e.stopCh)
}

// Pause pauses the simulation.
func (e *Engine) Pause() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.paused = true
}

// Resume unpauses the simulation.
func (e *Engine) Resume() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.paused = false
}

// IsPaused returns whether the engine is paused.
func (e *Engine) IsPaused() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.paused
}

// GetState returns a snapshot of the current engine state.
func (e *Engine) GetState() EngineState {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return EngineState{
		Cars:       e.Cars,
		Traffic:    e.Traffic,
		BestCar:    e.BestCar,
		Tick:       e.Tick,
		Generation: e.Generation,
	}
}

// SetBrainForAll replaces every AI car's brain with a clone of the given network.
func (e *Engine) SetBrainForAll(brain *nn.NeuralNetwork) {
	e.mu.Lock()
	defer e.mu.Unlock()
	for _, car := range e.Cars {
		if car.Brain != nil {
			car.Brain = cloneBrain(brain)
		}
	}
}

func (e *Engine) tickLoop() {
	ticker := time.NewTicker(time.Second / time.Duration(e.Config.TickRate))
	defer ticker.Stop()

	for {
		select {
		case <-e.stopCh:
			return
		case <-ticker.C:
			if e.IsPaused() {
				continue
			}
			e.tick()
			e.fanout()
		}
	}
}

func (e *Engine) tick() {
	e.mu.Lock()
	defer e.mu.Unlock()

	borders := e.Road.BorderSegments()

	// Update traffic
	for _, t := range e.Traffic {
		t.Update(borders, nil)
	}

	// Update AI cars — each car only writes to itself, safe for goroutines
	var wg sync.WaitGroup
	for _, car := range e.Cars {
		wg.Add(1)
		go func(c *Car) {
			defer wg.Done()
			c.Update(borders, e.Traffic)
		}(car)
	}
	wg.Wait()

	e.Tick++
	e.updateBestCar()

	// Check if all cars are damaged → end generation
	allDamaged := true
	for _, car := range e.Cars {
		if !car.Damaged {
			allDamaged = false
			break
		}
	}

	if allDamaged {
		e.endGeneration()
	}
}

func (e *Engine) updateBestCar() {
	if len(e.Cars) == 0 {
		return
	}
	best := e.Cars[0]
	for _, car := range e.Cars[1:] {
		if car.Y < best.Y {
			best = car
		}
	}
	e.BestCar = best
}

func (e *Engine) endGeneration() {
	e.Generation++
	log.Printf("Generation %d ended at tick %d, best distance: %.1f",
		e.Generation, e.Tick, e.BestCar.Distance())

	// For now, just reset with fresh random brains.
	// Phase 4 will add evolution (mutation + selection) here.
	e.resetCars()
}

func (e *Engine) resetCars() {
	spawnX := e.Road.GetLaneCenter(1)
	for _, car := range e.Cars {
		car.X = spawnX
		car.Y = 100
		car.StartY = 100
		car.Speed = 0
		car.Angle = 0
		car.Damaged = false
		car.TicksAlive = 0
		car.Controls = Controls{}
		car.Polygon = car.createPolygon()
		if car.CarType == CarTypeDummy {
			car.Controls.Forward = true
		}
		// Fresh random brain for now
		if car.Brain != nil {
			car.Brain = nn.NewNetwork(e.Config.NeuronCounts)
		}
	}
	// Reset traffic too
	e.Traffic = SpawnTraffic(e.Road, DefaultTraffic())
	e.Tick = 0
}

func (e *Engine) fanout() {
	e.mu.RLock()
	state := EngineState{
		Cars:       e.Cars,
		Traffic:    e.Traffic,
		BestCar:    e.BestCar,
		Tick:       e.Tick,
		Generation: e.Generation,
	}
	subs := make([]*Subscriber, 0, len(e.subscribers))
	for s := range e.subscribers {
		subs = append(subs, s)
	}
	e.mu.RUnlock()

	for _, s := range subs {
		select {
		case s.Ch <- state:
		default:
			// drop frame if subscriber is slow
		}
	}
}

// RoadState returns serializable road data for the init message.
func (e *Engine) RoadState() RoadState {
	return RoadState{
		X:         e.Road.X,
		Width:     e.Road.Width,
		LaneCount: e.Road.LaneCount,
		Left:      e.Road.Left,
		Right:     e.Road.Right,
		Top:       e.Road.Top,
		Bottom:    e.Road.Bottom,
		Borders: [2][2]types.Point{
			{e.Road.Borders[0][0], e.Road.Borders[0][1]},
			{e.Road.Borders[1][0], e.Road.Borders[1][1]},
		},
	}
}

// RoadState is the JSON-serializable road info.
type RoadState struct {
	X         float64              `json:"x"`
	Width     float64              `json:"width"`
	LaneCount int                  `json:"laneCount"`
	Left      float64              `json:"left"`
	Right     float64              `json:"right"`
	Top       float64              `json:"top"`
	Bottom    float64              `json:"bottom"`
	Borders   [2][2]types.Point    `json:"borders"`
}

func cloneBrain(src *nn.NeuralNetwork) *nn.NeuralNetwork {
	clone := &nn.NeuralNetwork{
		Levels: make([]*nn.Level, len(src.Levels)),
	}
	for i, l := range src.Levels {
		nl := &nn.Level{
			Inputs:  make([]float64, len(l.Inputs)),
			Outputs: make([]float64, len(l.Outputs)),
			Biases:  make([]float64, len(l.Biases)),
			Weights: make([][]float64, len(l.Weights)),
		}
		copy(nl.Biases, l.Biases)
		for j := range l.Weights {
			nl.Weights[j] = make([]float64, len(l.Weights[j]))
			copy(nl.Weights[j], l.Weights[j])
		}
		clone.Levels[i] = nl
	}
	return clone
}
