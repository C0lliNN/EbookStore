package main

import (
	"context"

	"github.com/ebookstore/internal/container"
	"github.com/ebookstore/internal/platform/config"
)

// Start HTTP Server and handle graceful shutdown
func main() {
	config.LoadConfiguration()

	ctx := context.Background()
	c := container.New()

	c.Start(ctx)
}
