package language

import (
	"testing"
)

func TestLoadLanguages(t *testing.T) {

	err := LoadLanguages()

	if err != nil {
		t.Fatal(err)
	}

	got := T(
		"en",
		"missing_key",
	)

	if got != "missing_key" {
		t.Fatal(
			"fallback failed",
		)
	}
}

func TestTranslationExists(t *testing.T) {

	err := LoadLanguages()

	if err != nil {
		t.Fatal(err)
	}

	got := T(
		"en",
		"title",
	)

	if got == "title" {
		t.Fatal(
			"translation missing",
		)
	}
}

func TestUnknownLanguage(t *testing.T) {

	err := LoadLanguages()

	if err != nil {
		t.Fatal(err)
	}

	got :=
		T(
			"fr",
			"hello",
		)

	if got != "hello" {
		t.Fatal(
			"should fallback to key",
		)
	}
}
