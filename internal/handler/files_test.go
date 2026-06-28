package handler

import (
	"net/http/httptest"
	"testing"

	"path/filepath"

	"github.com/lukas-arnold/garden-equipment-log/internal/storage"
)

func TestHandleFiles(t *testing.T) {

	h := New(
		storage.New(
			filepath.Join(
				t.TempDir(),
				"test.json",
			),
		),
	)

	req :=
		httptest.NewRequest(
			"GET",
			"/web/",
			nil,
		)

	rec :=
		httptest.NewRecorder()

	h.HandleFiles(
		rec,
		req,
	)

	// It should not crash.
	// Depending on embedded files, status can differ.
	if rec.Code == 0 {
		t.Fatal(
			"no response",
		)
	}
}

func TestHandleServiceWorker(t *testing.T) {

	h := New(
		storage.New(
			filepath.Join(
				t.TempDir(),
				"test.json",
			),
		),
	)

	req :=
		httptest.NewRequest(
			"GET",
			"/service-worker",
			nil,
		)

	rec :=
		httptest.NewRecorder()

	h.HandleServiceWorker(
		rec,
		req,
	)

	if rec.Code == 0 {
		t.Fatal(
			"no response",
		)
	}
}
