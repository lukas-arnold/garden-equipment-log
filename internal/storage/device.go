package storage

import (
	"slices"
	"sort"
	"time"

	"github.com/lukas-arnold/garden-equipment-log/internal/models"
	"github.com/lukas-arnold/garden-equipment-log/internal/utils"
)

func (s *Storage) AddDevice(input models.DeviceInput) error {
	storage, err := s.getEquipmentStorage()
	if err != nil {
		return err
	}

	device := models.Device{
		Id:               utils.Id(),
		DeviceInput:      input,
		OperationHistory: []models.DeviceOperation{},
	}

	storage.Devices = append(storage.Devices, device)
	return s.saveStorage(storage)
}

func (s *Storage) GetDevices() ([]models.Device, error) {
	storage, err := s.getEquipmentStorage()
	if err != nil {
		return nil, err
	}

	for i := range storage.Devices {
		enrichDevice(&storage.Devices[i])
	}

	return storage.Devices, nil
}

func (s *Storage) GetDevice(id int64) (models.Device, error) {
	storage, err := s.getEquipmentStorage()
	if err != nil {
		return models.Device{}, err
	}

	for _, device := range storage.Devices {
		if device.Id == id {
			enrichDevice(&device)
			return device, nil
		}
	}

	return models.Device{}, nil
}

func (s *Storage) GetDeviceIdByOperationId(id int64) (int64, error) {
	storage, err := s.getEquipmentStorage()
	if err != nil {
		return 0, err
	}

	for _, device := range storage.Devices {
		for _, operation := range device.OperationHistory {
			if operation.Id == id {
				return device.Id, nil
			}
		}
	}
	return 0, nil
}

func (s *Storage) UpdateDevice(device models.Device) error {
	storage, err := s.getEquipmentStorage()
	if err != nil {
		return err
	}

	for i := range storage.Devices {
		if storage.Devices[i].Id == device.Id {
			storage.Devices[i] = device
			break
		}
	}

	return s.saveStorage(storage)
}

func (s *Storage) DeleteDevice(id int64) error {
	storage, err := s.getEquipmentStorage()
	if err != nil {
		return err
	}

	idx := slices.IndexFunc(storage.Devices, func(d models.Device) bool { return d.Id == id })
	if idx != -1 {
		storage.Devices = slices.Delete(storage.Devices, idx, idx+1)
		return s.saveStorage(storage)
	}

	return nil
}

func enrichDevice(device *models.Device) {
	device.TotalOperations = len(device.OperationHistory)

	var totalMinutes float64
	var lastUsage time.Time

	for _, op := range device.OperationHistory {
		start, err1 := utils.ParseDateTime(op.StartTime)
		end, err2 := utils.ParseDateTime(op.EndTime)

		if err1 != nil || err2 != nil {
			continue
		}

		totalMinutes += max(0, end.Sub(start).Minutes())

		if start.After(lastUsage) {
			lastUsage = start
		}
	}

	device.TotalOperationTime = totalMinutes

	if totalMinutes > 0 {
		device.PricePerHour = device.PurchasePrice / (totalMinutes / 60)
	}

	if !lastUsage.IsZero() {
		device.LastUsageDate = lastUsage.Format(time.DateOnly)
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
