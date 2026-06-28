package storage

import (
	"slices"
	"sort"

	"github.com/lukas-arnold/garden-equipment-log/internal/models"
	"github.com/lukas-arnold/garden-equipment-log/internal/utils"
)

func (s *Storage) AddDeviceOperation(deviceId int64, input models.DeviceOperationInput) error {
	device, err := s.GetDevice(deviceId)
	if err != nil {
		return err
	}

	operation := models.DeviceOperation{
		Id:                   utils.Id(),
		DeviceOperationInput: input,
	}

	device.OperationHistory = append(device.OperationHistory, operation)
	return s.UpdateDevice(device)
}

func (s *Storage) GetDeviceOperation(id int64) (models.DeviceOperation, error) {
	devices, err := s.GetDevices()
	if err != nil {
		return models.DeviceOperation{}, err
	}

	for _, device := range devices {
		for _, op := range device.OperationHistory {
			if op.Id == id {
				return op, nil
			}
		}
	}
	return models.DeviceOperation{}, nil
}

func (s *Storage) UpdateDeviceOperation(id int64, input models.DeviceOperationInput) error {
	storage, err := s.getEquipmentStorage()
	if err != nil {
		return err
	}

	for i := range storage.Devices {
		for j := range storage.Devices[i].OperationHistory {
			if storage.Devices[i].OperationHistory[j].Id == id {
				storage.Devices[i].OperationHistory[j].DeviceOperationInput = input
				return s.saveStorage(storage)
			}
		}
	}
	return nil
}

func (s *Storage) DeleteDeviceOperation(id int64) error {
	storage, err := s.getEquipmentStorage()
	if err != nil {
		return err
	}

	for i := range storage.Devices {
		for j, op := range storage.Devices[i].OperationHistory {
			if op.Id == id {
				storage.Devices[i].OperationHistory = slices.Delete(
					storage.Devices[i].OperationHistory,
					j,
					j+1,
				)
				return s.saveStorage(storage)
			}
		}
	}
	return nil
}

func sortDeviceOperations(ops []models.DeviceOperation) []models.DeviceOperation {
	sort.Slice(ops, func(i, j int) bool {
		return ops[i].StartTime > ops[j].StartTime
	})
	return ops
}
