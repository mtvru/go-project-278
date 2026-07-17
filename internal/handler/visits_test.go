package handler_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/mtvru/go-project-278/internal/schemas"
)

func TestRedirect(t *testing.T) {
	router, _ := setupRouter(t)

	createLink(t, router, `{"original_url":"https://example.com/target","short_name":"goto"}`)

	rec := doJSON(t, router, http.MethodGet, "/r/goto", "")
	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "https://example.com/target" {
		t.Errorf("expected redirect to target, got %q", loc)
	}
}

func TestRedirectRecordsVisit(t *testing.T) {
	router, _ := setupRouter(t)

	createLink(t, router, `{"original_url":"https://example.com/target","short_name":"track"}`)
	doJSON(t, router, http.MethodGet, "/r/track", "")

	rec := doJSON(t, router, http.MethodGet, "/api/link_visits", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var visits []schemas.VisitResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &visits); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(visits) != 1 {
		t.Fatalf("expected 1 visit, got %d", len(visits))
	}
	if visits[0].Status != http.StatusFound {
		t.Errorf("expected status 302, got %d", visits[0].Status)
	}
	if visits[0].LinkID == 0 {
		t.Error("expected non-zero link_id")
	}
}

func TestRedirectNotFound(t *testing.T) {
	router, _ := setupRouter(t)

	rec := doJSON(t, router, http.MethodGet, "/r/missing", "")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestListVisitsPagination(t *testing.T) {
	router, _ := setupRouter(t)

	createLink(t, router, `{"original_url":"https://example.com/target","short_name":"multi"}`)
	for i := 0; i < 3; i++ {
		doJSON(t, router, http.MethodGet, "/r/multi", "")
	}

	rec := doJSON(t, router, http.MethodGet, "/api/link_visits?range=[0,2]", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var visits []schemas.VisitResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &visits); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(visits) != 2 {
		t.Errorf("expected 2 visits, got %d", len(visits))
	}
	if got := rec.Header().Get("Content-Range"); got != "link_visits 0-2/3" {
		t.Errorf("expected Content-Range 'link_visits 0-2/3', got %q", got)
	}
}
