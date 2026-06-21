package handler

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/lukas-arnold/garden-equipment-log/internal/configs"
	"github.com/lukas-arnold/garden-equipment-log/internal/language"
	"github.com/lukas-arnold/garden-equipment-log/internal/models"
	"github.com/lukas-arnold/garden-equipment-log/internal/storage"
	"github.com/lukas-arnold/garden-equipment-log/internal/utils"
)

type deviceFormData struct {
	Title       string
	Action      string
	SubmitLabel string
	Device      models.Device
}

type deviceOperationRow struct {
	models.DeviceOperation
	Time float64
}

type deviceHistoryPageData struct {
	Device        models.Device
	Chart         models.ChartModel
	OperationRows []deviceOperationRow
}

func HandleDevicesView(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("base.html").Funcs(getTemplateFuncs()).ParseFS(configs.GetWebFiles(), "templates/base.html", "templates/device/index.html"),
	)
	devices, err := storage.GetDevices()
	if err != nil {
		handleError(w, err, http.StatusNotFound)
		return
	}
	err = tmpl.Execute(w, devices)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
}

func HandleAddDeviceGet(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("base.html").
			Funcs(getTemplateFuncs()).
			ParseFS(configs.GetWebFiles(), "templates/base.html", "templates/device/add.html"),
	)
	err := tmpl.Execute(w, deviceFormData{
		Title:       language.T(configs.GetLanguage(), "addDevice"),
		Action:      "/device/add",
		SubmitLabel: language.T(configs.GetLanguage(), "save"),
		Device:      models.Device{},
	})
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
}

func HandleAddDevicePost(w http.ResponseWriter, r *http.Request) {
	purchasePrice, err := utils.ConvertFloat(r.FormValue("purchasePrice"))
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	err = storage.AddDevice(models.DeviceInput{
		Name:          r.FormValue("name"),
		PurchaseDate:  r.FormValue("purchaseDate"),
		PurchasePrice: purchasePrice,
	})
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/devices", http.StatusFound)
}

func HandleEditDevice(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("base.html").
			Funcs(getTemplateFuncs()).
			ParseFS(configs.GetWebFiles(), "templates/base.html", "templates/device/add.html"),
	)
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	device, err := storage.GetDevice(id)
	if err != nil {
		handleError(w, err, http.StatusNotFound)
		return
	}
	err = tmpl.Execute(w, deviceFormData{
		Title:       language.T(configs.GetLanguage(), "editDevice"),
		Action:      fmt.Sprintf("/device/save/%d", device.Id),
		SubmitLabel: language.T(configs.GetLanguage(), "save"),
		Device:      device,
	})
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
}

func HandleSaveDevice(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	device, err := storage.GetDevice(id)
	if err != nil {
		handleError(w, err, http.StatusNotFound)
		return
	}
	purchasePrice, err := utils.ConvertFloat(r.FormValue("purchasePrice"))
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	device.Name = r.FormValue("name")
	device.PurchaseDate = r.FormValue("purchaseDate")
	device.PurchasePrice = purchasePrice
	err = storage.UpdateDevice(device)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/devices", http.StatusFound)
}

func HandleDeleteDevice(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	err = storage.DeleteDevice(id)
	if err != nil {
		handleError(w, err, http.StatusNotFound)
		return
	}
	http.Redirect(w, r, "/devices", http.StatusFound)
}

func HandleDeviceHistory(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	device, err := storage.GetDevice(id)
	if err != nil {
		handleError(w, err, http.StatusNotFound)
		return
	}
	chart, err := storage.GetDeviceChart(id)
	if err != nil {
		handleError(w, err, http.StatusNotFound)
		return
	}
	rows := buildDeviceHistoryRows(device)
	view := deviceHistoryPageData{
		Device:        device,
		Chart:         chart,
		OperationRows: rows,
	}
	tmpl := template.Must(
		template.New("baseHistory.html").
			Funcs(getTemplateFuncs()).
			ParseFS(
				configs.GetWebFiles(),
				"templates/baseHistory.html",
				"templates/device/history.html",
			),
	)
	err = tmpl.Execute(w, view)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
}

func buildDeviceHistoryRows(device models.Device) []deviceOperationRow {
	rows := make([]deviceOperationRow, 0, len(device.OperationHistory))
	for _, op := range device.OperationHistory {
		var minutes float64
		if op.StartTime != "" && op.EndTime != "" {
			start, err1 := utils.ParseDateTime(op.StartTime)
			end, err2 := utils.ParseDateTime(op.EndTime)
			if err1 == nil && err2 == nil {
				minutes = end.Sub(start).Minutes()
				if minutes < 0 {
					minutes = 0
				}
			}
		}
		rows = append(rows, deviceOperationRow{
			DeviceOperation: op,
			Time:            minutes,
		})
	}
	return rows
}
