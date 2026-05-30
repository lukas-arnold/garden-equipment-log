package storage

import (
	"math"
	"os"
	"sort"

	"github.com/lukas-arnold/garden-equipment-log/internal/configs"
	"github.com/lukas-arnold/garden-equipment-log/internal/models"
	"github.com/lukas-arnold/garden-equipment-log/internal/utils"
)

func saveStorage(storage models.EquipmentStorage) error {
	storage = sortStorage(storage)
	bytes, err := utils.ConvertEquipmentStorageToBytes(storage)
	if err != nil {
		return err
	}
	err = os.WriteFile(configs.GetStorageFile(), []byte(bytes), 0666)
	if err != nil {
		return err
	}
	return nil
}

func readStorage() ([]byte, error) {
	checkStorage()
	bytes, err := os.ReadFile(configs.GetStorageFile())
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func checkStorage() {
	_, err := os.ReadFile(configs.GetStorageFile())
	if err != nil {
		saveStorage(models.EquipmentStorage{})
	}
}

func GetEquipmentStorage() (models.EquipmentStorage, error) {
	bytes, err := readStorage()
	if err != nil {
		return models.EquipmentStorage{}, err
	}
	storage, err := utils.ConvertBytesToEquipmentStorage(bytes)
	if err != nil {
		return models.EquipmentStorage{}, err
	}
	storage = sortStorage(storage)
	for i := range storage.Bottles {
		storage.Bottles[i].CurrentFillLevel = calculateBottleFillLevel(storage.Bottles[i])
	}
	for i := range storage.Devices {
		storage.Devices[i].LastUsageDate = latestDeviceOperationTimestamp(storage.Devices[i])
		storage.Devices[i].TotalOperations = calculateDeviceTotalOperations(storage.Devices[i])
		storage.Devices[i].TotalOperationHours = calculateDeviceTotalOperationHours(storage.Devices[i])
		storage.Devices[i].PricePerHour = calculateDevicePricePerHour(storage.Devices[i])
	}
	return storage, nil
}

func calculateDeviceTotalOperations(device models.Device) int {
	return len(device.OperationHistory)
}

func calculateDeviceTotalOperationHours(device models.Device) float64 {
	totalHours := 0.0
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
		if duration > 0 {
			totalHours += duration
		}
	}
	return totalHours
}

func calculateDevicePricePerHour(device models.Device) float64 {
	totalHours := calculateDeviceTotalOperationHours(device)
	if totalHours <= 0 {
		return 0
	}
	return device.PurchasePrice / totalHours
}

func calculateBottleFillLevel(bottle models.Bottle) float64 {
	currentWeight := bottle.InitialWeight
	if len(bottle.OperationHistory) > 0 {
		currentWeight = getLatestBottleWeight(bottle)
	}

	minW := math.Min(bottle.InitialWeight, bottle.FillingWeight)
	maxW := math.Max(bottle.InitialWeight, bottle.FillingWeight)
	if maxW == minW {
		return 0
	}
	capacity := maxW - minW

	// If operation weights are recorded as content volume (not total mass), detect
	// and interpret latest weight as content when it's <= capacity.
	if capacity > 0 && len(bottle.OperationHistory) > 0 {
		latest := getLatestBottleWeight(bottle)
		if latest <= capacity {
			// latest is content weight
			level := (latest / capacity) * 100
			if level < 0 {
				return 0
			}
			if level > 100 {
				return 100
			}
			return level
		}
	}

	if currentWeight < minW {
		currentWeight = minW
	}
	if currentWeight > maxW {
		currentWeight = maxW
	}
	level := (currentWeight - minW) / (maxW - minW) * 100
	if level < 0 {
		return 0
	}
	if level > 100 {
		return 100
	}
	return level
}

func getLatestBottleWeight(bottle models.Bottle) float64 {
	currentWeight := bottle.InitialWeight
	if len(bottle.OperationHistory) == 0 {
		return currentWeight
	}
	operations := append([]models.BottleOperation(nil), bottle.OperationHistory...)
	sort.Slice(operations, func(i, j int) bool {
		return operations[i].Date > operations[j].Date
	})
	return operations[0].Weight
}

func sortStorage(storage models.EquipmentStorage) models.EquipmentStorage {
	storage.Devices = sortDevices(storage.Devices)
	for i := range storage.Devices {
		storage.Devices[i].OperationHistory = sortDeviceOperations(storage.Devices[i].OperationHistory)
	}
	storage.Bottles = sortBottles(storage.Bottles)
	for i := range storage.Bottles {
		storage.Bottles[i].OperationHistory = sortBottleOperations(storage.Bottles[i].OperationHistory)
	}
	return storage
}
