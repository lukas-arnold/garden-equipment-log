package storage

import (
	"fmt"
	"sort"

	"github.com/lukas-arnold/garden-equipment-log/internal/configs"
	"github.com/lukas-arnold/garden-equipment-log/internal/language"
	"github.com/lukas-arnold/garden-equipment-log/internal/models"
	"github.com/lukas-arnold/garden-equipment-log/internal/utils"
)

func GetDeviceChart(deviceId int64) (models.ChartModel, error) {
	device, err := GetDevice(deviceId)
	if err != nil {
		return models.ChartModel{}, err
	}
	yearMinutes := make(map[int]float64)
	for _, op := range device.OperationHistory {
		if op.StartTime == "" || op.EndTime == "" {
			continue
		}
		start, err := utils.ParseDateTime(op.StartTime)
		if err != nil {
			continue
		}
		end, err := utils.ParseDateTime(op.EndTime)
		if err != nil {
			continue
		}
		minutes := end.Sub(start).Minutes()
		if minutes < 0 {
			minutes = 0
		}
		yearMinutes[start.Year()] += minutes
	}
	var years []int
	for year := range yearMinutes {
		years = append(years, year)
	}
	sort.Ints(years)
	var labels []string
	var values []float64
	for _, year := range years {
		labels = append(labels, fmt.Sprintf("%d", year))
		values = append(values, yearMinutes[year])
	}
	return models.ChartModel{
		Type:   "bar",
		Labels: labels,
		Sets: []models.ChartDataset{
			{
				Label: language.T(configs.GetLanguage(), "operationTime"),
				Unit:  "h/min",
				Data:  values,
			},
		},
	}, nil
}

func GetBottleChart(bottleId int64) (models.ChartModel, error) {
	bottle, err := GetBottle(bottleId)
	if err != nil {
		return models.ChartModel{}, err
	}
	ops := append([]models.BottleOperation(nil), bottle.OperationHistory...)
	sort.SliceStable(ops, func(i, j int) bool {
		return ops[i].Date < ops[j].Date
	})
	emptyWeight := bottle.InitialWeight - bottle.FillingWeight
	var labels []string
	labels = append(labels, bottle.PurchaseDate)
	var weights []float64
	weights = append(weights, bottle.InitialWeight)
	var emptyWeights []float64
	emptyWeights = append(emptyWeights, emptyWeight)
	for _, op := range ops {
		labels = append(labels, op.Date)
		weights = append(weights, op.Weight)
		emptyWeights = append(emptyWeights, emptyWeight)
	}
	return models.ChartModel{
		Type:   "line",
		Labels: labels,
		Sets: []models.ChartDataset{
			{
				Label: language.T(configs.GetLanguage(), "weight"),
				Unit:  "kg",
				Data:  weights,
			},
			{
				Label: language.T(configs.GetLanguage(), "emptyWeight"),
				Unit:  "kg",
				Data:  emptyWeights,
			},
		},
	}, nil
}
