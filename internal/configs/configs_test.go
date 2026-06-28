package configs

import (
	"os"
	"testing"
)

func TestGetPort(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		os.Unsetenv("PORT")
		if got := GetPort(); got != ":8080" {
			t.Errorf("expected :8080, got %s", got)
		}
	})

	t.Run("environment", func(t *testing.T) {
		os.Setenv("PORT", ":9999")
		defer os.Unsetenv("PORT")
		if got := GetPort(); got != ":9999" {
			t.Errorf("expected :9999, got %s", got)
		}
	})
}

func TestGetStorageFile(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		os.Unsetenv("STORAGE_FILE")
		if got := GetStorageFile(); got != "data/strength-tracker.json" {
			t.Errorf("expected data/strength-tracker.json, got %s", got)
		}
	})

	t.Run("environment", func(t *testing.T) {
		os.Setenv("STORAGE_FILE", "/tmp/test.json")
		defer os.Unsetenv("STORAGE_FILE")
		if got := GetStorageFile(); got != "/tmp/test.json" {
			t.Errorf("expected /tmp/test.json, got %s", got)
		}
	})
}

func TestGetLanguage(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		os.Unsetenv("LANGUAGE")
		if got := GetLanguage(); got != "en" {
			t.Errorf("expected en, got %s", got)
		}
	})

	t.Run("environment", func(t *testing.T) {
		os.Setenv("LANGUAGE", "de")
		defer os.Unsetenv("LANGUAGE")
		if got := GetLanguage(); got != "de" {
			t.Errorf("expected de, got %s", got)
		}
	})
}
