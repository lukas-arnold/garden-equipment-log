package utils

import (
	"testing"
	"time"

	"github.com/lukas-arnold/garden-equipment-log/internal/models"
)

func TestConvertEquipmentStorageToBytes(t *testing.T) {

	input := models.EquipmentStorage{}

	bytes, err := ConvertEquipmentStorageToBytes(input)

	if err != nil {
		t.Fatal(err)
	}

	if len(bytes) == 0 {
		t.Fatal("expected json bytes")
	}
}

func TestConvertBytesToEquipmentStorage(t *testing.T) {

	data := []byte(`{}`)

	storage, err := ConvertBytesToEquipmentStorage(data)

	if err != nil {
		t.Fatal(err)
	}

	_ = storage
}

func TestConvertBytesToEquipmentStorageInvalid(t *testing.T) {

	_, err := ConvertBytesToEquipmentStorage(
		[]byte("not json"),
	)

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestConvertId(t *testing.T) {

	got, err := ConvertId("123")

	if err != nil {
		t.Fatal(err)
	}

	if got != 123 {
		t.Fatalf(
			"got %d want %d",
			got,
			123,
		)
	}
}

func TestConvertIdInvalid(t *testing.T) {

	_, err := ConvertId("abc")

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestConvertFloat(t *testing.T) {

	got, err := ConvertFloat("42.5")

	if err != nil {
		t.Fatal(err)
	}

	if got != 42.5 {
		t.Fatalf(
			"got %f want %f",
			got,
			42.5,
		)
	}
}

func TestConvertFloatInvalid(t *testing.T) {

	_, err := ConvertFloat("abc")

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseDateTime(t *testing.T) {

	got, err := ParseDateTime(
		"2026-01-01T12:30",
	)

	if err != nil {
		t.Fatal(err)
	}

	expected := time.Date(
		2026,
		time.January,
		1,
		12,
		30,
		0,
		0,
		time.UTC,
	)

	if !got.Equal(expected) {
		t.Fatalf(
			"got %v want %v",
			got,
			expected,
		)
	}
}

func TestParseDateTimeInvalid(t *testing.T) {

	_, err := ParseDateTime(
		"2026-01-01",
	)

	if err == nil {
		t.Fatal("expected error")
	}
}
