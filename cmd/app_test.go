package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lukas-arnold/garden-equipment-log/internal/handler"
	"github.com/lukas-arnold/garden-equipment-log/internal/storage"
)

func TestCreateServer(t *testing.T) {
	store := storage.New(t.TempDir() + "/test.json")
	h := handler.New(store)
	server := createServer(h)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status OK, got %d", rec.Code)
	}
}
