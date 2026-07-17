package handler_test

import (
	"bytes"
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mtvru/go-project-278/db/migrations"
	"github.com/mtvru/go-project-278/internal/handler"
	"github.com/mtvru/go-project-278/internal/repository"
	"github.com/mtvru/go-project-278/internal/router"
	"github.com/mtvru/go-project-278/internal/service"
	"github.com/pressly/goose/v3"
)

const testBaseURL = "http://localhost:8080"

func setupRouter(t *testing.T) (*gin.Engine, *repository.Repository) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = os.Getenv("DATABASE_URL")
	}
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL/DATABASE_URL not set, skipping database tests")
	}

	migrate(t, dsn)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		t.Fatalf("connect to database: %v", err)
	}
	t.Cleanup(pool.Close)

	if _, err := pool.Exec(context.Background(), "TRUNCATE links, link_visits RESTART IDENTITY CASCADE"); err != nil {
		t.Fatalf("truncate tables: %v", err)
	}

	repo := repository.New(pool)
	linkService := service.NewLinkService(repo)
	visitService := service.NewVisitService(repo)
	linkHandler := handler.NewLinkHandler(linkService, testBaseURL)
	visitHandler := handler.NewVisitHandler(visitService)

	r := router.New(linkHandler, visitHandler, []string{"http://localhost:5173"})
	return r, repo
}

func migrate(t *testing.T, dsn string) {
	t.Helper()
	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("open migration db: %v", err)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			t.Logf("close migration db: %v", err)
		}
	}()

	goose.SetBaseFS(migrations.FS)
	if err := goose.SetDialect("postgres"); err != nil {
		t.Fatalf("set goose dialect: %v", err)
	}
	if err := goose.Up(sqlDB, "."); err != nil {
		t.Fatalf("run migrations: %v", err)
	}
}

func doJSON(t *testing.T, router http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func itoa(n int64) string {
	return strconv.FormatInt(n, 10)
}
