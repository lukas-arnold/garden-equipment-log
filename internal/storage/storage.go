package storage

import (
	"os"
	"path/filepath"

	"github.com/lukas-arnold/garden-equipment-log/internal/models"
	"github.com/lukas-arnold/garden-equipment-log/internal/utils"
)

type Storage struct {
	file string
}

func New(file string) *Storage {
	return &Storage{file: file}
}

func (s *Storage) saveStorage(storage models.EquipmentStorage) error {
	storage = sortStorage(storage)

	data, err := utils.ConvertEquipmentStorageToBytes(storage)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(s.file), 0755); err != nil {
		return err
	}

	return os.WriteFile(s.file, data, 0644)
}

func (s *Storage) readStorage() ([]byte, error) {
	if err := s.checkStorage(); err != nil {
		return nil, err
	}

	return os.ReadFile(s.file)
}

func (s *Storage) checkStorage() error {
	_, err := os.ReadFile(s.file)
	if err == nil {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(s.file), 0755); err != nil {
		return err
	}

	return s.saveStorage(models.EquipmentStorage{})
}

func (s *Storage) getEquipmentStorage() (models.EquipmentStorage, error) {
	bytes, err := s.readStorage()
	if err != nil {
		return models.EquipmentStorage{}, err
	}

	storage, err := utils.ConvertBytesToEquipmentStorage(bytes)
	if err != nil {
		return models.EquipmentStorage{}, err
	}

	return sortStorage(storage), nil
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
