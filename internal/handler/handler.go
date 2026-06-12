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

type deviceOperationsPageData struct {
	Device           models.Device
	Dates            []string
	Durations        []float64
	OperationMinutes []float64
}

type bottleOperationMetric struct {
	models.BottleOperation
	RestGas float64
	UsedGas float64
}

type bottleHistoryPageData struct {
	Bottle          models.Bottle
	Dates           []string
	Weights         []float64
	EmptyWeight     float64
	OperationRows   []bottleOperationMetric
	TotalUsedGas    float64
	TotalRestGas    float64
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

func parseDate(value string) (time.Time, error) {
	return time.Parse("2006-01-02", value)
}

func parseDateTime(value string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04", value)
}

func calculateDeviceOperationMinutes(device models.Device) []float64 {
	minutes := make([]float64, len(device.OperationHistory))
	for i, operation := range device.OperationHistory {
		if operation.StartTime == "" || operation.EndTime == "" {
			minutes[i] = 0
			continue
		}
		start, err := parseDateTime(operation.StartTime)
		if err != nil {
			minutes[i] = 0
			continue
		}
		end, err := parseDateTime(operation.EndTime)
		if err != nil {
			minutes[i] = 0
			continue
		}
		delta := end.Sub(start).Minutes()
		if delta < 0 {
			delta = 0
		}
		minutes[i] = delta
	}
	return minutes
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

func handleDevicesView(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("index.html").Funcs(getTemplateFuncs()).ParseFS(configs.GetWebFiles(), "templates/device/index.html"),
	)
	equipmentStorage, err := storage.GetEquipmentStorage()
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
	}
	err = tmpl.Execute(w, equipmentStorage)
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
}

func HandleBottlesView(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("index.html").Funcs(getTemplateFuncs()).ParseFS(configs.GetWebFiles(), "templates/bottle/index.html"),
	)
	equipmentStorage, err := storage.GetEquipmentStorage()
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
	}
	err = tmpl.Execute(w, equipmentStorage)
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
}

func HandleDevicesChart(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("chart.html").Funcs(getTemplateFuncs()).ParseFS(configs.GetWebFiles(), "templates/device/chart.html"),
	)
	chart, err := storage.GetDevicesForChart()
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
	}
	err = tmpl.Execute(w, chart)
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
}

func HandleBottlesChart(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("chart.html").Funcs(getTemplateFuncs()).ParseFS(configs.GetWebFiles(), "templates/bottle/chart.html"),
	)
	chart, err := storage.GetBottlesForChart()
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
	}
	err = tmpl.Execute(w, chart)
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
}

func HandleDeviceOperationChart(w http.ResponseWriter, r *http.Request) {
	deviceId := r.PathValue("deviceId")
	id, err := utils.ConvertId(deviceId)
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	chart, err := storage.GetDeviceOperationTimesForChart(id, 0)
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
		return
	}
	tmpl := template.Must(
		template.New("operationchart.html").Funcs(getTemplateFuncs()).ParseFS(configs.GetWebFiles(), "templates/device/operationchart.html"),
	)
	err = tmpl.Execute(w, chart)
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
}

func HandleDeviceOperationsView(w http.ResponseWriter, r *http.Request) {
	deviceId := r.PathValue("deviceId")
	id, err := utils.ConvertId(deviceId)
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	device, err := storage.GetDevice(id)
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
		return
	}
	chart, err := storage.GetDeviceOperationTimesForChart(id, 0)
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
		return
	}
	tmpl := template.Must(
		template.New("operations.html").Funcs(getTemplateFuncs()).ParseFS(configs.GetWebFiles(), "templates/device/operations.html"),
	)
	err = tmpl.Execute(w, deviceOperationsPageData{Device: device, Dates: chart.Dates, Durations: chart.Durations})
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
}

func HandleBottleWeightChart(w http.ResponseWriter, r *http.Request) {
	bottleId, err := utils.ConvertId(r.PathValue("bottleId"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	chart, err := storage.GetBottleWeightHistoryForChart(bottleId)
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
		return
	}
	tmpl := template.Must(
		template.New("weightchart.html").Funcs(getTemplateFuncs()).ParseFS(configs.GetWebFiles(), "templates/bottle/weightchart.html"),
	)
	err = tmpl.Execute(w, chart)
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
}

// New wrappers and compatibility handlers expected by cmd/main.go
func HandleAddDeviceGetOperation(w http.ResponseWriter, r *http.Request) {
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

func HandleAddDevicePostOperation(w http.ResponseWriter, r *http.Request) {
	HandleAddDeviceOperationPost(w, r)
}

func HandleEditDeviceOperation(w http.ResponseWriter, r *http.Request) {
	HandleEditDeviceOperationGet(w, r)
}

func HandleSaveDeviceOperation(w http.ResponseWriter, r *http.Request) {
	HandleSaveDeviceOperationPost(w, r)
}

func HandleDeviceHistory(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	device, err := storage.GetDevice(id)
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
		return
	}
	chart, err := storage.GetDeviceOperationMinutesPerYearForChart(id)
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
		return
	}
	tmpl := template.Must(
		template.New("history.html").Funcs(getTemplateFuncs()).ParseFS(configs.GetWebFiles(), "templates/device/history.html"),
	)
	view := deviceOperationsPageData{
		Device:           device,
		Dates:            chart.Dates,
		Durations:        chart.Durations,
		OperationMinutes: calculateDeviceOperationMinutes(device),
	}
	err = tmpl.Execute(w, view)
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
}

func HandleEditBottleOperation(w http.ResponseWriter, r *http.Request) {
	HandleEditBottleOperationGet(w, r)
}

func HandleSaveBottleOperation(w http.ResponseWriter, r *http.Request) {
	HandleSaveBottleOperationPost(w, r)
}

func HandleBottleHistory(w http.ResponseWriter, r *http.Request) {
	bottleId, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	chart, err := storage.GetBottleWeightHistoryForChart(bottleId)
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
		return
	}
	rows, totalUsed, totalRest, emptyWeight := buildBottleHistoryData(chart.Bottle)
	tmpl := template.Must(
		template.New("history.html").Funcs(getTemplateFuncs()).ParseFS(configs.GetWebFiles(), "templates/bottle/history.html"),
	)
	view := bottleHistoryPageData{
		Bottle:          chart.Bottle,
		Dates:           chart.Dates,
		Weights:         chart.Weights,
		EmptyWeight:     emptyWeight,
		OperationRows:   rows,
		TotalUsedGas:    totalUsed,
		TotalRestGas:    totalRest,
		TotalOperations: len(rows),
	}
	err = tmpl.Execute(w, view)
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
}
