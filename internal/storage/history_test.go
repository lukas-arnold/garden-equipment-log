package storage

import (
	"testing"

	"github.com/lukas-arnold/garden-equipment-log/internal/models"
)

func TestGetDeviceChart(t *testing.T) {
	store := testStorage(t)
	store.AddDevice(models.DeviceInput{Name: "Mower"})
	devices, _ := store.GetDevices()

	store.AddDeviceOperation(devices[0].Id, models.DeviceOperationInput{
		StartTime: "2024-01-01T10:00",
		EndTime:   "2024-01-01T11:30",
	})

	chart, err := store.GetDeviceChart(devices[0].Id)
	if err != nil {
		t.Fatal(err)
	}

	if chart.Type != "bar" || len(chart.Labels) != 1 || chart.Labels[0] != "2024" {
		t.Fatalf("unexpected chart structure: %+v", chart)
	}
	if chart.Sets[0].Data[0] != 90 {
		t.Fatalf("expected 90 minutes, got %f", chart.Sets[0].Data[0])
	}
}

func TestGetDeviceChartEmpty(t *testing.T) {
	store := testStorage(t)
	store.AddDevice(models.DeviceInput{Name: "Mower"})
	devices, _ := store.GetDevices()

	chart, err := store.GetDeviceChart(devices[0].Id)
	if err != nil {
		t.Fatal(err)
	}
	if len(chart.Labels) != 0 {
		t.Fatal("expected empty chart labels")
	}
}

func TestGetDeviceChartSorted(t *testing.T) {
	store := testStorage(t)
	store.AddDevice(models.DeviceInput{Name: "Mower"})
	devices, _ := store.GetDevices()

	store.AddDeviceOperation(devices[0].Id, models.DeviceOperationInput{StartTime: "2025-01-01T10:00", EndTime: "2025-01-01T11:00"})
	store.AddDeviceOperation(devices[0].Id, models.DeviceOperationInput{StartTime: "2024-01-01T10:00", EndTime: "2024-01-01T11:00"})

	chart, _ := store.GetDeviceChart(devices[0].Id)
	if chart.Labels[0] != "2024" || chart.Labels[1] != "2025" {
		t.Fatalf("labels not sorted correctly: %v", chart.Labels)
	}
}

func TestGetDeviceChartMultipleOperationsSameYear(t *testing.T) {
	store := testStorage(t)
	store.AddDevice(models.DeviceInput{Name: "Mower"})
	devices, _ := store.GetDevices()

	store.AddDeviceOperation(devices[0].Id, models.DeviceOperationInput{StartTime: "2024-01-01T10:00", EndTime: "2024-01-01T10:30"})
	store.AddDeviceOperation(devices[0].Id, models.DeviceOperationInput{StartTime: "2024-06-01T12:00", EndTime: "2024-06-01T12:45"})

	chart, _ := store.GetDeviceChart(devices[0].Id)
	if len(chart.Sets) != 1 || chart.Sets[0].Data[0] != 75 {
		t.Fatalf("expected 75 minutes, got %f", chart.Sets[0].Data[0])
	}
}

func TestGetBottleChart(t *testing.T) {
	store := testStorage(t)
	store.AddBottle(models.BottleInput{PurchaseDate: "2024-01-01", InitialWeight: 15, FillingWeight: 5})
	bottles, _ := store.GetBottles()

	store.AddBottleOperation(bottles[0].Id, models.BottleOperationInput{Date: "2024-02-01", Weight: 14})

	chart, err := store.GetBottleChart(bottles[0].Id)
	if err != nil {
		t.Fatal(err)
	}

	if chart.Type != "line" || len(chart.Labels) != 2 {
		t.Fatalf("unexpected chart structure: %+v", chart)
	}
	if chart.Sets[0].Data[0] != 15 || chart.Sets[0].Data[1] != 14 || chart.Sets[1].Data[0] != 10 {
		t.Fatal("weight data mismatch")
	}
}

func TestGetBottleChartEmpty(t *testing.T) {
	store := testStorage(t)
	store.AddBottle(models.BottleInput{PurchaseDate: "2024-01-01", InitialWeight: 15, FillingWeight: 5})
	bottles, _ := store.GetBottles()

	chart, _ := store.GetBottleChart(bottles[0].Id)
	if chart.Labels[0] != "2024-01-01" || chart.Sets[0].Data[0] != 15 {
		t.Fatalf("unexpected empty chart data: %+v", chart)
	}
}

func TestGetBottleChartSorted(t *testing.T) {
	store := testStorage(t)
	store.AddBottle(models.BottleInput{PurchaseDate: "2024-01-01", InitialWeight: 15, FillingWeight: 5})
	bottles, _ := store.GetBottles()

	store.AddBottleOperation(bottles[0].Id, models.BottleOperationInput{Date: "2024-03-01", Weight: 13})
	store.AddBottleOperation(bottles[0].Id, models.BottleOperationInput{Date: "2024-02-01", Weight: 14})

	chart, _ := store.GetBottleChart(bottles[0].Id)
	if chart.Labels[1] != "2024-02-01" || chart.Labels[2] != "2024-03-01" {
		t.Fatalf("bottle operations not sorted: %v", chart.Labels)
	}
}
