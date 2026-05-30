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

type deviceFormData struct {
	Title       string
	Action      string
	SubmitLabel string
	Device      models.Device
}

func HandleAddDeviceGet(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("add.html").Funcs(template.FuncMap{
			"T": func(key string) string {
				return language.T(configs.GetLanguage(), key)
			},
		}).ParseFS(configs.GetWebFiles(), "templates/device/add.html"),
	)
	err := tmpl.Execute(w, deviceFormData{
		Title:       language.T(configs.GetLanguage(), "addDevice"),
		Action:      "/device/add",
		SubmitLabel: language.T(configs.GetLanguage(), "save"),
		Device:      models.Device{},
	})
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
}

func HandleAddDevicePost(w http.ResponseWriter, r *http.Request) {
	purchasePrice, err := utils.ConvertFloat(r.FormValue("purchasePrice"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
	err = storage.AddDevice(models.DeviceInput{
		Name:          r.FormValue("name"),
		PurchaseDate:  r.FormValue("purchaseDate"),
		PurchasePrice: purchasePrice,
	})
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
	http.Redirect(w, r, "/devices", http.StatusFound)
}

func HandleEditDevice(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("add.html").Funcs(template.FuncMap{
			"T": func(key string) string {
				return language.T(configs.GetLanguage(), key)
			},
		}).ParseFS(configs.GetWebFiles(), "templates/device/add.html"),
	)
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
	device, err := storage.GetDevice(id)
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
	}
	err = tmpl.Execute(w, deviceFormData{
		Title:       language.T(configs.GetLanguage(), "editDevice"),
		Action:      fmt.Sprintf("/device/save/%d", device.Id),
		SubmitLabel: language.T(configs.GetLanguage(), "save"),
		Device:      device,
	})
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
}

func HandleSaveDevice(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
	device, err := storage.GetDevice(id)
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
	}
	purchasePrice, err := utils.ConvertFloat(r.FormValue("purchasePrice"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
	device.Name = r.FormValue("name")
	device.PurchaseDate = r.FormValue("purchaseDate")
	device.PurchasePrice = purchasePrice
	err = storage.UpdateDevice(device)
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
	http.Redirect(w, r, "/devices", http.StatusFound)
}

func HandleDeleteDevice(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
	err = storage.DeleteDevice(id)
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
	}
	http.Redirect(w, r, "/devices", http.StatusFound)
}
