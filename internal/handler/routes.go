package handler

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("GET /", h.HandleView)

	mux.HandleFunc("GET /web/", h.HandleFiles)
	mux.HandleFunc("GET /service-worker.js", h.HandleServiceWorker)

	mux.HandleFunc("GET /devices", h.HandleDevicesView)
	mux.HandleFunc("GET /device/add", h.HandleAddDeviceGet)
	mux.HandleFunc("POST /device/add", h.HandleAddDevicePost)
	mux.HandleFunc("GET /device/edit/{id}", h.HandleEditDevice)
	mux.HandleFunc("POST /device/save/{id}", h.HandleSaveDevice)
	mux.HandleFunc("GET /device/delete/{id}", h.HandleDeleteDevice)
	mux.HandleFunc("GET /device/history/{id}", h.HandleDeviceHistory)

	mux.HandleFunc("GET /device-operation/add/{deviceId}", h.HandleAddDeviceOperationGet)
	mux.HandleFunc("POST /device-operation/add/{deviceId}", h.HandleAddDeviceOperationPost)
	mux.HandleFunc("GET /device-operation/edit/{id}", h.HandleEditDeviceOperation)
	mux.HandleFunc("POST /device-operation/save/{id}", h.HandleSaveDeviceOperation)
	mux.HandleFunc("GET /device-operation/delete/{id}", h.HandleDeleteDeviceOperation)

	mux.HandleFunc("GET /bottles", h.HandleBottlesView)
	mux.HandleFunc("GET /bottle/add", h.HandleAddBottleGet)
	mux.HandleFunc("POST /bottle/add", h.HandleAddBottlePost)
	mux.HandleFunc("GET /bottle/edit/{id}", h.HandleEditBottle)
	mux.HandleFunc("POST /bottle/save/{id}", h.HandleSaveBottle)
	mux.HandleFunc("GET /bottle/delete/{id}", h.HandleDeleteBottle)
	mux.HandleFunc("GET /bottle/history/{id}", h.HandleBottleHistory)

	mux.HandleFunc("GET /bottle-operation/add/{bottleId}", h.HandleAddBottleOperationGet)
	mux.HandleFunc("POST /bottle-operation/add/{bottleId}", h.HandleAddBottleOperationPost)
	mux.HandleFunc("GET /bottle-operation/edit/{id}", h.HandleEditBottleOperation)
	mux.HandleFunc("POST /bottle-operation/save/{id}", h.HandleSaveBottleOperation)
	mux.HandleFunc("GET /bottle-operation/delete/{id}", h.HandleDeleteBottleOperation)
}
