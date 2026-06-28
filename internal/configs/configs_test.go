package configs

import (
	"os"
	"testing"
)

func TestGetPortDefault(t *testing.T) {

	os.Unsetenv("PORT")

	got := GetPort()

	if got != ":8080" {
		t.Fatalf(
			"got %s",
			got,
		)
	}
}

func TestGetPortEnv(t *testing.T) {

	os.Setenv(
		"PORT",
		":9999",
	)

	defer os.Unsetenv("PORT")

	got := GetPort()

	if got != ":9999" {
		t.Fatalf(
			"got %s",
			got,
		)
	}
}

func TestGetStorageFile(t *testing.T) {

	os.Unsetenv("STORAGE_FILE")

	got := GetStorageFile()

	if got != "data/garden-equipment-log.json" {
		t.Fatalf(
			"got %s",
			got,
		)
	}
}

func TestGetStorageFileEnv(t *testing.T) {

	old := os.Getenv("STORAGE_FILE")
	defer os.Setenv("STORAGE_FILE", old)

	os.Setenv(
		"STORAGE_FILE",
		"/tmp/test.json",
	)

	got := GetStorageFile()

	if got != "/tmp/test.json" {
		t.Fatalf(
			"got %s",
			got,
		)
	}
}

func TestGetLanguage(t *testing.T) {

	os.Unsetenv("LANGUAGE")

	got := GetLanguage()

	if got != "en" {
		t.Fatalf(
			"got %s",
			got,
		)
	}
}

func TestGetLanguageEnv(t *testing.T) {

	old := os.Getenv("LANGUAGE")
	defer os.Setenv("LANGUAGE", old)

	os.Setenv(
		"LANGUAGE",
		"de",
	)

	got := GetLanguage()

	if got != "de" {
		t.Fatalf(
			"got %s",
			got,
		)
	}
}
