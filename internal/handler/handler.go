package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/lukas-arnold/garden-equipment-log/internal/configs"
	"github.com/lukas-arnold/garden-equipment-log/internal/language"
	"github.com/lukas-arnold/garden-equipment-log/internal/storage"
)

type Handler struct {
	storage *storage.Storage
}

func New(storage *storage.Storage) *Handler {
	return &Handler{storage: storage}
}

func handleError(w http.ResponseWriter, err error, statusCode int) {
	log.Printf("HTTP %d: %v", statusCode, err)
	http.Error(w, http.StatusText(statusCode), statusCode)
}

func (h *Handler) renderTemplate(w http.ResponseWriter, file string, data any) {
	tmpl := template.Must(template.New("base.html").
		Funcs(getTemplateFuncs()).
		ParseFS(configs.GetWebFiles(), "templates/base.html", file))

	if err := tmpl.Execute(w, data); err != nil {
		handleError(w, err, http.StatusInternalServerError)
	}
}

func (h *Handler) renderHistoryTemplate(w http.ResponseWriter, file string, data any) {
	tmpl := template.Must(template.New("baseHistory.html").
		Funcs(getTemplateFuncs()).
		ParseFS(configs.GetWebFiles(), "templates/baseHistory.html", file))

	if err := tmpl.Execute(w, data); err != nil {
		handleError(w, err, http.StatusInternalServerError)
	}
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

		"formatOperationTime": func(value float64) string {
			totalMinutes := int64(value)
			hours := totalMinutes / 60
			minutes := totalMinutes % 60

			if hours > 0 {
				if minutes == 0 {
					return fmt.Sprintf("%d h", hours)
				}
				return fmt.Sprintf("%d h %d min", hours, minutes)
			}
			return fmt.Sprintf("%d min", minutes)
		},
	}
}

func (h *Handler) HandleView(w http.ResponseWriter, r *http.Request) {
	h.HandleDevicesView(w, r)
}
