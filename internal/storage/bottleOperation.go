package storage

import (
	"fmt"
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

func DeleteBottleOperation(id int64) error {
	storage, err := GetEquipmentStorage()
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

func UpdateBottleOperation(id int64, operationInput models.BottleOperationInput) error {
	storage, err := GetEquipmentStorage()
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

func GetBottleIdByOperationId(id int64) (int64, error) {
	storage, err := GetEquipmentStorage()
	if err != nil {
		return 0, err
	}
	for _, bottle := range storage.Bottles {
		for _, operation := range bottle.OperationHistory {
			if operation.Id == id {
				return bottle.Id, nil
			}
		}
	}
	return 0, fmt.Errorf("bottle operation %d not found", id)
}

func sortBottleOperations(operations []models.BottleOperation) []models.BottleOperation {
	sort.Slice(operations, func(i, j int) bool {
		return operations[i].Date > operations[j].Date
	})
	return operations
}
