package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/OJIEG65/car-go-ai/be/internal/sim"
)

// RegisterHandlers sets up REST API routes.
func RegisterHandlers(mux *http.ServeMux, engine *sim.Engine) {
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/api/config", handleConfig(engine))
	mux.HandleFunc("/api/state", handleState(engine))
	mux.HandleFunc("/api/brain/best", handleBestBrain(engine))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"ok"}`)
}

func handleConfig(engine *sim.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(engine.Config)
	}
}

func handleState(engine *sim.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		state := engine.GetState()
		json.NewEncoder(w).Encode(map[string]any{
			"tick":       state.Tick,
			"generation": state.Generation,
			"carCount":   len(state.Cars),
			"paused":     engine.IsPaused(),
		})
	}
}

func handleBestBrain(engine *sim.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		state := engine.GetState()
		if state.BestCar != nil && state.BestCar.Brain != nil {
			json.NewEncoder(w).Encode(state.BestCar.Brain)
		} else {
			http.Error(w, `{"error":"no best car"}`, http.StatusNotFound)
		}
	}
}
