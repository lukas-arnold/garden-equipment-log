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

func HandleAddBottleOperationPost(w http.ResponseWriter, r *http.Request) {
	bottleId, err := utils.ConvertId(r.PathValue("bottleId"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
	weight, err := utils.ConvertFloat(r.FormValue("weight"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
	err = storage.AddBottleOperation(bottleId, models.BottleOperationInput{
		Date:   r.FormValue("date"),
		Weight: weight,
	})
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
	http.Redirect(w, r, fmt.Sprintf("/bottle/history/%d", bottleId), http.StatusFound)
}

func HandleAddBottleOperationGet(w http.ResponseWriter, r *http.Request) {
	bottleId, err := utils.ConvertId(r.PathValue("bottleId"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	tmpl := template.Must(
		template.New("addOperation.html").Funcs(getTemplateFuncs()).ParseFS(configs.GetWebFiles(), "templates/bottle/addOperation.html"),
	)
	err = tmpl.Execute(w, struct {
		BottleId int64
	}{BottleId: bottleId})
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
}

func HandleDeleteBottleOperation(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
	bottleId, lookupErr := storage.GetBottleIdByOperationId(id)
	err = storage.DeleteBottleOperation(id)
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
	}
	redirectTarget := "/bottles"
	if lookupErr == nil {
		redirectTarget = fmt.Sprintf("/bottle/history/%d", bottleId)
	}
	http.Redirect(w, r, redirectTarget, http.StatusFound)
}

func HandleEditBottleOperationGet(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	operation, err := storage.GetBottleOperation(id)
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
		return
	}
	bottleId, err := storage.GetBottleIdByOperationId(id)
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
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
		errorHandling(w, 500)
		log.Print(err)
	}
}

func HandleSaveBottleOperationPost(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	bottleId, err := storage.GetBottleIdByOperationId(id)
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
		return
	}
	weight, err := utils.ConvertFloat(r.FormValue("weight"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	err = storage.UpdateBottleOperation(id, models.BottleOperationInput{
		Date:   r.FormValue("date"),
		Weight: weight,
	})
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/bottle/history/%d", bottleId), http.StatusFound)
}
