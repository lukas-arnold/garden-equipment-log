package storage

import (
	"path/filepath"
	"testing"

	"github.com/lukas-arnold/garden-equipment-log/internal/models"
)

func TestAddDeviceOperation(t *testing.T) {

	store := New(
		filepath.Join(
			t.TempDir(),
			"storage.json",
		),
	)

	store.AddDevice(
		models.DeviceInput{
			Name: "Mower",
		},
	)

	devices, _ := store.GetDevices()

	err := store.AddDeviceOperation(
		devices[0].Id,
		models.DeviceOperationInput{
			StartTime: "2024-01-01 10:00",
			EndTime:   "2024-01-01 11:00",
			Note:      "Test",
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	device, _ := store.GetDevice(
		devices[0].Id,
	)

	if len(device.OperationHistory) != 1 {
		t.Fatal("operation missing")
	}
}

func TestGetDeviceOperation(t *testing.T) {

	store := New(
		filepath.Join(
			t.TempDir(),
			"storage.json",
		),
	)

	store.AddDevice(
		models.DeviceInput{
			Name: "Mower",
		},
	)

	devices, _ := store.GetDevices()

	store.AddDeviceOperation(
		devices[0].Id,
		models.DeviceOperationInput{
			StartTime: "2024-01-01 10:00",
			EndTime:   "2024-01-01 11:00",
			Note:      "Test",
		},
	)

	device, _ := store.GetDevice(
		devices[0].Id,
	)

	operation, err := store.GetDeviceOperation(
		device.OperationHistory[0].Id,
	)

	if err != nil {
		t.Fatal(err)
	}

	if operation.Note != "Test" {
		t.Fatal("wrong operation")
	}
}

func TestUpdateDeviceOperation(t *testing.T) {

	store := New(
		filepath.Join(
			t.TempDir(),
			"storage.json",
		),
	)

	store.AddDevice(
		models.DeviceInput{
			Name: "Mower",
		},
	)

	devices, _ := store.GetDevices()

	store.AddDeviceOperation(
		devices[0].Id,
		models.DeviceOperationInput{
			StartTime: "2024-01-01 10:00",
			EndTime:   "2024-01-01 11:00",
			Note:      "Old",
		},
	)

	device, _ := store.GetDevice(
		devices[0].Id,
	)

	err := store.UpdateDeviceOperation(
		device.OperationHistory[0].Id,
		models.DeviceOperationInput{
			StartTime: "2024-01-02 09:00",
			EndTime:   "2024-01-02 10:30",
			Note:      "New",
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	operation, _ := store.GetDeviceOperation(
		device.OperationHistory[0].Id,
	)

	if operation.StartTime != "2024-01-02 09:00" {
		t.Fatal("start time not updated")
	}

	if operation.EndTime != "2024-01-02 10:30" {
		t.Fatal("end time not updated")
	}

	if operation.Note != "New" {
		t.Fatal("note not updated")
	}
}

func TestDeleteDeviceOperation(t *testing.T) {

	store := New(
		filepath.Join(
			t.TempDir(),
			"storage.json",
		),
	)

	store.AddDevice(
		models.DeviceInput{
			Name: "Mower",
		},
	)

	devices, _ := store.GetDevices()

	store.AddDeviceOperation(
		devices[0].Id,
		models.DeviceOperationInput{
			StartTime: "2024-01-01 10:00",
			EndTime:   "2024-01-01 11:00",
		},
	)

	device, _ := store.GetDevice(
		devices[0].Id,
	)

	err := store.DeleteDeviceOperation(
		device.OperationHistory[0].Id,
	)

	if err != nil {
		t.Fatal(err)
	}

	updated, _ := store.GetDevice(
		devices[0].Id,
	)

	if len(updated.OperationHistory) != 0 {
		t.Fatal("operation not deleted")
	}
}

func TestSortDeviceOperations(t *testing.T) {

	operations := []models.DeviceOperation{
		{
			DeviceOperationInput: models.DeviceOperationInput{
				StartTime: "2024-01-01 10:00",
			},
		},
		{
			DeviceOperationInput: models.DeviceOperationInput{
				StartTime: "2024-03-01 10:00",
			},
		},
		{
			DeviceOperationInput: models.DeviceOperationInput{
				StartTime: "2024-02-01 10:00",
			},
		},
	}

	sorted := sortDeviceOperations(
		operations,
	)

	if sorted[0].StartTime != "2024-03-01 10:00" {
		t.Fatal("not sorted")
	}

	if sorted[1].StartTime != "2024-02-01 10:00" {
		t.Fatal("not sorted")
	}

	if sorted[2].StartTime != "2024-01-01 10:00" {
		t.Fatal("not sorted")
	}
}
