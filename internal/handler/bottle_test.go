package handler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/lukas-arnold/garden-equipment-log/internal/models"
)

func TestHandleBottlesView(t *testing.T) {

	h := testHandler(t)

	err := h.storage.AddBottle(
		models.BottleInput{
			PurchaseDate:  "2026-01-01",
			PurchasePrice: 100,
			InitialWeight: 10,
			FillingWeight: 8,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(
		"GET",
		"/bottles",
		nil,
	)

	rec := httptest.NewRecorder()

	h.HandleBottlesView(
		rec,
		req,
	)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"got %d",
			rec.Code,
		)
	}
}

func TestAddBottleGet(t *testing.T) {

	h := testHandler(t)

	req := httptest.NewRequest(
		"GET",
		"/bottle/add",
		nil,
	)

	rec := httptest.NewRecorder()

	h.HandleAddBottleGet(
		rec,
		req,
	)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"got %d",
			rec.Code,
		)
	}
}

func TestAddBottlePost(t *testing.T) {

	h := testHandler(t)

	form := url.Values{}

	form.Set(
		"purchaseDate",
		"2024-01-01",
	)

	form.Set(
		"purchasePrice",
		"100",
	)

	form.Set(
		"initialWeight",
		"10",
	)

	form.Set(
		"fillingWeight",
		"8",
	)

	req := httptest.NewRequest(
		"POST",
		"/bottle/add",
		strings.NewReader(
			form.Encode(),
		),
	)

	req.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)

	rec := httptest.NewRecorder()

	h.HandleAddBottlePost(
		rec,
		req,
	)

	if rec.Code != http.StatusFound {
		t.Fatalf(
			"got %d",
			rec.Code,
		)
	}

	bottles, err := h.storage.GetBottles()

	if err != nil {
		t.Fatal(err)
	}

	if len(bottles) != 1 {
		t.Fatal(
			"bottle missing",
		)
	}

	if bottles[0].PurchasePrice != 100 {
		t.Fatal(
			"wrong bottle",
		)
	}
}

func TestEditBottle(t *testing.T) {

	h := testHandler(t)

	h.storage.AddBottle(
		models.BottleInput{
			PurchaseDate: "2024",
		},
	)

	bottles, _ := h.storage.GetBottles()

	req := httptest.NewRequest(
		"GET",
		"/bottle/edit",
		nil,
	)

	req.SetPathValue(
		"id",
		strconv.FormatInt(
			bottles[0].Id,
			10,
		),
	)

	rec := httptest.NewRecorder()

	h.HandleEditBottle(
		rec,
		req,
	)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"got %d",
			rec.Code,
		)
	}
}

func TestEditBottleInvalidID(t *testing.T) {

	h := testHandler(t)

	req := httptest.NewRequest(
		"GET",
		"/bottle/edit",
		nil,
	)

	req.SetPathValue(
		"id",
		"abc",
	)

	rec := httptest.NewRecorder()

	h.HandleEditBottle(
		rec,
		req,
	)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf(
			"got %d",
			rec.Code,
		)
	}
}

func TestSaveBottle(t *testing.T) {

	h := testHandler(t)

	h.storage.AddBottle(
		models.BottleInput{
			PurchasePrice: 10,
		},
	)

	bottles, _ := h.storage.GetBottles()

	form := url.Values{}

	form.Set(
		"purchaseDate",
		"2025",
	)

	form.Set(
		"purchasePrice",
		"200",
	)

	form.Set(
		"initialWeight",
		"12",
	)

	form.Set(
		"fillingWeight",
		"9",
	)

	req := httptest.NewRequest(
		"POST",
		"/bottle/save",
		strings.NewReader(
			form.Encode(),
		),
	)

	req.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)

	req.SetPathValue(
		"id",
		strconv.FormatInt(
			bottles[0].Id,
			10,
		),
	)

	rec := httptest.NewRecorder()

	h.HandleSaveBottle(
		rec,
		req,
	)

	if rec.Code != http.StatusFound {
		t.Fatalf(
			"got %d",
			rec.Code,
		)
	}

	updated, _ := h.storage.GetBottles()

	if updated[0].PurchasePrice != 200 {
		t.Fatal(
			"bottle not updated",
		)
	}
}

func TestSaveBottleInvalidID(t *testing.T) {

	h := testHandler(t)

	req := httptest.NewRequest(
		"POST",
		"/bottle/save",
		nil,
	)

	req.SetPathValue(
		"id",
		"abc",
	)

	rec := httptest.NewRecorder()

	h.HandleSaveBottle(
		rec,
		req,
	)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf(
			"got %d",
			rec.Code,
		)
	}
}

func TestDeleteBottle(t *testing.T) {

	h := testHandler(t)

	h.storage.AddBottle(
		models.BottleInput{},
	)

	bottles, _ := h.storage.GetBottles()

	req := httptest.NewRequest(
		"GET",
		"/bottle/delete",
		nil,
	)

	req.SetPathValue(
		"id",
		strconv.FormatInt(
			bottles[0].Id,
			10,
		),
	)

	rec := httptest.NewRecorder()

	h.HandleDeleteBottle(
		rec,
		req,
	)

	if rec.Code != http.StatusFound {
		t.Fatalf(
			"got %d",
			rec.Code,
		)
	}

	result, _ := h.storage.GetBottles()

	if len(result) != 0 {
		t.Fatal(
			"bottle not deleted",
		)
	}
}

func TestDeleteBottleInvalidID(t *testing.T) {

	h := testHandler(t)

	req := httptest.NewRequest(
		"GET",
		"/bottle/delete",
		nil,
	)

	req.SetPathValue(
		"id",
		"abc",
	)

	rec := httptest.NewRecorder()

	h.HandleDeleteBottle(
		rec,
		req,
	)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf(
			"got %d",
			rec.Code,
		)
	}
}

func TestBottleHistory(t *testing.T) {

	h := testHandler(t)

	h.storage.AddBottle(
		models.BottleInput{
			FillingWeight: 8,
			InitialWeight: 10,
		},
	)

	bottles, _ := h.storage.GetBottles()

	req := httptest.NewRequest(
		"GET",
		"/bottle/history",
		nil,
	)

	req.SetPathValue(
		"id",
		strconv.FormatInt(
			bottles[0].Id,
			10,
		),
	)

	rec := httptest.NewRecorder()

	h.HandleBottleHistory(
		rec,
		req,
	)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"got %d",
			rec.Code,
		)
	}
}

func TestBottleHistoryInvalidID(t *testing.T) {

	h := testHandler(t)

	req := httptest.NewRequest(
		"GET",
		"/bottle/history",
		nil,
	)

	req.SetPathValue(
		"id",
		"abc",
	)

	rec := httptest.NewRecorder()

	h.HandleBottleHistory(
		rec,
		req,
	)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf(
			"got %d",
			rec.Code,
		)
	}
}
