package storage

import (
	"testing"

	"github.com/lukas-arnold/garden-equipment-log/internal/models"
)

func TestAddBottle(t *testing.T) {
	store := testStorage(t)

	err := store.AddBottle(models.BottleInput{PurchaseDate: "2024-01-01"})
	if err != nil {
		t.Fatal(err)
	}

	bottles, _ := store.GetBottles()
	if len(bottles) != 1 {
		t.Fatal("bottle missing")
	}
}

func TestGetBottles(t *testing.T) {
	store := testStorage(t)
	store.AddBottle(models.BottleInput{PurchaseDate: "2024-01-01"})

	bottles, err := store.GetBottles()
	if err != nil {
		t.Fatal(err)
	}
	if len(bottles) != 1 {
		t.Fatal("bottle missing")
	}
}

func TestGetBottle(t *testing.T) {
	store := testStorage(t)
	store.AddBottle(models.BottleInput{PurchaseDate: "2024-01-01"})

	bottles, _ := store.GetBottles()
	bottle, err := store.GetBottle(bottles[0].Id)
	if err != nil {
		t.Fatal(err)
	}
	if bottle.Id != bottles[0].Id {
		t.Fatal("wrong bottle")
	}
}

func TestGetBottleIdByOperationId(t *testing.T) {
	store := testStorage(t)
	store.AddBottle(models.BottleInput{
		PurchaseDate:  "2024-01-01",
		InitialWeight: 15,
		FillingWeight: 5,
	})

	bottles, _ := store.GetBottles()
	store.AddBottleOperation(bottles[0].Id, models.BottleOperationInput{
		Date:   "2024-01-02",
		Weight: 14,
	})

	bottle, _ := store.GetBottle(bottles[0].Id)
	id, err := store.GetBottleIdByOperationId(bottle.OperationHistory[0].Id)
	if err != nil {
		t.Fatal(err)
	}
	if id != bottle.Id {
		t.Fatal("wrong bottle id")
	}
}

func TestUpdateBottle(t *testing.T) {
	store := testStorage(t)
	store.AddBottle(models.BottleInput{PurchaseDate: "2024-01-01"})

	bottles, _ := store.GetBottles()
	bottle := bottles[0]
	bottle.PurchasePrice = 99

	if err := store.UpdateBottle(bottle); err != nil {
		t.Fatal(err)
	}

	updated, _ := store.GetBottle(bottle.Id)
	if updated.PurchasePrice != 99 {
		t.Fatal("bottle not updated")
	}
}

func TestDeleteBottle(t *testing.T) {
	store := testStorage(t)
	store.AddBottle(models.BottleInput{PurchaseDate: "2024-01-01"})

	bottles, _ := store.GetBottles()
	if err := store.DeleteBottle(bottles[0].Id); err != nil {
		t.Fatal(err)
	}

	result, _ := store.GetBottles()
	if len(result) != 0 {
		t.Fatal("bottle not deleted")
	}
}

func TestEnrichBottle(t *testing.T) {
	bottle := models.Bottle{
		BottleInput: models.BottleInput{InitialWeight: 15, FillingWeight: 5},
		OperationHistory: []models.BottleOperation{
			{BottleOperationInput: models.BottleOperationInput{Weight: 13}},
		},
	}

	enrichBottle(&bottle)

	if bottle.TotalOperations != 1 {
		t.Fatal("wrong operations")
	}
	if bottle.UsedGas != 2 {
		t.Fatalf("got %f", bottle.UsedGas)
	}
	if bottle.RestGas != 3 {
		t.Fatalf("got %f", bottle.RestGas)
	}
}

func TestSortBottles(t *testing.T) {
	bottles := []models.Bottle{
		{Id: 2, BottleInput: models.BottleInput{PurchaseDate: "2024-01-01"}},
		{Id: 1, BottleInput: models.BottleInput{PurchaseDate: "2025-01-01"}},
	}

	sorted := sortBottles(bottles)
	if sorted[0].Id != 1 {
		t.Fatal("not sorted")
	}
}

func TestEnrichBottleEmpty(t *testing.T) {
	bottle := models.Bottle{
		BottleInput: models.BottleInput{InitialWeight: 15, FillingWeight: 5},
	}

	enrichBottle(&bottle)

	if bottle.TotalOperations != 0 {
		t.Fatal("wrong operations")
	}
	if bottle.UsedGas != 0 {
		t.Fatal("wrong used gas")
	}
	if bottle.RestGas != 5 {
		t.Fatal("wrong rest gas")
	}
}
