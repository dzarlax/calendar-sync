package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dzarlax/calendar-sync/internal/sync"
)

func main() {
	cfg := sync.Config{
		APIBaseURL:   getEnv("REST_API_URL", "http://calendar-mcp:8080"),
		APIKey:       getEnv("API_KEY", ""),
		SyncSource:   getEnv("SYNC_SOURCE", ""),
		SyncTarget:   getEnv("SYNC_TARGET", ""),
		StateFile:    getEnv("STATE_FILE", "/data/sync_state.json"),
		SyncInterval: 10 * time.Minute,
	}

	if cfg.SyncSource == "" || cfg.SyncTarget == "" {
		log.Fatal("SYNC_SOURCE and SYNC_TARGET are required")
	}
	if cfg.APIKey == "" {
		log.Fatal("API_KEY is required")
	}

	client := sync.NewRestClient(cfg.APIBaseURL, cfg.APIKey)
	state, err := sync.NewStateManager(cfg.StateFile)
	if err != nil {
		log.Fatalf("failed to init state manager: %v", err)
	}

	syncer := sync.NewSyncer(client, state, cfg.SyncSource, cfg.SyncTarget)
	scheduler := sync.NewScheduler(syncer, cfg.SyncInterval)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go scheduler.Start(ctx)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	sig := <-sigCh
	log.Printf("received %s, shutting down", sig)
	cancel()
	scheduler.Stop()
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
