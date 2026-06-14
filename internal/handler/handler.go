package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/lukas-arnold/garden-equipment-log/internal/configs"
	"github.com/lukas-arnold/garden-equipment-log/internal/language"
)

func handleError(w http.ResponseWriter, err error, statusCode int) {
	log.Printf("HTTP %d: %v", statusCode, err)
	http.Error(w, http.StatusText(statusCode), statusCode)
}

func getTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"T": func(key string) string {
			return language.T(configs.GetLanguage(), key)
		},
		"formatFloat": func(value float64, decimals int) string {
			format := fmt.Sprintf("%%.%df", decimals)
			return strings.ReplaceAll(fmt.Sprintf(format, value), ".", ",")
		},
		"formatEuro": func(value float64) string {
			return strings.ReplaceAll(fmt.Sprintf("%.2f €", value), ".", ",")
		},
	}
}

func HandleFiles(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/web/", http.FileServerFS(configs.GetWebFiles())).ServeHTTP(w, r)
}

func HandleServiceWorker(w http.ResponseWriter, r *http.Request) {
	http.ServeFileFS(w, r, configs.GetWebFiles(), "service-worker.js")
}

func HandleView(w http.ResponseWriter, r *http.Request) {
	HandleDevicesView(w, r)
}
