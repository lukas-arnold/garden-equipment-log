package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lukas-arnold/garden-equipment-log/internal/handler"
	"github.com/lukas-arnold/garden-equipment-log/internal/storage"
)

func TestCreateServer(t *testing.T) {

	store :=
		storage.New(
			t.TempDir() + "/test.json",
		)

	h :=
		handler.New(store)

	server :=
		createServer(h)

	req :=
		httptest.NewRequest(
			"GET",
			"/",
			nil,
		)

	rec :=
		httptest.NewRecorder()

	server.ServeHTTP(
		rec,
		req,
	)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"got %d",
			rec.Code,
		)
	}
}
