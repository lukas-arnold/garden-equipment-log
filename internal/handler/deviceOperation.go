package handler

import (
	"fmt"
	"net/http"

	"github.com/lukas-arnold/garden-equipment-log/internal/models"
	"github.com/lukas-arnold/garden-equipment-log/internal/utils"
)

func (h *Handler) HandleAddDeviceOperationGet(
	w http.ResponseWriter,
	r *http.Request,
) {
	deviceId, err := utils.ConvertId(
		r.PathValue("deviceId"),
	)

	if err != nil {
		handleError(
			w,
			err,
			http.StatusInternalServerError,
		)
		return
	}

	h.renderTemplate(
		w,
		"templates/device/addOperation.html",
		struct {
			DeviceId int64
		}{
			DeviceId: deviceId,
		},
	)
}

func (h *Handler) HandleAddDeviceOperationPost(
	w http.ResponseWriter,
	r *http.Request,
) {
	deviceId, err := utils.ConvertId(
		r.PathValue("deviceId"),
	)

	if err != nil {
		handleError(
			w,
			err,
			http.StatusInternalServerError,
		)
		return
	}

	err = h.storage.AddDeviceOperation(
		deviceId,
		models.DeviceOperationInput{
			StartTime: r.FormValue("startTime"),
			EndTime:   r.FormValue("endTime"),
			Note:      r.FormValue("note"),
		},
	)

	if err != nil {
		handleError(
			w,
			err,
			http.StatusInternalServerError,
		)
		return
	}

	http.Redirect(
		w,
		r,
		fmt.Sprintf(
			"/device/history/%d",
			deviceId,
		),
		http.StatusFound,
	)
}

func (h *Handler) HandleEditDeviceOperation(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := utils.ConvertId(
		r.PathValue("id"),
	)

	if err != nil {
		handleError(
			w,
			err,
			http.StatusInternalServerError,
		)
		return
	}

	operation, err := h.storage.GetDeviceOperation(id)

	if err != nil {
		handleError(
			w,
			err,
			http.StatusNotFound,
		)
		return
	}

	deviceId, err := h.storage.GetDeviceIdByOperationId(id)

	if err != nil {
		handleError(
			w,
			err,
			http.StatusNotFound,
		)
		return
	}

	h.renderTemplate(
		w,
		"templates/device/editOperation.html",
		struct {
			DeviceId  int64
			Operation models.DeviceOperation
		}{
			DeviceId:  deviceId,
			Operation: operation,
		},
	)
}

func (h *Handler) HandleSaveDeviceOperation(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := utils.ConvertId(
		r.PathValue("id"),
	)

	if err != nil {
		handleError(
			w,
			err,
			http.StatusInternalServerError,
		)
		return
	}

	deviceId, err := h.storage.GetDeviceIdByOperationId(id)

	if err != nil {
		handleError(
			w,
			err,
			http.StatusNotFound,
		)
		return
	}

	err = h.storage.UpdateDeviceOperation(
		id,
		models.DeviceOperationInput{
			StartTime: r.FormValue("startTime"),
			EndTime:   r.FormValue("endTime"),
			Note:      r.FormValue("note"),
		},
	)

	if err != nil {
		handleError(
			w,
			err,
			http.StatusInternalServerError,
		)
		return
	}

	http.Redirect(
		w,
		r,
		fmt.Sprintf(
			"/device/history/%d",
			deviceId,
		),
		http.StatusFound,
	)
}

func (h *Handler) HandleDeleteDeviceOperation(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := utils.ConvertId(
		r.PathValue("id"),
	)

	if err != nil {
		handleError(
			w,
			err,
			http.StatusInternalServerError,
		)
		return
	}

	deviceId, lookupErr :=
		h.storage.GetDeviceIdByOperationId(id)

	err = h.storage.DeleteDeviceOperation(id)

	if err != nil {
		handleError(
			w,
			err,
			http.StatusNotFound,
		)
		return
	}

	redirectTarget := "/devices"

	if lookupErr == nil {
		redirectTarget = fmt.Sprintf(
			"/device/history/%d",
			deviceId,
		)
	}

	http.Redirect(
		w,
		r,
		redirectTarget,
		http.StatusFound,
	)
}
