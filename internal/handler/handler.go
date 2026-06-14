package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/lukas-arnold/garden-equipment-log/internal/configs"
	"github.com/lukas-arnold/garden-equipment-log/internal/language"
	"github.com/lukas-arnold/garden-equipment-log/internal/models"
	"github.com/lukas-arnold/garden-equipment-log/internal/storage"
	"github.com/lukas-arnold/garden-equipment-log/internal/utils"
)

func errorHandling(w http.ResponseWriter, httpStatusCode int) {
	w.WriteHeader(httpStatusCode)
}

type deviceOperationRow struct {
	models.DeviceOperation
	Minutes float64
}

type deviceHistoryPageData struct {
	Device models.Device
	Chart  models.ChartModel

	OperationRows []deviceOperationRow
}

type bottleOperationMetric struct {
	models.BottleOperation
	RestGas float64
	UsedGas float64
}

type bottleHistoryPageData struct {
	Bottle models.Bottle
	Chart  models.ChartModel

	OperationRows   []bottleOperationMetric
	TotalUsedGas    float64
	TotalRestGas    float64
	FillLevel       float64
	TotalOperations int
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

func parseDateTime(value string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04", value)
}

func bottleMinMaxWeights(bottle models.Bottle) (float64, float64) {
	// Interpret fields as:
	// - InitialWeight: gross weight when bottle is full (bottle + gas)
	// - FillingWeight: mass of the gas when full
	// Then the empty bottle weight = InitialWeight - FillingWeight
	emptyWeight := bottle.InitialWeight - bottle.FillingWeight
	fullGasWeight := bottle.FillingWeight
	return emptyWeight, fullGasWeight
}

func bottleContentFromWeight(weight, minW, maxW float64) float64 {
	// Here minW is the empty bottle weight, maxW is the gas capacity (mass of gas when full)
	if maxW <= 0 {
		return 0
	}
	// Convert measured gross weight to gas content: currentGas = weight - emptyBottle
	currentGas := weight - minW
	if currentGas < 0 {
		currentGas = 0
	}
	if currentGas > maxW {
		currentGas = maxW
	}
	return currentGas
}

func buildBottleHistoryData(bottle models.Bottle) ([]bottleOperationMetric, float64, float64, float64) {
	operations := append([]models.BottleOperation(nil), bottle.OperationHistory...)

	sort.SliceStable(operations, func(i, j int) bool {
		return operations[i].Date < operations[j].Date
	})

	emptyBottleWeight, fillingWeight := bottleMinMaxWeights(bottle)

	rows := []bottleOperationMetric{}
	prevWeight := fillingWeight

	for _, op := range operations {
		weight := bottleContentFromWeight(op.Weight, emptyBottleWeight, fillingWeight)

		usedGas := 0.0
		if weight < prevWeight {
			usedGas = prevWeight - weight
		}

		rows = append(rows, bottleOperationMetric{
			BottleOperation: op,
			RestGas:         weight,
			UsedGas:         usedGas,
		})

		prevWeight = weight
	}

	latestweight := 0.0
	if len(operations) > 0 {
		last := operations[len(operations)-1]
		latestweight = bottleContentFromWeight(last.Weight, emptyBottleWeight, fillingWeight)
	}

	totalUsed := fillingWeight - latestweight
	if totalUsed < 0 {
		totalUsed = 0
	}

	sort.SliceStable(rows, func(i, j int) bool {
		return rows[i].Date > rows[j].Date
	})

	return rows, totalUsed, latestweight, emptyBottleWeight
}

func HandleFiles(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/web/", http.FileServerFS(configs.GetWebFiles())).ServeHTTP(w, r)
}

func HandleServiceWorker(w http.ResponseWriter, r *http.Request) {
	http.ServeFileFS(w, r, configs.GetWebFiles(), "service-worker.js")
}

func HandleView(w http.ResponseWriter, r *http.Request) {
	handleDevicesView(w, r)
}

func HandleDevicesView(w http.ResponseWriter, r *http.Request) {
	handleDevicesView(w, r)
}

func buildDeviceHistoryRows(device models.Device) []deviceOperationRow {

	rows := make([]deviceOperationRow, 0, len(device.OperationHistory))

	for _, op := range device.OperationHistory {

		var minutes float64

		if op.StartTime != "" && op.EndTime != "" {

			start, err1 := parseDateTime(op.StartTime)
			end, err2 := parseDateTime(op.EndTime)

			if err1 == nil && err2 == nil {
				minutes = end.Sub(start).Minutes()

				if minutes < 0 {
					minutes = 0
				}
			}
		}

		rows = append(rows, deviceOperationRow{
			DeviceOperation: op,
			Minutes:         minutes,
		})
	}

	return rows
}

func handleDevicesView(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("index.html").Funcs(getTemplateFuncs()).ParseFS(configs.GetWebFiles(), "templates/device/index.html"),
	)
	devices, err := storage.GetDevices()
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
	}
	err = tmpl.Execute(w, devices)
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
}

func HandleBottlesView(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("index.html").Funcs(getTemplateFuncs()).ParseFS(configs.GetWebFiles(), "templates/bottle/index.html"),
	)
	bottles, err := storage.GetBottles()
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
	}
	err = tmpl.Execute(w, bottles)
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
}

// New wrappers and compatibility handlers expected by cmd/main.go
func HandleAddDeviceOperationGet(w http.ResponseWriter, r *http.Request) {
	deviceId, err := utils.ConvertId(r.PathValue("deviceId"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	tmpl := template.Must(
		template.New("addOperation.html").Funcs(getTemplateFuncs()).ParseFS(configs.GetWebFiles(), "templates/device/addOperation.html"),
	)
	err = tmpl.Execute(w, struct{ DeviceId int64 }{DeviceId: deviceId})
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
}

func HandleEditDeviceOperation(w http.ResponseWriter, r *http.Request) {
	HandleEditDeviceOperationGet(w, r)
}

func HandleSaveDeviceOperation(w http.ResponseWriter, r *http.Request) {
	HandleSaveDeviceOperationPost(w, r)
}

func HandleDeviceHistory(w http.ResponseWriter, r *http.Request) {

	deviceId := r.PathValue("id")

	id, err := utils.ConvertId(deviceId)
	if err != nil {
		errorHandling(w, 500)
		return
	}

	device, err := storage.GetDevice(id)
	if err != nil {
		errorHandling(w, 404)
		return
	}

	chart, err := storage.GetDeviceChart(id)
	if err != nil {
		errorHandling(w, 500)
		return
	}

	rows := buildDeviceHistoryRows(device)

	view := deviceHistoryPageData{
		Device:        device,
		Chart:         chart,
		OperationRows: rows,
	}

	tmpl := template.Must(
		template.New("history.html").
			Funcs(getTemplateFuncs()).
			ParseFS(
				configs.GetWebFiles(),
				"templates/device/history.html",
			),
	)

	err = tmpl.Execute(w, view)
	if err != nil {
		errorHandling(w, 500)
		return
	}
}

func HandleEditBottleOperation(w http.ResponseWriter, r *http.Request) {
	HandleEditBottleOperationGet(w, r)
}

func HandleSaveBottleOperation(w http.ResponseWriter, r *http.Request) {
	HandleSaveBottleOperationPost(w, r)
}

func HandleBottleHistory(w http.ResponseWriter, r *http.Request) {
	bottleId := r.PathValue("id")

	id, err := utils.ConvertId(bottleId)
	if err != nil {
		errorHandling(w, 500)
		return
	}

	bottle, err := storage.GetBottle(id)
	if err != nil {
		errorHandling(w, 404)
		return
	}

	chart, err := storage.GetBottleChart(id)
	if err != nil {
		errorHandling(w, 500)
		return
	}

	rows, totalUsed, totalRest, _ := buildBottleHistoryData(bottle) // totalRest is not used in the view, but could be added if desired

	fillLevel := (totalRest / bottle.FillingWeight) * 100

	view := bottleHistoryPageData{
		Bottle:          bottle,
		Chart:           chart,
		OperationRows:   rows,
		TotalUsedGas:    totalUsed,
		TotalRestGas:    totalRest,
		FillLevel:       fillLevel,
		TotalOperations: len(rows),
	}

	fmt.Println(rows)

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
		errorHandling(w, 500)
		return
	}
}
