package handler

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/lukas-arnold/garden-equipment-log/internal/configs"
	"github.com/lukas-arnold/garden-equipment-log/internal/models"
	"github.com/lukas-arnold/garden-equipment-log/internal/storage"
	"github.com/lukas-arnold/garden-equipment-log/internal/utils"
)

func HandleAddDeviceOperationGet(w http.ResponseWriter, r *http.Request) {
	deviceId, err := utils.ConvertId(r.PathValue("deviceId"))
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	tmpl := template.Must(
		template.New("addOperation.html").Funcs(getTemplateFuncs()).ParseFS(configs.GetWebFiles(), "templates/device/addOperation.html"),
	)
	err = tmpl.Execute(w, struct{ DeviceId int64 }{DeviceId: deviceId})
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
}

func HandleAddDeviceOperationPost(w http.ResponseWriter, r *http.Request) {
	deviceId, err := utils.ConvertId(r.PathValue("deviceId"))
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	err = storage.AddDeviceOperation(deviceId, models.DeviceOperationInput{
		StartTime: r.FormValue("startTime"),
		EndTime:   r.FormValue("endTime"),
		Note:      r.FormValue("note"),
	})
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/device/history/%d", deviceId), http.StatusFound)
}

func HandleEditDeviceOperation(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	operation, err := storage.GetDeviceOperation(id)
	if err != nil {
		handleError(w, err, http.StatusNotFound)
		return
	}
	deviceId, err := storage.GetDeviceIdByOperationId(id)
	if err != nil {
		handleError(w, err, http.StatusNotFound)
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
		handleError(w, err, http.StatusInternalServerError)
		return
	}
}

func HandleSaveDeviceOperation(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	deviceId, err := storage.GetDeviceIdByOperationId(id)
	if err != nil {
		handleError(w, err, http.StatusNotFound)
		return
	}
	err = storage.UpdateDeviceOperation(id, models.DeviceOperationInput{
		StartTime: r.FormValue("startTime"),
		EndTime:   r.FormValue("endTime"),
		Note:      r.FormValue("note"),
	})
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/device/history/%d", deviceId), http.StatusFound)
}

func HandleDeleteDeviceOperation(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	deviceId, lookupErr := storage.GetDeviceIdByOperationId(id)
	err = storage.DeleteDeviceOperation(id)
	if err != nil {
		handleError(w, err, http.StatusNotFound)
		return
	}
	redirectTarget := "/devices"
	if lookupErr == nil {
		redirectTarget = fmt.Sprintf("/device/history/%d", deviceId)
	}
	http.Redirect(w, r, redirectTarget, http.StatusFound)
}
