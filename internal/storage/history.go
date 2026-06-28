package storage

import (
	"fmt"
	"sort"

	"github.com/lukas-arnold/garden-equipment-log/internal/configs"
	"github.com/lukas-arnold/garden-equipment-log/internal/language"
	"github.com/lukas-arnold/garden-equipment-log/internal/models"
	"github.com/lukas-arnold/garden-equipment-log/internal/utils"
)

func (s *Storage) GetDeviceChart(
	deviceId int64,
) (models.ChartModel, error) {
	device, err := s.GetDevice(deviceId)

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

	years := make([]int, 0)

	for year := range yearMinutes {
		years = append(years, year)
	}

	sort.Ints(years)

	labels := make([]string, 0)
	values := make([]float64, 0)

	for _, year := range years {
		labels = append(
			labels,
			fmt.Sprintf("%d", year),
		)

		values = append(
			values,
			yearMinutes[year],
		)
	}

	return models.ChartModel{
		Type:   "bar",
		Labels: labels,
		Sets: []models.ChartDataset{
			{
				Label: language.T(
					configs.GetLanguage(),
					"operationTime",
				),
				Unit: "h/min",
				Data: values,
			},
		},
	}, nil
}

func (s *Storage) GetBottleChart(
	bottleId int64,
) (models.ChartModel, error) {
	bottle, err := s.GetBottle(bottleId)

	if err != nil {
		return models.ChartModel{}, err
	}

	ops := append(
		[]models.BottleOperation(nil),
		bottle.OperationHistory...,
	)

	sort.SliceStable(ops, func(i, j int) bool {
		return ops[i].Date < ops[j].Date
	})

	emptyWeight := bottle.InitialWeight - bottle.FillingWeight

	labels := []string{
		bottle.PurchaseDate,
	}

	weights := []float64{
		bottle.InitialWeight,
	}

	emptyWeights := []float64{
		emptyWeight,
	}

	for _, op := range ops {
		labels = append(
			labels,
			op.Date,
		)

		weights = append(
			weights,
			op.Weight,
		)

		emptyWeights = append(
			emptyWeights,
			emptyWeight,
		)
	}

	return models.ChartModel{
		Type:   "line",
		Labels: labels,
		Sets: []models.ChartDataset{
			{
				Label: language.T(
					configs.GetLanguage(),
					"weight",
				),
				Unit: "kg",
				Data: weights,
			},
			{
				Label: language.T(
					configs.GetLanguage(),
					"emptyWeight",
				),
				Unit: "kg",
				Data: emptyWeights,
			},
		},
	}, nil
}
