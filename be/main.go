package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/OJIEG65/car-go-ai/be/internal/api"
	"github.com/OJIEG65/car-go-ai/be/internal/config"
	"github.com/OJIEG65/car-go-ai/be/internal/nn"
	"github.com/OJIEG65/car-go-ai/be/internal/sim"
)

func main() {
	cfg := config.DefaultConfig()

	engine := sim.NewEngine(cfg)
	engine.Start()
	defer engine.Stop()

	// Brain persistence store
	store, err := nn.NewStore("data/brains")
	if err != nil {
		log.Fatalf("failed to create brain store: %v", err)
	}

	// Resolve frontend directory
	feDir := "../fe"
	if d := os.Getenv("FE_DIR"); d != "" {
		feDir = d
	}

	server := api.NewServer(":8080", engine, store, feDir)

	// Graceful shutdown on SIGINT/SIGTERM
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		log.Println("Shutting down...")
		engine.Stop()
		server.Shutdown()
		os.Exit(0)
	}()

	log.Println("CarGoAi backend running on :8080")
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
