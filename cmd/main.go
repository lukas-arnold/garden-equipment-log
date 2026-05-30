package main

import (
	"log"
	"net/http"

	"github.com/lukas-arnold/garden-equipment-log/internal/configs"
	"github.com/lukas-arnold/garden-equipment-log/internal/handler"
	"github.com/lukas-arnold/garden-equipment-log/internal/language"
)

func main() {
	err := language.LoadLanguages()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /service-worker", handler.HandleServiceWorker)
	mux.HandleFunc("GET /web/", handler.HandleFiles)

	mux.HandleFunc("GET /", handler.HandleView)

	mux.HandleFunc("GET /devices", handler.HandleDevicesView)
	mux.HandleFunc("GET /device/add", handler.HandleAddDeviceGet)
	mux.HandleFunc("POST /device/add", handler.HandleAddDevicePost)
	mux.HandleFunc("GET /device/edit/{id}", handler.HandleEditDevice)
	mux.HandleFunc("POST /device/save/{id}", handler.HandleSaveDevice)
	mux.HandleFunc("GET /device/delete/{id}", handler.HandleDeleteDevice)
	mux.HandleFunc("GET /device/history/{id}", handler.HandleDeviceHistory)

	mux.HandleFunc("GET /device-operation/add/{deviceId}", handler.HandleAddDeviceGetOperation)
	mux.HandleFunc("POST /device-operation/add/{deviceId}", handler.HandleAddDevicePostOperation)
	mux.HandleFunc("GET /device-operation/edit/{id}", handler.HandleEditDeviceOperation)
	mux.HandleFunc("POST /device-operation/save/{id}", handler.HandleSaveDeviceOperation)
	mux.HandleFunc("GET /device-operation/delete/{id}", handler.HandleDeleteDeviceOperation)

	mux.HandleFunc("GET /bottles", handler.HandleBottlesView)
	mux.HandleFunc("GET /bottle/add", handler.HandleAddBottleGet)
	mux.HandleFunc("POST /bottle/add", handler.HandleAddBottlePost)
	mux.HandleFunc("GET /bottle/edit/{id}", handler.HandleEditBottle)
	mux.HandleFunc("POST /bottle/save/{id}", handler.HandleSaveBottle)
	mux.HandleFunc("GET /bottle/delete/{id}", handler.HandleDeleteBottle)
	mux.HandleFunc("GET /bottle/history/{id}", handler.HandleBottleHistory)

	mux.HandleFunc("GET /bottle-operation/add/{bottleId}", handler.HandleAddBottleOperationGet)
	mux.HandleFunc("POST /bottle-operation/add/{bottleId}", handler.HandleAddBottleOperationPost)
	mux.HandleFunc("GET /bottle-operation/edit/{id}", handler.HandleEditBottleOperation)
	mux.HandleFunc("POST /bottle-operation/save/{id}", handler.HandleSaveBottleOperation)
	mux.HandleFunc("GET /bottle-operation/delete/{id}", handler.HandleDeleteBottleOperation)

	log.Printf("Garden Equipment Log running on %s", configs.GetPort())
	log.Fatal(http.ListenAndServe(configs.GetPort(), mux))
}
