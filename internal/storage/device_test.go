package storage

import (
	"testing"

	"github.com/lukas-arnold/garden-equipment-log/internal/models"
)

func TestAddDevice(t *testing.T) {
	store := testStorage(t)
	err := store.AddDevice(models.DeviceInput{Name: "Mower"})
	if err != nil {
		t.Fatal(err)
	}

	devices, _ := store.GetDevices()
	if len(devices) != 1 || devices[0].Name != "Mower" {
		t.Fatal("device was not added or data mismatch")
	}
}

func TestGetDevices(t *testing.T) {
	store := testStorage(t)
	store.AddDevice(models.DeviceInput{Name: "Trimmer"})

	devices, err := store.GetDevices()
	if err != nil {
		t.Fatal(err)
	}
	if len(devices) != 1 {
		t.Fatal("expected 1 device, got", len(devices))
	}
}

func TestGetDevice(t *testing.T) {
	store := testStorage(t)
	store.AddDevice(models.DeviceInput{Name: "Chainsaw"})
	devices, _ := store.GetDevices()

	device, err := store.GetDevice(devices[0].Id)
	if err != nil {
		t.Fatal(err)
	}
	if device.Name != "Chainsaw" {
		t.Fatal("retrieved wrong device")
	}
}

func TestGetDeviceIdByOperationId(t *testing.T) {
	store := testStorage(t)
	store.AddDevice(models.DeviceInput{Name: "Mower"})
	devices, _ := store.GetDevices()

	store.AddDeviceOperation(devices[0].Id, models.DeviceOperationInput{
		StartTime: "2024-01-01 10:00",
		EndTime:   "2024-01-01 11:00",
	})

	device, _ := store.GetDevice(devices[0].Id)
	opID := device.OperationHistory[0].Id

	id, err := store.GetDeviceIdByOperationId(opID)
	if err != nil {
		t.Fatal(err)
	}
	if id != device.Id {
		t.Fatal("wrong device ID returned for operation")
	}
}

func TestUpdateDevice(t *testing.T) {
	store := testStorage(t)
	store.AddDevice(models.DeviceInput{Name: "Old"})
	devices, _ := store.GetDevices()

	device := devices[0]
	device.Name = "New"
	if err := store.UpdateDevice(device); err != nil {
		t.Fatal(err)
	}

	updated, _ := store.GetDevice(device.Id)
	if updated.Name != "New" {
		t.Fatal("device was not updated")
	}
}

func TestDeleteDevice(t *testing.T) {
	store := testStorage(t)
	store.AddDevice(models.DeviceInput{Name: "Mower"})
	devices, _ := store.GetDevices()

	if err := store.DeleteDevice(devices[0].Id); err != nil {
		t.Fatal(err)
	}

	result, _ := store.GetDevices()
	if len(result) != 0 {
		t.Fatal("device was not deleted")
	}
}

func TestEnrichDevice(t *testing.T) {
	device := models.Device{
		DeviceInput: models.DeviceInput{PurchasePrice: 120},
		OperationHistory: []models.DeviceOperation{
			{DeviceOperationInput: models.DeviceOperationInput{StartTime: "2024-01-03T12:00", EndTime: "2024-01-03T12:30"}},
			{DeviceOperationInput: models.DeviceOperationInput{StartTime: "2024-01-01T10:00", EndTime: "2024-01-01T11:00"}},
		},
	}

	enrichDevice(&device)

	if device.TotalOperations != 2 || device.TotalOperationTime != 90 {
		t.Fatalf("enrichment failed: ops %d, time %f", device.TotalOperations, device.TotalOperationTime)
	}
	if device.LastUsageDate != "2024-01-03" {
		t.Fatal("incorrect last usage date")
	}
	if device.PricePerHour != 80 {
		t.Fatalf("incorrect price per hour, got %f", device.PricePerHour)
	}
}

func TestSortDevices(t *testing.T) {
	devices := []models.Device{
		{Id: 2, DeviceInput: models.DeviceInput{PurchaseDate: "2024-01-01"}},
		{Id: 1, DeviceInput: models.DeviceInput{PurchaseDate: "2025-01-01"}},
	}

	sorted := sortDevices(devices)
	if sorted[0].Id != 1 {
		t.Fatal("devices not sorted correctly by purchase date")
	}
}

func TestLatestDeviceOperationTimestamp(t *testing.T) {
	device := models.Device{
		DeviceInput: models.DeviceInput{PurchaseDate: "2024-01-01"},
		OperationHistory: []models.DeviceOperation{
			{DeviceOperationInput: models.DeviceOperationInput{EndTime: "2024-02-01 10:00"}},
		},
	}

	if res := latestDeviceOperationTimestamp(device); res != "2024-02-01 10:00" {
		t.Fatal("wrong timestamp, got", res)
	}

	// Test case for no operations
	device.OperationHistory = nil
	if res := latestDeviceOperationTimestamp(device); res != "2024-01-01" {
		t.Fatal("wrong timestamp for empty history, got", res)
	}
}
