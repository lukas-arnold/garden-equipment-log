package storage

import (
	"path/filepath"
	"testing"

	"github.com/lukas-arnold/garden-equipment-log/internal/models"
)

func TestGetDeviceChart(t *testing.T) {

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
			StartTime: "2024-01-01T10:00",
			EndTime:   "2024-01-01T11:30",
		},
	)

	chart, err := store.GetDeviceChart(
		devices[0].Id,
	)

	if err != nil {
		t.Fatal(err)
	}

	if chart.Type != "bar" {
		t.Fatal("wrong chart type")
	}

	if len(chart.Labels) != 1 {
		t.Fatal("missing label")
	}

	if chart.Labels[0] != "2024" {
		t.Fatal("wrong year")
	}

	if chart.Sets[0].Data[0] != 90 {
		t.Fatalf(
			"got %f",
			chart.Sets[0].Data[0],
		)
	}
}

func TestGetDeviceChartEmpty(t *testing.T) {

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

	chart, err := store.GetDeviceChart(
		devices[0].Id,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(chart.Labels) != 0 {
		t.Fatal("expected empty chart")
	}
}

func TestGetDeviceChartSorted(t *testing.T) {

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
			StartTime: "2025-01-01T10:00",
			EndTime:   "2025-01-01T11:00",
		},
	)

	store.AddDeviceOperation(
		devices[0].Id,
		models.DeviceOperationInput{
			StartTime: "2024-01-01T10:00",
			EndTime:   "2024-01-01T11:00",
		},
	)

	chart, _ := store.GetDeviceChart(
		devices[0].Id,
	)

	if chart.Labels[0] != "2024" {
		t.Fatal("not sorted")
	}

	if chart.Labels[1] != "2025" {
		t.Fatal("not sorted")
	}
}

func TestGetDeviceChartMultipleOperationsSameYear(t *testing.T) {

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
			StartTime: "2024-01-01T10:00",
			EndTime:   "2024-01-01T10:30",
		},
	)

	store.AddDeviceOperation(
		devices[0].Id,
		models.DeviceOperationInput{
			StartTime: "2024-06-01T12:00",
			EndTime:   "2024-06-01T12:45",
		},
	)

	chart, err := store.GetDeviceChart(
		devices[0].Id,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(chart.Sets) != 1 {
		t.Fatal("missing dataset")
	}

	if chart.Sets[0].Data[0] != 75 {
		t.Fatalf(
			"got %f",
			chart.Sets[0].Data[0],
		)
	}
}

func TestGetBottleChart(t *testing.T) {

	store := New(
		filepath.Join(
			t.TempDir(),
			"storage.json",
		),
	)

	store.AddBottle(
		models.BottleInput{
			PurchaseDate:  "2024-01-01",
			InitialWeight: 15,
			FillingWeight: 5,
		},
	)

	bottles, _ := store.GetBottles()

	store.AddBottleOperation(
		bottles[0].Id,
		models.BottleOperationInput{
			Date:   "2024-02-01",
			Weight: 14,
		},
	)

	chart, err := store.GetBottleChart(
		bottles[0].Id,
	)

	if err != nil {
		t.Fatal(err)
	}

	if chart.Type != "line" {
		t.Fatal("wrong chart type")
	}

	if len(chart.Labels) != 2 {
		t.Fatal("wrong labels")
	}

	if chart.Sets[0].Data[0] != 15 {
		t.Fatal("wrong initial weight")
	}

	if chart.Sets[0].Data[1] != 14 {
		t.Fatal("wrong operation weight")
	}

	if chart.Sets[1].Data[0] != 10 {
		t.Fatal("wrong empty weight")
	}
}

func TestGetBottleChartEmpty(t *testing.T) {

	store := New(
		filepath.Join(
			t.TempDir(),
			"storage.json",
		),
	)

	store.AddBottle(
		models.BottleInput{
			PurchaseDate:  "2024-01-01",
			InitialWeight: 15,
			FillingWeight: 5,
		},
	)

	bottles, _ := store.GetBottles()

	chart, err := store.GetBottleChart(
		bottles[0].Id,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(chart.Labels) != 1 {
		t.Fatal("wrong labels")
	}

	if chart.Labels[0] != "2024-01-01" {
		t.Fatal("wrong purchase date")
	}

	if len(chart.Sets[0].Data) != 1 {
		t.Fatal("wrong weights")
	}

	if chart.Sets[0].Data[0] != 15 {
		t.Fatal("wrong initial weight")
	}
}

func TestGetBottleChartSorted(t *testing.T) {

	store := New(
		filepath.Join(
			t.TempDir(),
			"storage.json",
		),
	)

	store.AddBottle(
		models.BottleInput{
			PurchaseDate:  "2024-01-01",
			InitialWeight: 15,
			FillingWeight: 5,
		},
	)

	bottles, _ := store.GetBottles()

	store.AddBottleOperation(
		bottles[0].Id,
		models.BottleOperationInput{
			Date:   "2024-03-01",
			Weight: 13,
		},
	)

	store.AddBottleOperation(
		bottles[0].Id,
		models.BottleOperationInput{
			Date:   "2024-02-01",
			Weight: 14,
		},
	)

	chart, _ := store.GetBottleChart(
		bottles[0].Id,
	)

	if chart.Labels[1] != "2024-02-01" {
		t.Fatal("not sorted")
	}

	if chart.Labels[2] != "2024-03-01" {
		t.Fatal("not sorted")
	}
}
