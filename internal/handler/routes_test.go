package handler

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/lukas-arnold/garden-equipment-log/internal/storage"
)

func TestRoutes(t *testing.T) {

	store := storage.New(
		filepath.Join(
			t.TempDir(),
			"test.json",
		),
	)

	h := New(
		store,
	)

	mux := http.NewServeMux()

	RegisterRoutes(
		mux,
		h,
	)

	tests := []struct {
		method string
		path   string
	}{
		{
			method: "GET",
			path:   "/",
		},
		{
			method: "GET",
			path:   "/devices",
		},
		{
			method: "GET",
			path:   "/device/add",
		},
		{
			method: "POST",
			path:   "/device/add",
		},
		{
			method: "GET",
			path:   "/device/edit/1",
		},
		{
			method: "POST",
			path:   "/device/save/1",
		},
		{
			method: "GET",
			path:   "/device/delete/1",
		},
		{
			method: "GET",
			path:   "/device/history/1",
		},

		{
			method: "GET",
			path:   "/device-operation/add/1",
		},
		{
			method: "POST",
			path:   "/device-operation/add/1",
		},
		{
			method: "GET",
			path:   "/device-operation/edit/1",
		},
		{
			method: "POST",
			path:   "/device-operation/save/1",
		},
		{
			method: "GET",
			path:   "/device-operation/delete/1",
		},

		{
			method: "GET",
			path:   "/bottles",
		},
		{
			method: "GET",
			path:   "/bottle/add",
		},
		{
			method: "POST",
			path:   "/bottle/add",
		},
		{
			method: "GET",
			path:   "/bottle/edit/1",
		},
		{
			method: "POST",
			path:   "/bottle/save/1",
		},
		{
			method: "GET",
			path:   "/bottle/delete/1",
		},
		{
			method: "GET",
			path:   "/bottle/history/1",
		},

		{
			method: "GET",
			path:   "/bottle-operation/add/1",
		},
		{
			method: "POST",
			path:   "/bottle-operation/add/1",
		},
		{
			method: "GET",
			path:   "/bottle-operation/edit/1",
		},
		{
			method: "POST",
			path:   "/bottle-operation/save/1",
		},
		{
			method: "GET",
			path:   "/bottle-operation/delete/1",
		},

		{
			method: "GET",
			path:   "/service-worker",
		},
		{
			method: "GET",
			path:   "/web/",
		},
	}

	for _, test := range tests {

		req := httptest.NewRequest(
			test.method,
			test.path,
			nil,
		)

		rec := httptest.NewRecorder()

		mux.ServeHTTP(
			rec,
			req,
		)

		if rec.Code == http.StatusNotFound {
			t.Fatalf(
				"route missing: %s %s",
				test.method,
				test.path,
			)
		}
	}
}
