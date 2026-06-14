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

func HandleBottlesView(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("index.html").Funcs(getTemplateFuncs()).ParseFS(configs.GetWebFiles(), "templates/bottle/index.html"),
	)
	bottles, err := storage.GetBottles()
	if err != nil {
		handleError(w, err, http.StatusNotFound)
		return
	}
	err = tmpl.Execute(w, bottles)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
}

func HandleAddBottleOperationGet(w http.ResponseWriter, r *http.Request) {
	bottleId, err := utils.ConvertId(r.PathValue("bottleId"))
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	tmpl := template.Must(
		template.New("addOperation.html").Funcs(getTemplateFuncs()).ParseFS(configs.GetWebFiles(), "templates/bottle/addOperation.html"),
	)
	err = tmpl.Execute(w, struct {
		BottleId int64
	}{BottleId: bottleId})
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
}

func HandleAddBottleOperationPost(w http.ResponseWriter, r *http.Request) {
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
	err = storage.AddBottleOperation(bottleId, models.BottleOperationInput{
		Date:   r.FormValue("date"),
		Weight: weight,
	})
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/bottle/history/%d", bottleId), http.StatusFound)
}

func HandleEditBottleOperation(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	operation, err := storage.GetBottleOperation(id)
	if err != nil {
		handleError(w, err, http.StatusNotFound)
		return
	}
	bottleId, err := storage.GetBottleIdByOperationId(id)
	if err != nil {
		handleError(w, err, http.StatusNotFound)
		return
	}
	tmpl := template.Must(
		template.New("editOperation.html").Funcs(getTemplateFuncs()).ParseFS(configs.GetWebFiles(), "templates/bottle/editOperation.html"),
	)
	err = tmpl.Execute(w, struct {
		BottleId  int64
		Operation models.BottleOperation
	}{BottleId: bottleId, Operation: operation})
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
}

func HandleSaveBottleOperation(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	bottleId, err := storage.GetBottleIdByOperationId(id)
	if err != nil {
		handleError(w, err, http.StatusNotFound)
		return
	}
	weight, err := utils.ConvertFloat(r.FormValue("weight"))
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	err = storage.UpdateBottleOperation(id, models.BottleOperationInput{
		Date:   r.FormValue("date"),
		Weight: weight,
	})
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/bottle/history/%d", bottleId), http.StatusFound)
}

func HandleDeleteBottleOperation(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	bottleId, lookupErr := storage.GetBottleIdByOperationId(id)
	err = storage.DeleteBottleOperation(id)
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
