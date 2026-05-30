package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/lukas-arnold/garden-equipment-log/internal/configs"
	"github.com/lukas-arnold/garden-equipment-log/internal/models"
	"github.com/lukas-arnold/garden-equipment-log/internal/storage"
	"github.com/lukas-arnold/garden-equipment-log/internal/utils"
)

func HandleAddDeviceOperationPost(w http.ResponseWriter, r *http.Request) {
	deviceId, err := utils.ConvertId(r.PathValue("deviceId"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
	err = storage.AddDeviceOperation(deviceId, models.DeviceOperationInput{
		StartTime: r.FormValue("startTime"),
		EndTime:   r.FormValue("endTime"),
		Note:      r.FormValue("note"),
	})
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
	http.Redirect(w, r, fmt.Sprintf("/device/history/%d", deviceId), http.StatusFound)
}

func HandleDeleteDeviceOperation(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
	deviceId, lookupErr := storage.GetDeviceIdByOperationId(id)
	err = storage.DeleteDeviceOperation(id)
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
	}
	redirectTarget := "/devices"
	if lookupErr == nil {
		redirectTarget = fmt.Sprintf("/device/history/%d", deviceId)
	}
	http.Redirect(w, r, redirectTarget, http.StatusFound)
}

func HandleEditDeviceOperationGet(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	operation, err := storage.GetDeviceOperation(id)
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
		return
	}
	deviceId, err := storage.GetDeviceIdByOperationId(id)
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
		return
	}
	tmpl := template.Must(
		template.New("editOperation.html").Funcs(getTemplateFuncs()).ParseFS(configs.GetWebFiles(), "templates/device/editOperation.html"),
	)
	err = tmpl.Execute(w, struct {
		DeviceId  int64
		Operation models.DeviceOperation
	}{DeviceId: deviceId, Operation: operation})
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
}

func HandleSaveDeviceOperationPost(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	deviceId, err := storage.GetDeviceIdByOperationId(id)
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
		return
	}
	err = storage.UpdateDeviceOperation(id, models.DeviceOperationInput{
		StartTime: r.FormValue("startTime"),
		EndTime:   r.FormValue("endTime"),
		Note:      r.FormValue("note"),
	})
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/device/history/%d", deviceId), http.StatusFound)
}
