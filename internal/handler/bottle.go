package handler

import (
	"fmt"
	"html/template"
	"log"
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
		errorHandling(w, 500)
		log.Print(err)
	}
}

func HandleAddBottlePost(w http.ResponseWriter, r *http.Request) {
	purchasePrice, err := utils.ConvertFloat(r.FormValue("purchasePrice"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
	initialWeight, err := utils.ConvertFloat(r.FormValue("initialWeight"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
	fillingWeight, err := utils.ConvertFloat(r.FormValue("fillingWeight"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
	err = storage.AddBottle(models.BottleInput{
		PurchaseDate:  r.FormValue("purchaseDate"),
		PurchasePrice: purchasePrice,
		InitialWeight: initialWeight,
		FillingWeight: fillingWeight,
	})
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
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
		errorHandling(w, 500)
		log.Print(err)
	}
	bottle, err := storage.GetBottle(id)
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
	}
	err = tmpl.Execute(w, bottleFormData{
		Title:       language.T(configs.GetLanguage(), "editBottle"),
		Action:      fmt.Sprintf("/bottle/save/%d", bottle.Id),
		SubmitLabel: language.T(configs.GetLanguage(), "save"),
		Bottle:      bottle,
	})
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
}

func HandleSaveBottle(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
	bottle, err := storage.GetBottle(id)
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
	}
	purchasePrice, err := utils.ConvertFloat(r.FormValue("purchasePrice"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
	initialWeight, err := utils.ConvertFloat(r.FormValue("initialWeight"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
	fillingWeight, err := utils.ConvertFloat(r.FormValue("fillingWeight"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
	bottle.PurchaseDate = r.FormValue("purchaseDate")
	bottle.PurchasePrice = purchasePrice
	bottle.InitialWeight = initialWeight
	bottle.FillingWeight = fillingWeight
	err = storage.UpdateBottle(bottle)
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
	http.Redirect(w, r, "/bottles", http.StatusFound)
}

func HandleDeleteBottle(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
	err = storage.DeleteBottle(id)
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
	}
	http.Redirect(w, r, "/bottles", http.StatusFound)
}
