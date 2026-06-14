package storage

import (
	"slices"
	"sort"

	"github.com/lukas-arnold/garden-equipment-log/internal/models"
	"github.com/lukas-arnold/garden-equipment-log/internal/utils"
)

func AddBottleOperation(bottleId int64, operationInput models.BottleOperationInput) error {
	bottle, err := GetBottle(bottleId)
	if err != nil {
		return err
	}
	operation := models.BottleOperation{Id: utils.Id(), BottleOperationInput: operationInput}
	bottle.OperationHistory = append(bottle.OperationHistory, operation)
	err = UpdateBottle(bottle)
	if err != nil {
		return err
	}
	return nil
}

func GetBottleOperation(id int64) (models.BottleOperation, error) {
	bottles, err := GetBottles()
	if err != nil {
		return models.BottleOperation{}, err
	}
	var operation models.BottleOperation
	for _, bottle := range bottles {
		for _, value := range bottle.OperationHistory {
			if value.Id == id {
				operation = value
			}
		}
	}
	return operation, nil
}

func UpdateBottleOperation(id int64, operationInput models.BottleOperationInput) error {
	storage, err := getEquipmentStorage()
	if err != nil {
		return err
	}
	for i := range storage.Bottles {
		for j := range storage.Bottles[i].OperationHistory {
			if storage.Bottles[i].OperationHistory[j].Id == id {
				storage.Bottles[i].OperationHistory[j].Date = operationInput.Date
				storage.Bottles[i].OperationHistory[j].Weight = operationInput.Weight
				return saveStorage(storage)
			}
		}
	}
	return nil
}

func DeleteBottleOperation(id int64) error {
	storage, err := getEquipmentStorage()
	if err != nil {
		return err
	}
	for i := range storage.Bottles {
		for j := range storage.Bottles[i].OperationHistory {
			if storage.Bottles[i].OperationHistory[j].Id == id {
				storage.Bottles[i].OperationHistory = slices.Delete(storage.Bottles[i].OperationHistory, j, j+1)
				return saveStorage(storage)
			}
		}
	}
	return nil
}

func sortBottleOperations(operations []models.BottleOperation) []models.BottleOperation {
	sort.Slice(operations, func(i, j int) bool {
		return operations[i].Date > operations[j].Date
	})
	return operations
}
