package storage

import (
	"slices"
	"sort"

	"github.com/lukas-arnold/garden-equipment-log/internal/models"
	"github.com/lukas-arnold/garden-equipment-log/internal/utils"
)

func (s *Storage) AddBottle(input models.BottleInput) error {
	storage, err := s.getEquipmentStorage()
	if err != nil {
		return err
	}

	bottle := models.Bottle{
		Id:               utils.Id(),
		BottleInput:      input,
		OperationHistory: []models.BottleOperation{},
	}

	storage.Bottles = append(storage.Bottles, bottle)
	return s.saveStorage(storage)
}

func (s *Storage) GetBottles() ([]models.Bottle, error) {
	storage, err := s.getEquipmentStorage()
	if err != nil {
		return nil, err
	}

	for i := range storage.Bottles {
		enrichBottle(&storage.Bottles[i])
	}
	return storage.Bottles, nil
}

func (s *Storage) GetBottle(id int64) (models.Bottle, error) {
	bottles, err := s.GetBottles()
	if err != nil {
		return models.Bottle{}, err
	}

	for _, bottle := range bottles {
		if bottle.Id == id {
			return bottle, nil
		}
	}
	return models.Bottle{}, nil
}

func (s *Storage) GetBottleIdByOperationId(id int64) (int64, error) {
	storage, err := s.getEquipmentStorage()
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
	return 0, nil
}

func (s *Storage) UpdateBottle(bottle models.Bottle) error {
	storage, err := s.getEquipmentStorage()
	if err != nil {
		return err
	}

	for i := range storage.Bottles {
		if storage.Bottles[i].Id == bottle.Id {
			storage.Bottles[i] = bottle
			break
		}
	}
	return s.saveStorage(storage)
}

func (s *Storage) DeleteBottle(id int64) error {
	storage, err := s.getEquipmentStorage()
	if err != nil {
		return err
	}

	for i, b := range storage.Bottles {
		if b.Id == id {
			storage.Bottles = slices.Delete(storage.Bottles, i, i+1)
			return s.saveStorage(storage)
		}
	}
	return nil
}

func enrichBottle(bottle *models.Bottle) {
	bottle.TotalOperations = len(bottle.OperationHistory)
	usedGas := 0.0

	if len(bottle.OperationHistory) > 0 {
		lastWeight := bottle.OperationHistory[0].Weight
		usedGas = bottle.InitialWeight - lastWeight
	}

	bottle.UsedGas = usedGas
	bottle.RestGas = bottle.FillingWeight - usedGas
}

func sortBottles(bottles []models.Bottle) []models.Bottle {
	sort.Slice(bottles, func(i, j int) bool {
		if bottles[i].PurchaseDate != bottles[j].PurchaseDate {
			return bottles[i].PurchaseDate > bottles[j].PurchaseDate
		}
		return bottles[i].Id < bottles[j].Id
	})
	return bottles
}
