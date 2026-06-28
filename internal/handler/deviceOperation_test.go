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

func TestHandleAddDeviceOperationGet(t *testing.T) {
	h := testHandler(t)

	req := httptest.NewRequest("GET", "/device/operation/add", nil)
	req.SetPathValue("deviceId", "1")
	rec := httptest.NewRecorder()

	h.HandleAddDeviceOperationGet(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status OK, got %d", rec.Code)
	}
}

func TestHandleAddDeviceOperationPost(t *testing.T) {
	h := testHandler(t)
	err := h.storage.AddDevice(models.DeviceInput{Name: "Chainsaw"})
	if err != nil {
		t.Fatal(err)
	}
	devices, _ := h.storage.GetDevices()

	form := url.Values{}
	form.Set("startTime", "2026-01-01 10:00")
	form.Set("endTime", "2026-01-01 11:00")
	form.Set("note", "test")

	req := httptest.NewRequest("POST", "/device/operation/add", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetPathValue("deviceId", strconv.FormatInt(devices[0].Id, 10))
	rec := httptest.NewRecorder()

	h.HandleAddDeviceOperationPost(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected status Found, got %d", rec.Code)
	}

	device, _ := h.storage.GetDevice(devices[0].Id)
	if len(device.OperationHistory) != 1 {
		t.Fatal("operation missing from history")
	}
}

func TestHandleEditDeviceOperation(t *testing.T) {
	h := testHandler(t)
	h.storage.AddDevice(models.DeviceInput{Name: "Trimmer"})
	devices, _ := h.storage.GetDevices()
	h.storage.AddDeviceOperation(devices[0].Id, models.DeviceOperationInput{
		StartTime: "2026-01-01 10:00",
		EndTime:   "2026-01-01 11:00",
	})

	device, _ := h.storage.GetDevice(devices[0].Id)
	opID := device.OperationHistory[0].Id

	req := httptest.NewRequest("GET", "/device/operation/edit", nil)
	req.SetPathValue("id", strconv.FormatInt(opID, 10))
	rec := httptest.NewRecorder()

	h.HandleEditDeviceOperation(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status OK, got %d", rec.Code)
	}
}

func TestHandleEditDeviceOperationInvalidID(t *testing.T) {
	h := testHandler(t)

	req := httptest.NewRequest("GET", "/device/operation/edit", nil)
	req.SetPathValue("id", "abc")
	rec := httptest.NewRecorder()

	h.HandleEditDeviceOperation(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected InternalServerError, got %d", rec.Code)
	}
}

func TestHandleSaveDeviceOperation(t *testing.T) {
	h := testHandler(t)
	h.storage.AddDevice(models.DeviceInput{Name: "Saw"})
	devices, _ := h.storage.GetDevices()
	h.storage.AddDeviceOperation(devices[0].Id, models.DeviceOperationInput{})

	device, _ := h.storage.GetDevice(devices[0].Id)
	opID := device.OperationHistory[0].Id

	form := url.Values{}
	form.Set("startTime", "2026-02-01 10:00")
	form.Set("endTime", "2026-02-01 12:00")
	form.Set("note", "updated")

	req := httptest.NewRequest("POST", "/device/operation/save", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetPathValue("id", strconv.FormatInt(opID, 10))
	rec := httptest.NewRecorder()

	h.HandleSaveDeviceOperation(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected status Found, got %d", rec.Code)
	}
}

func TestHandleSaveDeviceOperationInvalidID(t *testing.T) {
	h := testHandler(t)

	req := httptest.NewRequest("POST", "/device/operation/save", nil)
	req.SetPathValue("id", "abc")
	rec := httptest.NewRecorder()

	h.HandleSaveDeviceOperation(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected InternalServerError, got %d", rec.Code)
	}
}

func TestHandleDeleteDeviceOperation(t *testing.T) {
	h := testHandler(t)
	h.storage.AddDevice(models.DeviceInput{})
	devices, _ := h.storage.GetDevices()
	h.storage.AddDeviceOperation(devices[0].Id, models.DeviceOperationInput{})

	device, _ := h.storage.GetDevice(devices[0].Id)
	opID := device.OperationHistory[0].Id

	req := httptest.NewRequest("GET", "/device/operation/delete", nil)
	req.SetPathValue("id", strconv.FormatInt(opID, 10))
	rec := httptest.NewRecorder()

	h.HandleDeleteDeviceOperation(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected status Found, got %d", rec.Code)
	}
}

func TestHandleDeleteDeviceOperationInvalidID(t *testing.T) {
	h := testHandler(t)

	req := httptest.NewRequest("GET", "/device/operation/delete", nil)
	req.SetPathValue("id", "abc")
	rec := httptest.NewRecorder()

	h.HandleDeleteDeviceOperation(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected InternalServerError, got %d", rec.Code)
	}
}
