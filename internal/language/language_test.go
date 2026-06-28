package language

import (
	"testing"
)

func TestLoadLanguages(t *testing.T) {
	if err := LoadLanguages(); err != nil {
		t.Fatalf("failed to load languages: %v", err)
	}

	t.Run("fallback on missing key", func(t *testing.T) {
		got := T("en", "missing_key")
		if got != "missing_key" {
			t.Errorf("expected fallback to key, got %q", got)
		}
	})

	t.Run("translation exists", func(t *testing.T) {
		got := T("en", "title")
		if got == "title" {
			t.Errorf("expected actual translation, got the key itself")
		}
	})

	t.Run("unknown language fallback", func(t *testing.T) {
		got := T("fr", "hello")
		if got != "hello" {
			t.Errorf("expected fallback to key for unknown language, got %q", got)
		}
	})
}
