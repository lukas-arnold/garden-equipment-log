package storage

import (
	"slices"
	"sort"

	"github.com/lukas-arnold/garden-equipment-log/internal/models"
	"github.com/lukas-arnold/garden-equipment-log/internal/utils"
)

func AddBottle(bottle models.BottleInput) error {
	storage, err := GetEquipmentStorage()
	if err != nil {
		return err
	}
	newBottle := models.Bottle{Id: utils.Id(), BottleInput: bottle, OperationHistory: []models.BottleOperation{}}
	storage.Bottles = append(storage.Bottles, newBottle)
	err = saveStorage(storage)
	if err != nil {
		return err
	}
	return nil
}

func GetBottles() ([]models.Bottle, error) {
	storage, err := GetEquipmentStorage()
	if err != nil {
		return nil, err
	}
	return storage.Bottles, nil
}

func GetBottle(id int64) (models.Bottle, error) {
	bottles, err := GetBottles()
	if err != nil {
		return models.Bottle{}, err
	}
	var bottle models.Bottle
	for _, value := range bottles {
		if value.Id == id {
			bottle = value
		}
	}
	return bottle, nil
}

func UpdateBottle(bottle models.Bottle) error {
	storage, err := GetEquipmentStorage()
	if err != nil {
		return err
	}
	for i := range storage.Bottles {
		if storage.Bottles[i].Id == bottle.Id {
			storage.Bottles[i].PurchaseDate = bottle.PurchaseDate
			storage.Bottles[i].PurchasePrice = bottle.PurchasePrice
			storage.Bottles[i].InitialWeight = bottle.InitialWeight
			storage.Bottles[i].FillingWeight = bottle.FillingWeight
			storage.Bottles[i].OperationHistory = bottle.OperationHistory
		}
	}
	err = saveStorage(storage)
	if err != nil {
		return err
	}
	return nil
}

func DeleteBottle(id int64) error {
	storage, err := GetEquipmentStorage()
	if err != nil {
		return err
	}
	var index int
	for i := range storage.Bottles {
		if storage.Bottles[i].Id == id {
			index = i
		}
	}
	storage.Bottles = slices.Delete(storage.Bottles, index, index+1)
	err = saveStorage(storage)
	if err != nil {
		return err
	}
	return nil
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
