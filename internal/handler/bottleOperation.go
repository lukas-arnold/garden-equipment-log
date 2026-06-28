package handler

import (
	"fmt"
	"net/http"

	"github.com/lukas-arnold/garden-equipment-log/internal/models"
	"github.com/lukas-arnold/garden-equipment-log/internal/utils"
)

func (h *Handler) HandleAddBottleOperationGet(w http.ResponseWriter, r *http.Request) {
	bottleId, err := utils.ConvertId(r.PathValue("bottleId"))
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	h.renderTemplate(w, "templates/bottle/addOperation.html", struct{ BottleId int64 }{
		BottleId: bottleId,
	})
}

func (h *Handler) HandleAddBottleOperationPost(w http.ResponseWriter, r *http.Request) {
	bottleId, err := utils.ConvertId(r.PathValue("bottleId"))
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	weight, err := utils.ConvertFloat(r.FormValue("weight"))
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	err = h.storage.AddBottleOperation(bottleId, models.BottleOperationInput{
		Date:   r.FormValue("date"),
		Weight: weight,
	})
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/bottle/history/%d", bottleId), http.StatusFound)
}

func (h *Handler) HandleEditBottleOperation(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	operation, err := h.storage.GetBottleOperation(id)
	if err != nil {
		handleError(w, err, http.StatusNotFound)
		return
	}

	bottleId, err := h.storage.GetBottleIdByOperationId(id)
	if err != nil {
		handleError(w, err, http.StatusNotFound)
		return
	}

	h.renderTemplate(w, "templates/bottle/editOperation.html", struct {
		BottleId  int64
		Operation models.BottleOperation
	}{
		BottleId:  bottleId,
		Operation: operation,
	})
}

func (h *Handler) HandleSaveBottleOperation(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	bottleId, err := h.storage.GetBottleIdByOperationId(id)
	if err != nil {
		handleError(w, err, http.StatusNotFound)
		return
	}

	weight, err := utils.ConvertFloat(r.FormValue("weight"))
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	err = h.storage.UpdateBottleOperation(id, models.BottleOperationInput{
		Date:   r.FormValue("date"),
		Weight: weight,
	})
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/bottle/history/%d", bottleId), http.StatusFound)
}

func (h *Handler) HandleDeleteBottleOperation(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	bottleId, lookupErr := h.storage.GetBottleIdByOperationId(id)
	err = h.storage.DeleteBottleOperation(id)
	if err != nil {
		handleError(w, err, http.StatusNotFound)
		return
	}

	redirectTarget := "/bottles"
	if lookupErr == nil {
		redirectTarget = fmt.Sprintf("/bottle/history/%d", bottleId)
	}

	http.Redirect(w, r, redirectTarget, http.StatusFound)
}
