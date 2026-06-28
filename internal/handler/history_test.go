package handler

import (
	"testing"

	"github.com/lukas-arnold/garden-equipment-log/internal/models"
)

func TestBuildBottleHistoryRows(t *testing.T) {

	bottle := models.Bottle{
		BottleInput: models.BottleInput{
			InitialWeight: 10,
			FillingWeight: 8,
		},
		OperationHistory: []models.BottleOperation{
			{
				Id: 1,
				BottleOperationInput: models.BottleOperationInput{
					Date:   "2024-01-01",
					Weight: 6,
				},
			},
		},
	}

	rows := buildBottleHistoryRows(
		bottle,
	)

	if len(rows) != 1 {
		t.Fatal(
			"expected one row",
		)
	}

	if rows[0].Id != 1 {
		t.Fatal(
			"wrong operation",
		)
	}

	// initial 10kg -> current 6kg = 4kg used
	if rows[0].UsedGas != 4 {
		t.Fatalf(
			"expected used gas 4 got %f",
			rows[0].UsedGas,
		)
	}

	// filling weight 8kg - used 4kg = 4kg remaining
	if rows[0].RestGas != 4 {
		t.Fatalf(
			"expected rest gas 4 got %f",
			rows[0].RestGas,
		)
	}
}

func TestBuildBottleHistoryRowsMultiple(t *testing.T) {

	bottle := models.Bottle{
		BottleInput: models.BottleInput{
			InitialWeight: 15,
			FillingWeight: 12,
		},
		OperationHistory: []models.BottleOperation{
			{
				Id: 1,
				BottleOperationInput: models.BottleOperationInput{
					Date:   "2024-01-01",
					Weight: 10,
				},
			},
			{
				Id: 2,
				BottleOperationInput: models.BottleOperationInput{
					Date:   "2024-01-02",
					Weight: 8,
				},
			},
		},
	}

	rows := buildBottleHistoryRows(
		bottle,
	)

	if len(rows) != 2 {
		t.Fatal(
			"expected two rows",
		)
	}

	if rows[0].UsedGas != 5 {
		t.Fatalf(
			"expected first used gas 5 got %f",
			rows[0].UsedGas,
		)
	}

	if rows[1].UsedGas != 7 {
		t.Fatalf(
			"expected second used gas 7 got %f",
			rows[1].UsedGas,
		)
	}

	if rows[1].RestGas != 5 {
		t.Fatalf(
			"expected second rest gas 5 got %f",
			rows[1].RestGas,
		)
	}
}

func TestBuildBottleHistoryRowsEmpty(t *testing.T) {

	bottle := models.Bottle{}

	rows := buildBottleHistoryRows(
		bottle,
	)

	if len(rows) != 0 {
		t.Fatal(
			"expected no rows",
		)
	}
}

func TestBuildDeviceHistoryRows(t *testing.T) {

	device := models.Device{
		OperationHistory: []models.DeviceOperation{
			{
				Id: 1,
				DeviceOperationInput: models.DeviceOperationInput{
					StartTime: "2024-01-01T10:00",
					EndTime:   "2024-01-01T10:30",
					Note:      "test",
				},
			},
		},
	}

	rows := buildDeviceHistoryRows(
		device,
	)

	if len(rows) != 1 {
		t.Fatal(
			"expected one row",
		)
	}

	if rows[0].Id != 1 {
		t.Fatal(
			"wrong operation",
		)
	}

	if rows[0].Time != 30 {
		t.Fatalf(
			"expected 30 minutes got %f",
			rows[0].Time,
		)
	}
}

func TestBuildDeviceHistoryRowsEmptyTime(t *testing.T) {

	device := models.Device{
		OperationHistory: []models.DeviceOperation{
			{
				Id: 1,
				DeviceOperationInput: models.DeviceOperationInput{
					Note: "missing time",
				},
			},
		},
	}

	rows := buildDeviceHistoryRows(
		device,
	)

	if len(rows) != 1 {
		t.Fatal(
			"expected row",
		)
	}

	if rows[0].Time != 0 {
		t.Fatalf(
			"expected zero time got %f",
			rows[0].Time,
		)
	}
}

func TestBuildDeviceHistoryRowsInvalidTime(t *testing.T) {

	device := models.Device{
		OperationHistory: []models.DeviceOperation{
			{
				Id: 1,
				DeviceOperationInput: models.DeviceOperationInput{
					StartTime: "abc",
					EndTime:   "xyz",
				},
			},
		},
	}

	rows := buildDeviceHistoryRows(
		device,
	)

	if rows[0].Time != 0 {
		t.Fatalf(
			"expected zero time got %f",
			rows[0].Time,
		)
	}
}

func TestBuildDeviceHistoryRowsNegativeDuration(t *testing.T) {

	device := models.Device{
		OperationHistory: []models.DeviceOperation{
			{
				Id: 1,
				DeviceOperationInput: models.DeviceOperationInput{
					StartTime: "2024-01-01T12:00",
					EndTime:   "2024-01-01T11:00",
				},
			},
		},
	}

	rows := buildDeviceHistoryRows(
		device,
	)

	if rows[0].Time != 0 {
		t.Fatalf(
			"expected negative duration clamped to zero got %f",
			rows[0].Time,
		)
	}
}
