package handler

import (
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/lukas-arnold/garden-equipment-log/internal/storage"
)

func TestHandleFiles(t *testing.T) {
	h := New(storage.New(filepath.Join(t.TempDir(), "test.json")))

	req := httptest.NewRequest("GET", "/web/", nil)
	rec := httptest.NewRecorder()

	h.HandleFiles(rec, req)

	if rec.Code == 0 {
		t.Fatal("no response code set")
	}
}

func TestHandleServiceWorker(t *testing.T) {
	h := New(storage.New(filepath.Join(t.TempDir(), "test.json")))

	req := httptest.NewRequest("GET", "/service-worker", nil)
	rec := httptest.NewRecorder()

	h.HandleServiceWorker(rec, req)

	if rec.Code == 0 {
		t.Fatal("no response code set")
	}
}
