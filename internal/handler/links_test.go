package handler_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/mtvru/go-project-278/internal/schemas"
)

func decodeLink(t *testing.T, body []byte) schemas.LinkResponse {
	t.Helper()
	var resp schemas.LinkResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("decode link: %v (body: %s)", err, body)
	}
	return resp
}

func createLink(t *testing.T, router http.Handler, body string) schemas.LinkResponse {
	t.Helper()
	rec := doJSON(t, router, http.MethodPost, "/api/links", body)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create link: expected 201, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	return decodeLink(t, rec.Body.Bytes())
}

func TestCreateLink(t *testing.T) {
	router, _ := setupRouter(t)

	link := createLink(t, router, `{"original_url":"https://example.com/long","short_name":"exmpl"}`)
	if link.ID == 0 {
		t.Error("expected non-zero id")
	}
	if link.ShortName != "exmpl" {
		t.Errorf("expected short_name exmpl, got %q", link.ShortName)
	}
	if link.OriginalURL != "https://example.com/long" {
		t.Errorf("unexpected original_url %q", link.OriginalURL)
	}
	if want := testBaseURL + "/r/exmpl"; link.ShortURL != want {
		t.Errorf("expected short_url %q, got %q", want, link.ShortURL)
	}
}

func TestCreateLinkGeneratesShortName(t *testing.T) {
	router, _ := setupRouter(t)

	link := createLink(t, router, `{"original_url":"https://example.com/long"}`)
	if link.ShortName == "" {
		t.Fatal("expected generated short_name")
	}
	if want := testBaseURL + "/r/" + link.ShortName; link.ShortURL != want {
		t.Errorf("expected short_url %q, got %q", want, link.ShortURL)
	}
}

func TestCreateLinkValidationError(t *testing.T) {
	router, _ := setupRouter(t)

	rec := doJSON(t, router, http.MethodPost, "/api/links", `{"original_url":"not-a-url"}`)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	var resp struct {
		Errors map[string]string `json:"errors"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if _, ok := resp.Errors["original_url"]; !ok {
		t.Errorf("expected error for original_url, got %v", resp.Errors)
	}
}

func TestCreateLinkShortNameTooShort(t *testing.T) {
	router, _ := setupRouter(t)

	rec := doJSON(t, router, http.MethodPost, "/api/links", `{"original_url":"https://example.com","short_name":"ab"}`)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d (body: %s)", rec.Code, rec.Body.String())
	}
}

func TestCreateLinkInvalidJSON(t *testing.T) {
	router, _ := setupRouter(t)

	rec := doJSON(t, router, http.MethodPost, "/api/links", `{"original_url":`)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	var resp struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Error != "invalid request" {
		t.Errorf("expected 'invalid request', got %q", resp.Error)
	}
}

func TestCreateLinkDuplicateShortName(t *testing.T) {
	router, _ := setupRouter(t)

	createLink(t, router, `{"original_url":"https://example.com/a","short_name":"dup"}`)
	rec := doJSON(t, router, http.MethodPost, "/api/links", `{"original_url":"https://example.com/b","short_name":"dup"}`)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	var resp struct {
		Errors map[string]string `json:"errors"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Errors["short_name"] != "short name already in use" {
		t.Errorf("unexpected error: %v", resp.Errors)
	}
}

func TestGetLink(t *testing.T) {
	router, _ := setupRouter(t)

	created := createLink(t, router, `{"original_url":"https://example.com","short_name":"getme"}`)
	rec := doJSON(t, router, http.MethodGet, "/api/links/"+itoa(created.ID), "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	got := decodeLink(t, rec.Body.Bytes())
	if got.ID != created.ID {
		t.Errorf("expected id %d, got %d", created.ID, got.ID)
	}
}

func TestGetLinkNotFound(t *testing.T) {
	router, _ := setupRouter(t)

	rec := doJSON(t, router, http.MethodGet, "/api/links/999999", "")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestUpdateLink(t *testing.T) {
	router, _ := setupRouter(t)

	created := createLink(t, router, `{"original_url":"https://example.com","short_name":"before"}`)
	rec := doJSON(t, router, http.MethodPut, "/api/links/"+itoa(created.ID),
		`{"original_url":"https://updated.com","short_name":"after"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	got := decodeLink(t, rec.Body.Bytes())
	if got.OriginalURL != "https://updated.com" || got.ShortName != "after" {
		t.Errorf("update not applied: %+v", got)
	}
}

func TestUpdateLinkNotFound(t *testing.T) {
	router, _ := setupRouter(t)

	rec := doJSON(t, router, http.MethodPut, "/api/links/999999",
		`{"original_url":"https://updated.com","short_name":"after"}`)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestDeleteLink(t *testing.T) {
	router, _ := setupRouter(t)

	created := createLink(t, router, `{"original_url":"https://example.com","short_name":"delme"}`)
	rec := doJSON(t, router, http.MethodDelete, "/api/links/"+itoa(created.ID), "")
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}

	rec = doJSON(t, router, http.MethodGet, "/api/links/"+itoa(created.ID), "")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 after delete, got %d", rec.Code)
	}
}

func TestDeleteLinkNotFound(t *testing.T) {
	router, _ := setupRouter(t)

	rec := doJSON(t, router, http.MethodDelete, "/api/links/999999", "")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestListLinksPagination(t *testing.T) {
	router, _ := setupRouter(t)

	for i := 0; i < 3; i++ {
		createLink(t, router, `{"original_url":"https://example.com/`+itoa(int64(i))+`"}`)
	}

	rec := doJSON(t, router, http.MethodGet, "/api/links?range=[0,2]", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var links []schemas.LinkResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &links); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(links) != 2 {
		t.Errorf("expected 2 links, got %d", len(links))
	}
	if got := rec.Header().Get("Content-Range"); got != "links 0-2/3" {
		t.Errorf("expected Content-Range 'links 0-2/3', got %q", got)
	}
}

func TestListLinksPaginationOffset(t *testing.T) {
	router, _ := setupRouter(t)

	for i := 0; i < 3; i++ {
		createLink(t, router, `{"original_url":"https://example.com/`+itoa(int64(i))+`"}`)
	}

	rec := doJSON(t, router, http.MethodGet, "/api/links?range=[1,3]", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var links []schemas.LinkResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &links); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(links) != 2 {
		t.Errorf("expected 2 links, got %d", len(links))
	}
	if got := rec.Header().Get("Content-Range"); got != "links 1-3/3" {
		t.Errorf("expected Content-Range 'links 1-3/3', got %q", got)
	}
}

func TestListLinksEmptyReturnsArray(t *testing.T) {
	router, _ := setupRouter(t)

	rec := doJSON(t, router, http.MethodGet, "/api/links", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if body := rec.Body.String(); body != "[]" {
		t.Errorf("expected '[]', got %q", body)
	}
}
