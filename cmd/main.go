package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ebookstore/internal/log"
	"github.com/ebookstore/internal/platform/config"
)

// Start HTTP Server and handle graceful shutdown
func main() {
	config.LoadConfiguration()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	server := NewServer()
	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			log.Default().Fatalf("failed starting server: %v", err)
		}
	}()
	log.Default().Info("starting HTTP server")

	<-done
	log.Default().Info("shutting down HTTP server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Default().Fatalf("failed to shutdown HTTP server")
	}

	log.Default().Info("server was shutdown properly")

}
