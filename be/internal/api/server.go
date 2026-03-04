package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/OJIEG65/car-go-ai/be/internal/nn"
	"github.com/OJIEG65/car-go-ai/be/internal/sim"
)

// Server wraps the HTTP server with graceful shutdown.
type Server struct {
	httpServer *http.Server
	engine     *sim.Engine
}

// NewServer creates the API server wired to the simulation engine.
func NewServer(addr string, engine *sim.Engine, store *nn.Store, feDir string) *Server {
	mux := http.NewServeMux()

	// REST endpoints
	RegisterHandlers(mux, engine)

	// WebSocket endpoint
	mux.HandleFunc("/ws", HandleWS(engine, store))

	// Serve frontend static files
	if feDir != "" {
		fs := http.FileServer(http.Dir(feDir))
		mux.Handle("/", fs)
	}

	return &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: CORS(mux),
		},
		engine: engine,
	}
}

// Start begins listening. Blocks until shutdown.
func (s *Server) Start() error {
	log.Printf("API server listening on %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.httpServer.Shutdown(ctx)
}
