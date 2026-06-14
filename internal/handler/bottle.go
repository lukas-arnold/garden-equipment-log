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

type bottleFormData struct {
	Title       string
	Action      string
	SubmitLabel string
	Bottle      models.Bottle
}

type bottleOperationRow struct {
	models.BottleOperation
	RestGas float64
	UsedGas float64
}

type bottleHistoryPageData struct {
	Bottle        models.Bottle
	Chart         models.ChartModel
	OperationRows []bottleOperationRow
	FillLevel     float64
}

func HandleAddBottleGet(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("add.html").Funcs(template.FuncMap{
			"T": func(key string) string {
				return language.T(configs.GetLanguage(), key)
			},
		}).ParseFS(configs.GetWebFiles(), "templates/bottle/add.html"),
	)
	err := tmpl.Execute(w, bottleFormData{
		Title:       language.T(configs.GetLanguage(), "addBottle"),
		Action:      "/bottle/add",
		SubmitLabel: language.T(configs.GetLanguage(), "save"),
		Bottle:      models.Bottle{},
	})
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
}

func HandleAddBottlePost(w http.ResponseWriter, r *http.Request) {
	purchasePrice, err := utils.ConvertFloat(r.FormValue("purchasePrice"))
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	initialWeight, err := utils.ConvertFloat(r.FormValue("initialWeight"))
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	fillingWeight, err := utils.ConvertFloat(r.FormValue("fillingWeight"))
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	err = storage.AddBottle(models.BottleInput{
		PurchaseDate:  r.FormValue("purchaseDate"),
		PurchasePrice: purchasePrice,
		InitialWeight: initialWeight,
		FillingWeight: fillingWeight,
	})
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/bottles", http.StatusFound)
}

func HandleEditBottle(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("add.html").Funcs(template.FuncMap{
			"T": func(key string) string {
				return language.T(configs.GetLanguage(), key)
			},
		}).ParseFS(configs.GetWebFiles(), "templates/bottle/add.html"),
	)
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	bottle, err := storage.GetBottle(id)
	if err != nil {
		handleError(w, err, http.StatusNotFound)
		return
	}
	err = tmpl.Execute(w, bottleFormData{
		Title:       language.T(configs.GetLanguage(), "editBottle"),
		Action:      fmt.Sprintf("/bottle/save/%d", bottle.Id),
		SubmitLabel: language.T(configs.GetLanguage(), "save"),
		Bottle:      bottle,
	})
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
}

func HandleSaveBottle(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	bottle, err := storage.GetBottle(id)
	if err != nil {
		handleError(w, err, http.StatusNotFound)
		return
	}
	purchasePrice, err := utils.ConvertFloat(r.FormValue("purchasePrice"))
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	initialWeight, err := utils.ConvertFloat(r.FormValue("initialWeight"))
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	fillingWeight, err := utils.ConvertFloat(r.FormValue("fillingWeight"))
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	bottle.PurchaseDate = r.FormValue("purchaseDate")
	bottle.PurchasePrice = purchasePrice
	bottle.InitialWeight = initialWeight
	bottle.FillingWeight = fillingWeight
	err = storage.UpdateBottle(bottle)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/bottles", http.StatusFound)
}

func HandleDeleteBottle(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	err = storage.DeleteBottle(id)
	if err != nil {
		handleError(w, err, http.StatusNotFound)
		return
	}
	http.Redirect(w, r, "/bottles", http.StatusFound)
}

func HandleBottleHistory(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	bottle, err := storage.GetBottle(id)
	if err != nil {
		handleError(w, err, http.StatusNotFound)
		return
	}
	chart, err := storage.GetBottleChart(id)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	rows := buildBottleHistoryRows(bottle)
	fillLevel := (bottle.RestGas / bottle.FillingWeight) * 100
	view := bottleHistoryPageData{
		Bottle:        bottle,
		Chart:         chart,
		OperationRows: rows,
		FillLevel:     fillLevel,
	}
	tmpl := template.Must(
		template.New("history.html").
			Funcs(getTemplateFuncs()).
			ParseFS(
				configs.GetWebFiles(),
				"templates/bottle/history.html",
			),
	)
	err = tmpl.Execute(w, view)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
}

func buildBottleHistoryRows(bottle models.Bottle) []bottleOperationRow {
	rows := make([]bottleOperationRow, 0, len(bottle.OperationHistory))
	for _, op := range bottle.OperationHistory {
		usedGas := bottle.InitialWeight - op.Weight
		restGas := bottle.FillingWeight - usedGas
		rows = append(rows, bottleOperationRow{
			BottleOperation: op,
			RestGas:         restGas,
			UsedGas:         usedGas,
		})
	}
	return rows
}
