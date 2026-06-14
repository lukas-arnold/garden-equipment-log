package storage

import (
	"slices"
	"sort"

	"github.com/lukas-arnold/garden-equipment-log/internal/models"
	"github.com/lukas-arnold/garden-equipment-log/internal/utils"
)

func AddDeviceOperation(deviceId int64, operationInput models.DeviceOperationInput) error {
	device, err := GetDevice(deviceId)
	if err != nil {
		return err
	}
	operation := models.DeviceOperation{Id: utils.Id(), DeviceOperationInput: operationInput}
	device.OperationHistory = append(device.OperationHistory, operation)
	err = UpdateDevice(device)
	if err != nil {
		return err
	}
	return nil
}

func GetDeviceOperation(id int64) (models.DeviceOperation, error) {
	devices, err := GetDevices()
	if err != nil {
		return models.DeviceOperation{}, err
	}
	var operation models.DeviceOperation
	for _, device := range devices {
		for _, value := range device.OperationHistory {
			if value.Id == id {
				operation = value
			}
		}
	}
	return operation, nil
}

func DeleteDeviceOperation(id int64) error {
	storage, err := getEquipmentStorage()
	if err != nil {
		return err
	}
	for i := range storage.Devices {
		for j := range storage.Devices[i].OperationHistory {
			if storage.Devices[i].OperationHistory[j].Id == id {
				storage.Devices[i].OperationHistory = slices.Delete(storage.Devices[i].OperationHistory, j, j+1)
				return saveStorage(storage)
			}
		}
	}
	return nil
}

func UpdateDeviceOperation(id int64, operationInput models.DeviceOperationInput) error {
	storage, err := getEquipmentStorage()
	if err != nil {
		return err
	}
	for i := range storage.Devices {
		for j := range storage.Devices[i].OperationHistory {
			if storage.Devices[i].OperationHistory[j].Id == id {
				storage.Devices[i].OperationHistory[j].StartTime = operationInput.StartTime
				storage.Devices[i].OperationHistory[j].EndTime = operationInput.EndTime
				storage.Devices[i].OperationHistory[j].Note = operationInput.Note
				return saveStorage(storage)
			}
		}
	}
	return nil
}

func sortDeviceOperations(operations []models.DeviceOperation) []models.DeviceOperation {
	sort.Slice(operations, func(i, j int) bool {
		return operations[i].StartTime > operations[j].StartTime
	})
	return operations
}
