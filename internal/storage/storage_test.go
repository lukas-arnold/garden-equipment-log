package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/lukas-arnold/garden-equipment-log/internal/models"
)

func testStorage(t *testing.T) *Storage {
	t.Helper()

	return New(filepath.Join(t.TempDir(), "storage.json"))
}

func TestStorageCreatesFile(t *testing.T) {
	store := testStorage(t)
	_, _ = store.readStorage()

	if _, err := os.Stat(store.file); os.IsNotExist(err) {
		t.Fatal("expected storage file to be created")
	}
}

func TestStorageCreatesMissingDirectory(t *testing.T) {
	dir := t.TempDir()
	store := New(filepath.Join(dir, "nested", "storage.json"))

	err := store.checkStorage()
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(filepath.Join(dir, "nested"))
	if err != nil {
		t.Fatal("expected directory to be created")
	}
}

func TestSaveStoragePersistsData(t *testing.T) {
	store := testStorage(t)
	input := models.EquipmentStorage{
		Devices: []models.Device{{DeviceInput: models.DeviceInput{Name: "Mower"}}},
	}

	if err := store.saveStorage(input); err != nil {
		t.Fatal(err)
	}

	result, err := store.getEquipmentStorage()
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Devices) != 1 || result.Devices[0].Name != "Mower" {
		t.Fatalf("data mismatch: expected 1 device named Mower, got %v", result.Devices)
	}
}

func TestSaveStorageSortsDevices(t *testing.T) {
	store := testStorage(t)
	input := models.EquipmentStorage{
		Devices: []models.Device{
			{DeviceInput: models.DeviceInput{Name: "Zebra", PurchaseDate: "2026-01-01"}},
			{DeviceInput: models.DeviceInput{Name: "Bench", PurchaseDate: "2026-02-01"}},
		},
	}

	if err := store.saveStorage(input); err != nil {
		t.Fatal(err)
	}

	result, _ := store.getEquipmentStorage()
	if result.Devices[0].Name != "Bench" {
		t.Fatalf("expected sorted device (Bench), got %s", result.Devices[0].Name)
	}
}

func TestSaveStorageSortsBottleOperations(t *testing.T) {
	store := testStorage(t)
	input := models.EquipmentStorage{
		Bottles: []models.Bottle{
			{
				OperationHistory: []models.BottleOperation{
					{BottleOperationInput: models.BottleOperationInput{Date: "2026-02-01"}},
					{BottleOperationInput: models.BottleOperationInput{Date: "2026-01-01"}},
				},
			},
		},
	}

	if err := store.saveStorage(input); err != nil {
		t.Fatal(err)
	}

	result, _ := store.getEquipmentStorage()
	ops := result.Bottles[0].OperationHistory
	if ops[0].Date != "2026-02-01" {
		t.Fatalf("expected sorted operations, got first date: %s", ops[0].Date)
	}
}

func TestReadStorageCreatesEmptyStorage(t *testing.T) {
	store := testStorage(t)
	data, err := store.readStorage()
	if err != nil || len(data) == 0 {
		t.Fatal("expected initialized empty storage file")
	}
}

func TestGetEquipmentStorageEmpty(t *testing.T) {
	store := testStorage(t)
	result, err := store.getEquipmentStorage()
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Devices) != 0 || len(result.Bottles) != 0 {
		t.Fatal("expected empty storage structure")
	}
}
