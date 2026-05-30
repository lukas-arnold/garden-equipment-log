package storage

import (
	"fmt"
	"sort"
	"time"

	"github.com/lukas-arnold/garden-equipment-log/internal/models"
)

func GetDevicesForChart() (models.DevicesForChart, error) {
	devices, err := GetDevices()
	if err != nil {
		return models.DevicesForChart{}, err
	}
	yearSet := map[int]struct{}{}
	deviceYearHours := make(map[int64]map[int]float64)
	for _, device := range devices {
		yearHours := map[int]float64{}
		for _, operation := range device.OperationHistory {
			if operation.StartTime == "" || operation.EndTime == "" {
				continue
			}
			start, err := parseDateTime(operation.StartTime)
			if err != nil {
				continue
			}
			end, err := parseDateTime(operation.EndTime)
			if err != nil {
				continue
			}
			duration := end.Sub(start).Hours()
			if duration < 0 {
				duration = 0
			}
			year := start.Year()
			yearSet[year] = struct{}{}
			yearHours[year] += duration
		}
		deviceYearHours[device.Id] = yearHours
	}
	var years []int
	for year := range yearSet {
		years = append(years, year)
	}
	sort.Ints(years)
	var labels []string
	for _, year := range years {
		labels = append(labels, fmt.Sprintf("%d", year))
	}
	var hours [][]float64
	for _, device := range devices {
		row := make([]float64, len(years))
		for i, year := range years {
			row[i] = deviceYearHours[device.Id][year]
		}
		hours = append(hours, row)
	}
	return models.DevicesForChart{Devices: devices, Labels: labels, Hours: hours}, nil
}

func parseDateTime(value string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04", value)
}

func GetBottlesForChart() (models.BottlesForChart, error) {
	bottles, err := GetBottles()
	if err != nil {
		return models.BottlesForChart{}, err
	}
	var names []string
	var weights []float64
	for _, bottle := range bottles {
		names = append(names, fmt.Sprintf("Bottle %d", bottle.Id))
		if len(bottle.OperationHistory) > 0 {
			operations := append([]models.BottleOperation(nil), bottle.OperationHistory...)
			sort.Slice(operations, func(i, j int) bool {
				return operations[i].Date > operations[j].Date
			})
			weights = append(weights, operations[0].Weight)
		} else {
			weights = append(weights, bottle.InitialWeight)
		}
	}
	return models.BottlesForChart{Bottles: bottles, Names: names, LastWeights: weights}, nil
}

func GetDeviceOperationTimesForChart(deviceId int64, year int) (models.DeviceOperationTimesForChart, error) {
	device, err := GetDevice(deviceId)
	if err != nil {
		return models.DeviceOperationTimesForChart{}, err
	}
	// Build per-operation chart data (chronological)
	ops := append([]models.DeviceOperation(nil), device.OperationHistory...)
	// filter operations with valid start/end and optional year
	var filtered []models.DeviceOperation
	for _, op := range ops {
		if op.StartTime == "" || op.EndTime == "" {
			continue
		}
		if year != 0 {
			t, err := parseDateTime(op.StartTime)
			if err != nil {
				continue
			}
			if t.Year() != year {
				continue
			}
		}
		filtered = append(filtered, op)
	}
	// sort ascending by start time for chart (oldest -> newest)
	sort.SliceStable(filtered, func(i, j int) bool {
		return filtered[i].StartTime < filtered[j].StartTime
	})

	var dateLabels []string
	var durations []float64
	for _, op := range filtered {
		start, err := parseDateTime(op.StartTime)
		if err != nil {
			continue
		}
		end, err := parseDateTime(op.EndTime)
		if err != nil {
			continue
		}
		delta := end.Sub(start).Minutes()
		if delta < 0 {
			delta = 0
		}
		dateLabels = append(dateLabels, op.StartTime)
		durations = append(durations, delta)
	}

	// keep device.OperationHistory as-is for table display (already sorted by storage)
	return models.DeviceOperationTimesForChart{Device: device, Dates: dateLabels, Durations: durations}, nil
}

// GetDeviceOperationMinutesPerYearForChart aggregates operation durations (in minutes)
// per year for the given device and returns the result as a DeviceOperationTimesForChart
// where Dates are year strings and Durations are total minutes for that year.
func GetDeviceOperationMinutesPerYearForChart(deviceId int64) (models.DeviceOperationTimesForChart, error) {
	device, err := GetDevice(deviceId)
	if err != nil {
		return models.DeviceOperationTimesForChart{}, err
	}
	// aggregate minutes per year
	yearTotals := map[int]float64{}
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
		delta := end.Sub(start).Minutes()
		if delta < 0 {
			delta = 0
		}
		year := start.Year()
		yearTotals[year] += delta
	}
	var years []int
	for y := range yearTotals {
		years = append(years, y)
	}
	sort.Ints(years)
	var labels []string
	var durations []float64
	for _, y := range years {
		labels = append(labels, fmt.Sprintf("%d", y))
		durations = append(durations, yearTotals[y])
	}
	return models.DeviceOperationTimesForChart{Device: device, Dates: labels, Durations: durations}, nil
}

func GetBottleWeightHistoryForChart(bottleId int64) (models.BottleWeightHistoryForChart, error) {
	bottle, err := GetBottle(bottleId)
	if err != nil {
		return models.BottleWeightHistoryForChart{}, err
	}
	opHistory := append([]models.BottleOperation(nil), bottle.OperationHistory...)
	sort.SliceStable(opHistory, func(i, j int) bool {
		return opHistory[i].Date < opHistory[j].Date
	})
	var dates []string
	var weights []float64
	for _, operation := range opHistory {
		dates = append(dates, operation.Date)
		weights = append(weights, operation.Weight)
	}
	bottle.OperationHistory = opHistory
	return models.BottleWeightHistoryForChart{Bottle: bottle, Dates: dates, Weights: weights}, nil
}
