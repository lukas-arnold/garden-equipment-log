package storage

import (
	"slices"
	"sort"

	"github.com/lukas-arnold/garden-equipment-log/internal/models"
	"github.com/lukas-arnold/garden-equipment-log/internal/utils"
)

func (s *Storage) AddBottleOperation(bottleId int64, input models.BottleOperationInput) error {
	bottle, err := s.GetBottle(bottleId)
	if err != nil {
		return err
	}

	operation := models.BottleOperation{
		Id:                   utils.Id(),
		BottleOperationInput: input,
	}

	bottle.OperationHistory = append(bottle.OperationHistory, operation)
	return s.UpdateBottle(bottle)
}

func (s *Storage) GetBottleOperation(id int64) (models.BottleOperation, error) {
	bottles, err := s.GetBottles()
	if err != nil {
		return models.BottleOperation{}, err
	}

	for _, bottle := range bottles {
		for _, op := range bottle.OperationHistory {
			if op.Id == id {
				return op, nil
			}
		}
	}
	return models.BottleOperation{}, nil
}

func (s *Storage) UpdateBottleOperation(id int64, input models.BottleOperationInput) error {
	storage, err := s.getEquipmentStorage()
	if err != nil {
		return err
	}

	for i := range storage.Bottles {
		for j := range storage.Bottles[i].OperationHistory {
			if storage.Bottles[i].OperationHistory[j].Id == id {
				storage.Bottles[i].OperationHistory[j].Date = input.Date
				storage.Bottles[i].OperationHistory[j].Weight = input.Weight
				return s.saveStorage(storage)
			}
		}
	}
	return nil
}

func (s *Storage) DeleteBottleOperation(id int64) error {
	storage, err := s.getEquipmentStorage()
	if err != nil {
		return err
	}

	for i := range storage.Bottles {
		for j, op := range storage.Bottles[i].OperationHistory {
			if op.Id == id {
				storage.Bottles[i].OperationHistory = slices.Delete(
					storage.Bottles[i].OperationHistory,
					j,
					j+1,
				)
				return s.saveStorage(storage)
			}
		}
	}
	return nil
}

func sortBottleOperations(ops []models.BottleOperation) []models.BottleOperation {
	sort.Slice(ops, func(i, j int) bool {
		return ops[i].Date > ops[j].Date
	})
	return ops
}
