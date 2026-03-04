package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/OJIEG65/car-go-ai/be/internal/nn"
	"github.com/OJIEG65/car-go-ai/be/internal/sim"
	"github.com/OJIEG65/car-go-ai/be/internal/types"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

// HandleWS upgrades to WebSocket and streams simulation state.
func HandleWS(engine *sim.Engine, store *nn.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("ws upgrade error: %v", err)
			return
		}
		defer conn.Close()

		sub := engine.Subscribe()
		defer engine.Unsubscribe(sub)

		// Send init message with road state and config
		sendInit(conn, engine)

		// Read pump (handles client messages)
		go readPump(conn, engine, store)

		// Write pump (streams state)
		writePump(conn, sub, engine)
	}
}

func sendInit(conn *websocket.Conn, engine *sim.Engine) {
	roadData, _ := json.Marshal(engine.RoadState())
	cfgData, _ := json.Marshal(engine.Config)

	initPayload := types.InitPayload{
		Road:   roadData,
		Config: cfgData,
	}
	payload, _ := json.Marshal(initPayload)

	msg := types.WSMessage{
		Type:    "init",
		Payload: payload,
	}

	conn.SetWriteDeadline(time.Now().Add(writeWait))
	conn.WriteJSON(msg)
}

func writePump(conn *websocket.Conn, sub *sim.Subscriber, engine *sim.Engine) {
	pingTicker := time.NewTicker(pingPeriod)
	defer pingTicker.Stop()

	for {
		select {
		case state, ok := <-sub.Ch:
			if !ok {
				return
			}
			worldState := buildWorldState(state, engine)
			payload, _ := json.Marshal(worldState)

			msg := types.WSMessage{
				Type:    "state",
				Payload: payload,
			}

			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteJSON(msg); err != nil {
				log.Printf("ws write error: %v", err)
				return
			}

		case <-pingTicker.C:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func readPump(conn *websocket.Conn, engine *sim.Engine, store *nn.Store) {
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		var msg types.WSMessage
		if err := conn.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("ws read error: %v", err)
			}
			return
		}

		switch msg.Type {
		case "pause":
			engine.Pause()
		case "resume":
			engine.Resume()
		case "reset":
			engine.Reset()
		case "save_brain":
			state := engine.GetState()
			if state.BestCar != nil && state.BestCar.Brain != nil && store != nil {
				name := "best"
				if msg.Payload != nil {
					var p struct{ Name string `json:"name"` }
					json.Unmarshal(msg.Payload, &p)
					if p.Name != "" {
						name = p.Name
					}
				}
				if err := store.Save(name, state.BestCar.Brain); err != nil {
					log.Printf("save brain error: %v", err)
				} else {
					log.Printf("brain saved as %q", name)
				}
			}
		case "load_brain":
			if store != nil {
				name := "best"
				if msg.Payload != nil {
					var p struct{ Name string `json:"name"` }
					json.Unmarshal(msg.Payload, &p)
					if p.Name != "" {
						name = p.Name
					}
				}
				brain, err := store.Load(name)
				if err != nil {
					log.Printf("load brain error: %v", err)
				} else {
					engine.SetBrainForAll(brain)
					log.Printf("brain %q loaded into all cars", name)
				}
			}
		default:
			log.Printf("unknown ws message type: %s", msg.Type)
		}
	}
}

func buildWorldState(state sim.EngineState, engine *sim.Engine) types.WorldState {
	cars := make([]types.CarState, len(state.Cars))
	bestIndex := 0
	for i, c := range state.Cars {
		cars[i] = types.CarState{
			X:       c.X,
			Y:       c.Y,
			Angle:   c.Angle,
			Damaged: c.Damaged,
			Polygon: c.Polygon,
			Speed:   c.Speed,
		}
		if state.BestCar != nil && c == state.BestCar {
			bestIndex = i
		}
	}

	traffic := make([]types.CarState, len(state.Traffic))
	for i, c := range state.Traffic {
		traffic[i] = types.CarState{
			X:       c.X,
			Y:       c.Y,
			Angle:   c.Angle,
			Damaged: c.Damaged,
			Polygon: c.Polygon,
			Speed:   c.Speed,
		}
	}

	ws := types.WorldState{
		Cars:       cars,
		Traffic:    traffic,
		BestIndex:  bestIndex,
		Tick:       state.Tick,
		Generation: state.Generation,
	}

	// Include best car's sensor data for rendering
	if state.BestCar != nil && state.BestCar.Sensor != nil {
		sensor := state.BestCar.Sensor
		sensorState := &types.SensorState{
			Rays:     make([][2]types.Point, len(sensor.Rays)),
			Readings: sensor.Readings,
		}
		for i, ray := range sensor.Rays {
			sensorState.Rays[i] = [2]types.Point{ray[0], ray[1]}
		}
		ws.BestSensor = sensorState
	}

	// Include best car's brain for network visualization
	if state.BestCar != nil && state.BestCar.Brain != nil {
		// Brain data will be sent periodically, not every tick (to save bandwidth)
		// For now we skip it — the frontend will request it via "get_brain"
	}

	return ws
}
