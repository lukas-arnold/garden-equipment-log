package storage

import (
	"slices"
	"sort"
	"time"

	"github.com/lukas-arnold/garden-equipment-log/internal/models"
	"github.com/lukas-arnold/garden-equipment-log/internal/utils"
)

func AddDevice(device models.DeviceInput) error {
	storage, err := GetEquipmentStorage()
	if err != nil {
		return err
	}
	newDevice := models.Device{Id: utils.Id(), DeviceInput: device, OperationHistory: []models.DeviceOperation{}}
	storage.Devices = append(storage.Devices, newDevice)
	err = saveStorage(storage)
	if err != nil {
		return err
	}
	return nil
}

func GetDevices() ([]models.Device, error) {
	storage, err := GetEquipmentStorage()
	if err != nil {
		return nil, err
	}
	for i := range storage.Devices {
		enrichDevice(&storage.Devices[i])
	}
	return storage.Devices, nil
}

func GetDevice(id int64) (models.Device, error) {
	devices, err := GetDevices()
	if err != nil {
		return models.Device{}, err
	}
	var device models.Device
	for _, value := range devices {
		if value.Id == id {
			device = value
		}
	}
	return device, nil
}

func UpdateDevice(device models.Device) error {
	storage, err := GetEquipmentStorage()
	if err != nil {
		return err
	}
	for i := range storage.Devices {
		if storage.Devices[i].Id == device.Id {
			storage.Devices[i].Name = device.Name
			storage.Devices[i].PurchaseDate = device.PurchaseDate
			storage.Devices[i].PurchasePrice = device.PurchasePrice
			storage.Devices[i].OperationHistory = device.OperationHistory
		}
	}
	err = saveStorage(storage)
	if err != nil {
		return err
	}
	return nil
}

func DeleteDevice(id int64) error {
	storage, err := GetEquipmentStorage()
	if err != nil {
		return err
	}
	var index int
	for i := range storage.Devices {
		if storage.Devices[i].Id == id {
			index = i
		}
	}
	storage.Devices = slices.Delete(storage.Devices, index, index+1)
	err = saveStorage(storage)
	if err != nil {
		return err
	}
	return nil
}

func enrichDevice(device *models.Device) {
	device.TotalOperations = len(device.OperationHistory)
	var totalHours float64
	var lastUsage time.Time
	for _, op := range device.OperationHistory {
		if op.StartTime == "" || op.EndTime == "" {
			continue
		}
		start, err := time.Parse("2006-01-02T15:04", op.StartTime)
		if err != nil {
			continue
		}
		end, err := time.Parse("2006-01-02T15:04", op.EndTime)
		if err != nil {
			continue
		}
		hours := end.Sub(start).Hours()
		if hours < 0 {
			hours = 0
		}
		totalHours += hours
		if start.After(lastUsage) {
			lastUsage = start
		}
	}
	device.TotalOperationHours = totalHours
	if totalHours > 0 {
		device.PricePerHour = device.PurchasePrice / totalHours
	}
	if !lastUsage.IsZero() {
		device.LastUsageDate = lastUsage.Format("2006-01-02")
	}
}

func sortDevices(devices []models.Device) []models.Device {
	sort.Slice(devices, func(i, j int) bool {
		iLast := latestDeviceOperationTimestamp(devices[i])
		jLast := latestDeviceOperationTimestamp(devices[j])
		if iLast != jLast {
			return iLast > jLast
		}
		if devices[i].PurchaseDate != devices[j].PurchaseDate {
			return devices[i].PurchaseDate > devices[j].PurchaseDate
		}
		return devices[i].Id < devices[j].Id
	})
	return devices
}

func latestDeviceOperationTimestamp(device models.Device) string {
	if len(device.OperationHistory) == 0 {
		return device.PurchaseDate
	}
	lastOp := device.OperationHistory[0]
	if lastOp.EndTime != "" {
		return lastOp.EndTime
	}
	return lastOp.StartTime
}
