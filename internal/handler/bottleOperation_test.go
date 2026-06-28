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

func TestHandleAddBottleOperationGet(t *testing.T) {
	h := testHandler(t)

	req := httptest.NewRequest("GET", "/bottle/operation/add", nil)
	req.SetPathValue("bottleId", "1")
	rec := httptest.NewRecorder()

	h.HandleAddBottleOperationGet(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status OK, got %d", rec.Code)
	}
}

func TestHandleAddBottleOperationGetInvalidID(t *testing.T) {
	h := testHandler(t)

	req := httptest.NewRequest("GET", "/bottle/operation/add", nil)
	req.SetPathValue("bottleId", "abc")
	rec := httptest.NewRecorder()

	h.HandleAddBottleOperationGet(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected InternalServerError, got %d", rec.Code)
	}
}

func TestHandleAddBottleOperationPost(t *testing.T) {
	h := testHandler(t)
	err := h.storage.AddBottle(models.BottleInput{InitialWeight: 10, FillingWeight: 8})
	if err != nil {
		t.Fatal(err)
	}
	bottles, _ := h.storage.GetBottles()

	form := url.Values{}
	form.Set("date", "2026-01-02")
	form.Set("weight", "7")

	req := httptest.NewRequest("POST", "/bottle/operation/add", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetPathValue("bottleId", strconv.FormatInt(bottles[0].Id, 10))
	rec := httptest.NewRecorder()

	h.HandleAddBottleOperationPost(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected status Found, got %d", rec.Code)
	}

	result, _ := h.storage.GetBottle(bottles[0].Id)
	if len(result.OperationHistory) != 1 {
		t.Fatal("operation missing from history")
	}
}

func TestHandleEditBottleOperation(t *testing.T) {
	h := testHandler(t)
	h.storage.AddBottle(models.BottleInput{InitialWeight: 10, FillingWeight: 8})
	bottles, _ := h.storage.GetBottles()
	h.storage.AddBottleOperation(bottles[0].Id, models.BottleOperationInput{Date: "2026-01-01", Weight: 9})

	bottle, _ := h.storage.GetBottle(bottles[0].Id)
	opID := bottle.OperationHistory[0].Id

	req := httptest.NewRequest("GET", "/bottle/operation/edit", nil)
	req.SetPathValue("id", strconv.FormatInt(opID, 10))
	rec := httptest.NewRecorder()

	h.HandleEditBottleOperation(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status OK, got %d", rec.Code)
	}
}

func TestHandleEditBottleOperationInvalidID(t *testing.T) {
	h := testHandler(t)

	req := httptest.NewRequest("GET", "/bottle/operation/edit", nil)
	req.SetPathValue("id", "abc")
	rec := httptest.NewRecorder()

	h.HandleEditBottleOperation(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected InternalServerError, got %d", rec.Code)
	}
}

func TestHandleSaveBottleOperation(t *testing.T) {
	h := testHandler(t)
	h.storage.AddBottle(models.BottleInput{InitialWeight: 10, FillingWeight: 8})
	bottles, _ := h.storage.GetBottles()
	h.storage.AddBottleOperation(bottles[0].Id, models.BottleOperationInput{Date: "2026-01-01", Weight: 9})

	bottle, _ := h.storage.GetBottle(bottles[0].Id)
	opID := bottle.OperationHistory[0].Id

	form := url.Values{}
	form.Set("date", "2026-02-01")
	form.Set("weight", "8")

	req := httptest.NewRequest("POST", "/bottle/operation/save", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetPathValue("id", strconv.FormatInt(opID, 10))
	rec := httptest.NewRecorder()

	h.HandleSaveBottleOperation(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected status Found, got %d", rec.Code)
	}
}

func TestHandleSaveBottleOperationInvalidID(t *testing.T) {
	h := testHandler(t)

	req := httptest.NewRequest("POST", "/bottle/operation/save", nil)
	req.SetPathValue("id", "abc")
	rec := httptest.NewRecorder()

	h.HandleSaveBottleOperation(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected InternalServerError, got %d", rec.Code)
	}
}

func TestHandleDeleteBottleOperation(t *testing.T) {
	h := testHandler(t)
	h.storage.AddBottle(models.BottleInput{})
	bottles, _ := h.storage.GetBottles()
	h.storage.AddBottleOperation(bottles[0].Id, models.BottleOperationInput{})

	bottle, _ := h.storage.GetBottle(bottles[0].Id)
	opID := bottle.OperationHistory[0].Id

	req := httptest.NewRequest("GET", "/bottle/operation/delete", nil)
	req.SetPathValue("id", strconv.FormatInt(opID, 10))
	rec := httptest.NewRecorder()

	h.HandleDeleteBottleOperation(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected status Found, got %d", rec.Code)
	}
}

func TestHandleDeleteBottleOperationInvalidID(t *testing.T) {
	h := testHandler(t)

	req := httptest.NewRequest("GET", "/bottle/operation/delete", nil)
	req.SetPathValue("id", "abc")
	rec := httptest.NewRecorder()

	h.HandleDeleteBottleOperation(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected InternalServerError, got %d", rec.Code)
	}
}
