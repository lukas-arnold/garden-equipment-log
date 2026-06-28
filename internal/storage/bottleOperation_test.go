package storage

import (
	"path/filepath"
	"testing"

	"github.com/lukas-arnold/garden-equipment-log/internal/models"
)

func TestAddBottleOperation(t *testing.T) {

	store := New(
		filepath.Join(
			t.TempDir(),
			"storage.json",
		),
	)

	store.AddBottle(
		models.BottleInput{
			PurchaseDate: "2024-01-01",
		},
	)

	bottles, _ := store.GetBottles()

	err := store.AddBottleOperation(
		bottles[0].Id,
		models.BottleOperationInput{
			Date:   "2024-01-02",
			Weight: 14,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	bottle, _ := store.GetBottle(
		bottles[0].Id,
	)

	if len(bottle.OperationHistory) != 1 {
		t.Fatal("operation missing")
	}
}

func TestGetBottleOperation(t *testing.T) {

	store := New(
		filepath.Join(
			t.TempDir(),
			"storage.json",
		),
	)

	store.AddBottle(
		models.BottleInput{
			PurchaseDate: "2024-01-01",
		},
	)

	bottles, _ := store.GetBottles()

	store.AddBottleOperation(
		bottles[0].Id,
		models.BottleOperationInput{
			Date:   "2024-01-02",
			Weight: 14,
		},
	)

	bottle, _ := store.GetBottle(
		bottles[0].Id,
	)

	operation, err := store.GetBottleOperation(
		bottle.OperationHistory[0].Id,
	)

	if err != nil {
		t.Fatal(err)
	}

	if operation.Weight != 14 {
		t.Fatal("wrong operation")
	}
}

func TestUpdateBottleOperation(t *testing.T) {

	store := New(
		filepath.Join(
			t.TempDir(),
			"storage.json",
		),
	)

	store.AddBottle(
		models.BottleInput{
			PurchaseDate: "2024-01-01",
		},
	)

	bottles, _ := store.GetBottles()

	store.AddBottleOperation(
		bottles[0].Id,
		models.BottleOperationInput{
			Date:   "2024-01-02",
			Weight: 14,
		},
	)

	bottle, _ := store.GetBottle(
		bottles[0].Id,
	)

	err := store.UpdateBottleOperation(
		bottle.OperationHistory[0].Id,
		models.BottleOperationInput{
			Date:   "2024-01-03",
			Weight: 13,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	operation, _ := store.GetBottleOperation(
		bottle.OperationHistory[0].Id,
	)

	if operation.Date != "2024-01-03" {
		t.Fatal("date not updated")
	}

	if operation.Weight != 13 {
		t.Fatal("weight not updated")
	}
}

func TestDeleteBottleOperation(t *testing.T) {

	store := New(
		filepath.Join(
			t.TempDir(),
			"storage.json",
		),
	)

	store.AddBottle(
		models.BottleInput{
			PurchaseDate: "2024-01-01",
		},
	)

	bottles, _ := store.GetBottles()

	store.AddBottleOperation(
		bottles[0].Id,
		models.BottleOperationInput{
			Date:   "2024-01-02",
			Weight: 14,
		},
	)

	bottle, _ := store.GetBottle(
		bottles[0].Id,
	)

	err := store.DeleteBottleOperation(
		bottle.OperationHistory[0].Id,
	)

	if err != nil {
		t.Fatal(err)
	}

	updated, _ := store.GetBottle(
		bottles[0].Id,
	)

	if len(updated.OperationHistory) != 0 {
		t.Fatal("operation not deleted")
	}
}

func TestSortBottleOperations(t *testing.T) {

	operations := []models.BottleOperation{
		{
			BottleOperationInput: models.BottleOperationInput{
				Date: "2024-01-01",
			},
		},
		{
			BottleOperationInput: models.BottleOperationInput{
				Date: "2024-03-01",
			},
		},
		{
			BottleOperationInput: models.BottleOperationInput{
				Date: "2024-02-01",
			},
		},
	}

	sorted := sortBottleOperations(
		operations,
	)

	if sorted[0].Date != "2024-03-01" {
		t.Fatal("not sorted")
	}

	if sorted[1].Date != "2024-02-01" {
		t.Fatal("not sorted")
	}

	if sorted[2].Date != "2024-01-01" {
		t.Fatal("not sorted")
	}
}
