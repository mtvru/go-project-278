package app_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mtvru/go-project-278/internal/handler"
	"github.com/mtvru/go-project-278/internal/repository"
	"github.com/mtvru/go-project-278/internal/router"
	"github.com/mtvru/go-project-278/internal/service"
)

const testBaseURL = "http://localhost:8080"

func TestApp(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := repository.New(nil)
	linkService := service.NewLinkService(repo)
	visitService := service.NewVisitService(repo)
	linkHandler := handler.NewLinkHandler(linkService, testBaseURL)
	visitHandler := handler.NewVisitHandler(visitService)

	r := router.New(linkHandler, visitHandler, []string{"http://localhost:5173"})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if rec.Body.String() != "pong" {
		t.Fatalf("expected body %q, got %q", "pong", rec.Body.String())
	}
}
