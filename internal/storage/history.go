package storage

import (
	"fmt"
	"sort"

	"github.com/lukas-arnold/garden-equipment-log/internal/configs"
	"github.com/lukas-arnold/garden-equipment-log/internal/language"
	"github.com/lukas-arnold/garden-equipment-log/internal/models"
	"github.com/lukas-arnold/garden-equipment-log/internal/utils"
)

func (s *Storage) GetDeviceChart(deviceId int64) (models.ChartModel, error) {
	device, err := s.GetDevice(deviceId)
	if err != nil {
		return models.ChartModel{}, err
	}

	yearMinutes := make(map[int]float64)
	for _, op := range device.OperationHistory {
		start, err1 := utils.ParseDateTime(op.StartTime)
		end, err2 := utils.ParseDateTime(op.EndTime)
		if err1 != nil || err2 != nil {
			continue
		}

		yearMinutes[start.Year()] += max(0, end.Sub(start).Minutes())
	}

	years := make([]int, 0, len(yearMinutes))
	for year := range yearMinutes {
		years = append(years, year)
	}
	sort.Ints(years)

	labels := make([]string, len(years))
	values := make([]float64, len(years))
	for i, year := range years {
		labels[i] = fmt.Sprintf("%d", year)
		values[i] = yearMinutes[year]
	}

	return models.ChartModel{
		Type:   "bar",
		Labels: labels,
		Sets: []models.ChartDataset{{
			Label: language.T(configs.GetLanguage(), "operationTime"),
			Unit:  "h/min",
			Data:  values,
		}},
	}, nil
}

func (s *Storage) GetBottleChart(bottleId int64) (models.ChartModel, error) {
	bottle, err := s.GetBottle(bottleId)
	if err != nil {
		return models.ChartModel{}, err
	}

	// Sort operations by date for the line chart
	ops := append([]models.BottleOperation(nil), bottle.OperationHistory...)
	sort.SliceStable(ops, func(i, j int) bool {
		return ops[i].Date < ops[j].Date
	})

	emptyWeight := bottle.InitialWeight - bottle.FillingWeight
	labels := []string{bottle.PurchaseDate}
	weights := []float64{bottle.InitialWeight}
	emptyWeights := []float64{emptyWeight}

	for _, op := range ops {
		labels = append(labels, op.Date)
		weights = append(weights, op.Weight)
		emptyWeights = append(emptyWeights, emptyWeight)
	}

	lang := configs.GetLanguage()
	return models.ChartModel{
		Type:   "line",
		Labels: labels,
		Sets: []models.ChartDataset{
			{
				Label: language.T(lang, "weight"),
				Unit:  "kg",
				Data:  weights,
			},
			{
				Label: language.T(lang, "emptyWeight"),
				Unit:  "kg",
				Data:  emptyWeights,
			},
		},
	}, nil
}
