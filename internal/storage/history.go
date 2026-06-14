package storage

import (
	"fmt"
	"sort"
	"time"

	"github.com/lukas-arnold/garden-equipment-log/internal/models"
)

func parseDateTime(value string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04", value)
}

func GetDeviceChart(deviceId int64) (models.ChartModel, error) {

	device, err := GetDevice(deviceId)
	if err != nil {
		return models.ChartModel{}, err
	}

	yearHours := make(map[int]float64)

	for _, op := range device.OperationHistory {

		if op.StartTime == "" || op.EndTime == "" {
			continue
		}

		start, err := parseDateTime(op.StartTime)
		if err != nil {
			continue
		}

		end, err := parseDateTime(op.EndTime)
		if err != nil {
			continue
		}

		hours := end.Sub(start).Hours()

		if hours < 0 {
			hours = 0
		}

		yearHours[start.Year()] += hours
	}

	var years []int

	for year := range yearHours {
		years = append(years, year)
	}

	sort.Ints(years)

	var labels []string
	var values []float64

	for _, year := range years {
		labels = append(labels, fmt.Sprintf("%d", year))
		values = append(values, yearHours[year])
	}

	return models.ChartModel{
		Type:   "bar",
		Labels: labels,
		Sets: []models.ChartDataset{
			{
				Label: "Operation Hours",
				Unit:  "h",
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
				Label: "Weight",
				Unit:  "kg",
				Data:  weights,
			},
			{
				Label: "Empty Weight",
				Unit:  "kg",
				Data:  emptyWeights,
			},
		},
	}, nil
}
