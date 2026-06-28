package handler

import (
	"net/http"

	"github.com/lukas-arnold/garden-equipment-log/internal/configs"
)

func (h *Handler) HandleFiles(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServerFS(configs.GetWebFiles())
	handler := http.StripPrefix("/web/", fs)

	handler.ServeHTTP(w, r)
}

func (h *Handler) HandleServiceWorker(w http.ResponseWriter, r *http.Request) {
	http.ServeFileFS(
		w,
		r,
		configs.GetWebFiles(),
		"service-worker.js",
	)
}
