package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mtvru/go-project-278/internal/config"
	"github.com/mtvru/go-project-278/internal/handler"
	"github.com/mtvru/go-project-278/internal/repository"
	"github.com/mtvru/go-project-278/internal/router"
	"github.com/mtvru/go-project-278/internal/service"
)

type App struct {
	router *gin.Engine
	pool   *pgxpool.Pool
}

func NewApp(ctx context.Context, cfg *config.Config) (*App, error) {
	if cfg.SentryDSN != "" {
		if err := sentry.Init(sentry.ClientOptions{Dsn: cfg.SentryDSN}); err != nil {
			log.Printf("sentry init: %v", err)
		}
	}

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	repo := repository.New(pool)
	linkService := service.NewLinkService(repo)
	visitService := service.NewVisitService(repo)

	linkHandler := handler.NewLinkHandler(linkService, cfg.BaseURL)
	visitHandler := handler.NewVisitHandler(visitService)

	r := router.New(linkHandler, visitHandler, cfg.AllowOrigins)

	return &App{
		router: r,
		pool:   pool,
	}, nil
}

func (a *App) Run(addr string) error {
	defer func() {
		a.pool.Close()
		sentry.Flush(2 * time.Second)
	}()
	return a.router.Run(addr)
}
