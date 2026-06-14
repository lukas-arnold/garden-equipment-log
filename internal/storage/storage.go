package storage

import (
	"os"

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

func getEquipmentStorage() (models.EquipmentStorage, error) {
	bytes, err := readStorage()
	if err != nil {
		return models.EquipmentStorage{}, err
	}
	storage, err := utils.ConvertBytesToEquipmentStorage(bytes)
	if err != nil {
		return models.EquipmentStorage{}, err
	}
	storage = sortStorage(storage)
	return storage, nil
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
