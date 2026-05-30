package utils

import (
	"encoding/json"
	"strconv"

	"github.com/lukas-arnold/garden-equipment-log/internal/models"
)

func ConvertEquipmentStorageToBytes(storage models.EquipmentStorage) ([]byte, error) {
	bytes, err := json.Marshal(storage)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func ConvertBytesToEquipmentStorage(bytes []byte) (models.EquipmentStorage, error) {
	var storage models.EquipmentStorage
	err := json.Unmarshal(bytes, &storage)
	if err != nil {
		return storage, err
	}
	return storage, nil
}

func ConvertId(idStr string) (int64, error) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func ConvertFloat(floatStr string) (float64, error) {
	value, err := strconv.ParseFloat(floatStr, 64)
	if err != nil {
		return -1, err
	}
	return value, nil
}
