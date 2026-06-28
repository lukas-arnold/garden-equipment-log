package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/lukas-arnold/garden-equipment-log/internal/models"
)

func newTestStorage(t *testing.T) *Storage {
	t.Helper()

	return New(
		filepath.Join(
			t.TempDir(),
			"storage.json",
		),
	)
}

func TestStorageCreatesFile(t *testing.T) {

	store := newTestStorage(t)

	err := store.checkStorage()

	if err != nil {
		t.Fatal(err)
	}

	data, err := store.readStorage()

	if err != nil {
		t.Fatal(err)
	}

	if len(data) == 0 {
		t.Fatal("expected storage file")
	}
}

func TestStorageCreatesMissingDirectory(t *testing.T) {

	dir := t.TempDir()

	store := New(
		filepath.Join(
			dir,
			"nested",
			"storage.json",
		),
	)

	err := store.checkStorage()

	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(
		filepath.Join(
			dir,
			"nested",
		),
	)

	if err != nil {
		t.Fatal("expected directory to be created")
	}
}

func TestSaveStoragePersistsData(t *testing.T) {

	store := newTestStorage(t)

	input := models.EquipmentStorage{
		Devices: []models.Device{
			{
				DeviceInput: models.DeviceInput{
					Name: "Mower",
				},
			},
		},
	}

	err := store.saveStorage(input)

	if err != nil {
		t.Fatal(err)
	}

	result, err := store.getEquipmentStorage()

	if err != nil {
		t.Fatal(err)
	}

	if len(result.Devices) != 1 {
		t.Fatalf(
			"expected 1 device got %d",
			len(result.Devices),
		)
	}

	if result.Devices[0].Name != "Mower" {
		t.Fatalf(
			"wrong device %s",
			result.Devices[0].Name,
		)
	}
}

func TestSaveStorageSortsDevices(t *testing.T) {

	store := newTestStorage(t)

	err := store.saveStorage(
		models.EquipmentStorage{
			Devices: []models.Device{
				{
					DeviceInput: models.DeviceInput{
						Name:         "Zebra",
						PurchaseDate: "2026-01-01",
					},
				},
				{
					DeviceInput: models.DeviceInput{
						Name:         "Bench",
						PurchaseDate: "2026-02-01",
					},
				},
			},
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	result, err := store.getEquipmentStorage()

	if err != nil {
		t.Fatal(err)
	}

	if result.Devices[0].Name != "Bench" {
		t.Fatalf(
			"expected sorted device got %s",
			result.Devices[0].Name,
		)
	}
}

func TestSaveStorageSortsBottleOperations(t *testing.T) {

	store := newTestStorage(t)

	err := store.saveStorage(
		models.EquipmentStorage{
			Bottles: []models.Bottle{
				{
					OperationHistory: []models.BottleOperation{
						{
							BottleOperationInput: models.BottleOperationInput{
								Date: "2026-02-01",
							},
						},
						{
							BottleOperationInput: models.BottleOperationInput{
								Date: "2026-01-01",
							},
						},
					},
				},
			},
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	result, err := store.getEquipmentStorage()

	if err != nil {
		t.Fatal(err)
	}

	ops := result.Bottles[0].OperationHistory

	if ops[0].Date != "2026-02-01" {
		t.Fatalf(
			"expected sorted operations got %s",
			ops[0].Date,
		)
	}
}

func TestReadStorageCreatesEmptyStorage(t *testing.T) {

	store := newTestStorage(t)

	data, err := store.readStorage()

	if err != nil {
		t.Fatal(err)
	}

	if len(data) == 0 {
		t.Fatal(
			"expected initialized storage",
		)
	}
}

func TestGetEquipmentStorageEmpty(t *testing.T) {

	store := newTestStorage(t)

	result, err := store.getEquipmentStorage()

	if err != nil {
		t.Fatal(err)
	}

	if len(result.Devices) != 0 {
		t.Fatal(
			"expected no devices",
		)
	}

	if len(result.Bottles) != 0 {
		t.Fatal(
			"expected no bottles",
		)
	}
}
