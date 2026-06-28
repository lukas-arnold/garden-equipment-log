package handler

import (
	"fmt"
	"net/http"

	"github.com/lukas-arnold/garden-equipment-log/internal/configs"
	"github.com/lukas-arnold/garden-equipment-log/internal/language"
	"github.com/lukas-arnold/garden-equipment-log/internal/models"
	"github.com/lukas-arnold/garden-equipment-log/internal/utils"
)

type bottleFormData struct {
	Title       string
	Action      string
	SubmitLabel string
	Bottle      models.Bottle
}

type bottleHistoryPageData struct {
	Bottle        models.Bottle
	Chart         models.ChartModel
	OperationRows []bottleOperationRow
	FillLevel     float64
}

func (h *Handler) HandleBottlesView(
	w http.ResponseWriter,
	r *http.Request,
) {
	bottles, err := h.storage.GetBottles()

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
		"templates/bottle/index.html",
		bottles,
	)
}

func (h *Handler) HandleAddBottleGet(
	w http.ResponseWriter,
	r *http.Request,
) {
	h.renderTemplate(
		w,
		"templates/bottle/add.html",
		bottleFormData{
			Title: language.T(
				configs.GetLanguage(),
				"addBottle",
			),
			Action: "/bottle/add",
			SubmitLabel: language.T(
				configs.GetLanguage(),
				"save",
			),
			Bottle: models.Bottle{},
		},
	)
}

func (h *Handler) HandleAddBottlePost(
	w http.ResponseWriter,
	r *http.Request,
) {
	purchasePrice, err := utils.ConvertFloat(
		r.FormValue("purchasePrice"),
	)

	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	initialWeight, err := utils.ConvertFloat(
		r.FormValue("initialWeight"),
	)

	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	fillingWeight, err := utils.ConvertFloat(
		r.FormValue("fillingWeight"),
	)

	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	err = h.storage.AddBottle(
		models.BottleInput{
			PurchaseDate:  r.FormValue("purchaseDate"),
			PurchasePrice: purchasePrice,
			InitialWeight: initialWeight,
			FillingWeight: fillingWeight,
		},
	)

	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(
		w,
		r,
		"/bottles",
		http.StatusFound,
	)
}

func (h *Handler) HandleEditBottle(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := utils.ConvertId(
		r.PathValue("id"),
	)

	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	bottle, err := h.storage.GetBottle(id)

	if err != nil {
		handleError(w, err, http.StatusNotFound)
		return
	}

	h.renderTemplate(
		w,
		"templates/bottle/add.html",
		bottleFormData{
			Title: language.T(
				configs.GetLanguage(),
				"editBottle",
			),

			Action: fmt.Sprintf(
				"/bottle/save/%d",
				bottle.Id,
			),

			SubmitLabel: language.T(
				configs.GetLanguage(),
				"save",
			),

			Bottle: bottle,
		},
	)
}

func (h *Handler) HandleSaveBottle(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := utils.ConvertId(
		r.PathValue("id"),
	)

	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	bottle, err := h.storage.GetBottle(id)

	if err != nil {
		handleError(w, err, http.StatusNotFound)
		return
	}

	purchasePrice, err := utils.ConvertFloat(
		r.FormValue("purchasePrice"),
	)

	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	initialWeight, err := utils.ConvertFloat(
		r.FormValue("initialWeight"),
	)

	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	fillingWeight, err := utils.ConvertFloat(
		r.FormValue("fillingWeight"),
	)

	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	bottle.PurchaseDate = r.FormValue("purchaseDate")
	bottle.PurchasePrice = purchasePrice
	bottle.InitialWeight = initialWeight
	bottle.FillingWeight = fillingWeight

	err = h.storage.UpdateBottle(bottle)

	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(
		w,
		r,
		"/bottles",
		http.StatusFound,
	)
}

func (h *Handler) HandleDeleteBottle(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := utils.ConvertId(
		r.PathValue("id"),
	)

	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	err = h.storage.DeleteBottle(id)

	if err != nil {
		handleError(w, err, http.StatusNotFound)
		return
	}

	http.Redirect(
		w,
		r,
		"/bottles",
		http.StatusFound,
	)
}

func (h *Handler) HandleBottleHistory(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := utils.ConvertId(
		r.PathValue("id"),
	)

	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	bottle, err := h.storage.GetBottle(id)

	if err != nil {
		handleError(w, err, http.StatusNotFound)
		return
	}

	chart, err := h.storage.GetBottleChart(id)

	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	view := bottleHistoryPageData{
		Bottle:        bottle,
		Chart:         chart,
		OperationRows: buildBottleHistoryRows(bottle),
		FillLevel:     (bottle.RestGas / bottle.FillingWeight) * 100,
	}

	h.renderHistoryTemplate(
		w,
		"templates/bottle/history.html",
		view,
	)
}
