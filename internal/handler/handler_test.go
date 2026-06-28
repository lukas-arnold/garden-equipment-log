package handler

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/lukas-arnold/garden-equipment-log/internal/storage"
)

func testHandler(t *testing.T) *Handler {

	t.Helper()

	store := storage.New(
		filepath.Join(
			t.TempDir(),
			"test.json",
		),
	)

	return New(store)
}

func TestTemplateFuncs(t *testing.T) {

	funcs := getTemplateFuncs()

	if funcs == nil {
		t.Fatal("expected template funcs")
	}

	for name, fn := range funcs {

		if fn == nil {
			t.Fatalf(
				"template func %s is nil",
				name,
			)
		}
	}
}

func TestTemplateFuncsFormat(t *testing.T) {

	funcs := getTemplateFuncs()

	formatFloat :=
		funcs["formatFloat"].(func(float64, int) string)

	formatEuro :=
		funcs["formatEuro"].(func(float64) string)

	formatTime :=
		funcs["formatOperationTime"].(func(float64) string)

	if formatFloat(12.345, 2) != "12,35" {
		t.Fatal(
			"formatFloat failed",
		)
	}

	if formatEuro(19.5) != "19,50 €" {
		t.Fatal(
			"formatEuro failed",
		)
	}

	if formatTime(125) != "2 h 5 min" {
		t.Fatal(
			"formatOperationTime failed",
		)
	}
}

func TestFormatOperationTimeHoursOnly(t *testing.T) {

	funcs := getTemplateFuncs()

	format :=
		funcs["formatOperationTime"].(func(float64) string)

	if format(120) != "2 h" {
		t.Fatal(
			"wrong format",
		)
	}
}

func TestFormatOperationTimeMinutesOnly(t *testing.T) {

	funcs := getTemplateFuncs()

	format :=
		funcs["formatOperationTime"].(func(float64) string)

	if format(45) != "45 min" {
		t.Fatal(
			"wrong format",
		)
	}
}

func TestTemplateTranslationFunc(t *testing.T) {

	funcs := getTemplateFuncs()

	tFunc :=
		funcs["T"].(func(string) string)

	result := tFunc("something")

	if result == "" {
		t.Fatal(
			"expected translation result",
		)
	}
}

func TestRenderTemplate(t *testing.T) {

	h := testHandler(t)

	rec :=
		httptest.NewRecorder()

	h.renderTemplate(
		rec,
		"templates/device/add.html",
		nil,
	)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected 200 got %d",
			rec.Code,
		)
	}
}

func TestRenderTemplateWithData(t *testing.T) {

	h := testHandler(t)

	rec :=
		httptest.NewRecorder()

	h.renderTemplate(
		rec,
		"templates/device/edit.html",
		struct {
			Name string
		}{
			Name: "Test",
		},
	)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected 200 got %d",
			rec.Code,
		)
	}
}

func TestRenderHistoryTemplate(t *testing.T) {

	h := testHandler(t)

	rec :=
		httptest.NewRecorder()

	h.renderHistoryTemplate(
		rec,
		"templates/device/history.html",
		nil,
	)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected 200 got %d",
			rec.Code,
		)
	}
}

func TestRenderHistoryTemplateWithData(t *testing.T) {

	h := testHandler(t)

	rec :=
		httptest.NewRecorder()

	h.renderHistoryTemplate(
		rec,
		"templates/device/history.html",
		struct {
			Name string
		}{
			Name: "History",
		},
	)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected 200 got %d",
			rec.Code,
		)
	}
}
