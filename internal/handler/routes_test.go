package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoutes(t *testing.T) {
	h := testHandler(t)
	mux := http.NewServeMux()
	RegisterRoutes(mux, h)

	tests := []struct {
		method string
		path   string
	}{
		{"GET", "/"},

		{"GET", "/service-worker.js"},
		{"GET", "/web/"},

		{"GET", "/devices"},
		{"GET", "/device/add"},
		{"POST", "/device/add"},
		{"GET", "/device/edit/1"},
		{"POST", "/device/save/1"},
		{"GET", "/device/delete/1"},
		{"GET", "/device/history/1"},

		{"GET", "/device-operation/add/1"},
		{"POST", "/device-operation/add/1"},
		{"GET", "/device-operation/edit/1"},
		{"POST", "/device-operation/save/1"},
		{"GET", "/device-operation/delete/1"},

		{"GET", "/bottles"},
		{"GET", "/bottle/add"},
		{"POST", "/bottle/add"},
		{"GET", "/bottle/edit/1"},
		{"POST", "/bottle/save/1"},
		{"GET", "/bottle/delete/1"},
		{"GET", "/bottle/history/1"},

		{"GET", "/bottle-operation/add/1"},
		{"POST", "/bottle-operation/add/1"},
		{"GET", "/bottle-operation/edit/1"},
		{"POST", "/bottle-operation/save/1"},
		{"GET", "/bottle-operation/delete/1"},
	}

	for _, test := range tests {
		req := httptest.NewRequest(test.method, test.path, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code == http.StatusNotFound {
			t.Errorf("route missing or returned 404: %s %s", test.method, test.path)
		}
	}
}
