package sim

// TrafficConfig defines a single traffic car spawn.
type TrafficConfig struct {
	Lane     int
	Y        float64
	MaxSpeed float64
}

// DefaultTraffic returns the traffic layout matching the original JS.
func DefaultTraffic() []TrafficConfig {
	return []TrafficConfig{
		{Lane: 1, Y: -100, MaxSpeed: 2},
	}
}

// SpawnTraffic creates traffic cars from config on the given road.
func SpawnTraffic(road *Road, configs []TrafficConfig) []*Car {
	cars := make([]*Car, len(configs))
	for i, tc := range configs {
		x := road.GetLaneCenter(tc.Lane)
		cars[i] = NewCar(x, tc.Y, 30, 50, CarTypeDummy, tc.MaxSpeed, 0.2, 0.05, 0.03, nil, 0, 0, 0)
	}
	return cars
}
