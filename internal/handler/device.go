package handler

import (
	"fmt"
	"net/http"

	"github.com/lukas-arnold/garden-equipment-log/internal/configs"
	"github.com/lukas-arnold/garden-equipment-log/internal/language"
	"github.com/lukas-arnold/garden-equipment-log/internal/models"
	"github.com/lukas-arnold/garden-equipment-log/internal/utils"
)

type deviceFormData struct {
	Title       string
	Action      string
	SubmitLabel string
	Device      models.Device
}

type deviceHistoryPageData struct {
	Device        models.Device
	Chart         models.ChartModel
	OperationRows []deviceOperationRow
}

func (h *Handler) HandleDevicesView(w http.ResponseWriter, r *http.Request) {
	devices, err := h.storage.GetDevices()
	if err != nil {
		handleError(w, err, http.StatusNotFound)
		return
	}

	h.renderTemplate(w, "templates/device/index.html", devices)
}

func (h *Handler) HandleAddDeviceGet(w http.ResponseWriter, r *http.Request) {
	lang := configs.GetLanguage()

	h.renderTemplate(w, "templates/device/add.html", deviceFormData{
		Title:       language.T(lang, "addDevice"),
		Action:      "/device/add",
		SubmitLabel: language.T(lang, "save"),
		Device:      models.Device{},
	})
}

func (h *Handler) HandleAddDevicePost(w http.ResponseWriter, r *http.Request) {
	purchasePrice, err := utils.ConvertFloat(r.FormValue("purchasePrice"))
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	err = h.storage.AddDevice(models.DeviceInput{
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

func (h *Handler) HandleEditDevice(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	device, err := h.storage.GetDevice(id)
	if err != nil {
		handleError(w, err, http.StatusNotFound)
		return
	}

	lang := configs.GetLanguage()
	h.renderTemplate(w, "templates/device/add.html", deviceFormData{
		Title:       language.T(lang, "editDevice"),
		Action:      fmt.Sprintf("/device/save/%d", device.Id),
		SubmitLabel: language.T(lang, "save"),
		Device:      device,
	})
}

func (h *Handler) HandleSaveDevice(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	device, err := h.storage.GetDevice(id)
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

	if err := h.storage.UpdateDevice(device); err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/devices", http.StatusFound)
}

func (h *Handler) HandleDeleteDevice(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	if err := h.storage.DeleteDevice(id); err != nil {
		handleError(w, err, http.StatusNotFound)
		return
	}

	http.Redirect(w, r, "/devices", http.StatusFound)
}

func (h *Handler) HandleDeviceHistory(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	device, err := h.storage.GetDevice(id)
	if err != nil {
		handleError(w, err, http.StatusNotFound)
		return
	}

	chart, err := h.storage.GetDeviceChart(id)
	if err != nil {
		handleError(w, err, http.StatusNotFound)
		return
	}

	h.renderHistoryTemplate(w, "templates/device/history.html", deviceHistoryPageData{
		Device:        device,
		Chart:         chart,
		OperationRows: buildDeviceHistoryRows(device),
	})
}
