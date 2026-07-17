package main

import (
	"context"
	"log"

	"github.com/mtvru/go-project-278/internal/app"
	"github.com/mtvru/go-project-278/internal/config"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()
	application, err := app.NewApp(ctx, &cfg)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	if err := application.Run(":8080"); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
