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

func TestHandleDevicesView(t *testing.T) {

	h := testHandler(t)

	req :=
		httptest.NewRequest(
			"GET",
			"/devices",
			nil,
		)

	rec :=
		httptest.NewRecorder()

	h.HandleDevicesView(
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

func TestHandleAddDeviceGet(t *testing.T) {

	h := testHandler(t)

	req :=
		httptest.NewRequest(
			"GET",
			"/device/add",
			nil,
		)

	rec :=
		httptest.NewRecorder()

	h.HandleAddDeviceGet(
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

func TestHandleAddDevicePost(t *testing.T) {

	h := testHandler(t)

	form := url.Values{}

	form.Set(
		"name",
		"Mower",
	)

	form.Set(
		"purchaseDate",
		"2024-01-01",
	)

	form.Set(
		"purchasePrice",
		"100",
	)

	req :=
		httptest.NewRequest(
			"POST",
			"/device/add",
			strings.NewReader(
				form.Encode(),
			),
		)

	req.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)

	rec :=
		httptest.NewRecorder()

	h.HandleAddDevicePost(
		rec,
		req,
	)

	if rec.Code != http.StatusFound {
		t.Fatalf(
			"got %d",
			rec.Code,
		)
	}

	devices, err :=
		h.storage.GetDevices()

	if err != nil {
		t.Fatal(err)
	}

	if len(devices) != 1 {
		t.Fatal(
			"device missing",
		)
	}

	if devices[0].Name != "Mower" {
		t.Fatal(
			"wrong device",
		)
	}
}

func TestHandleEditDevice(t *testing.T) {

	h := testHandler(t)

	h.storage.AddDevice(
		models.DeviceInput{
			Name: "Mower",
		},
	)

	devices, _ :=
		h.storage.GetDevices()

	req :=
		httptest.NewRequest(
			"GET",
			"/device/edit",
			nil,
		)

	req.SetPathValue(
		"id",
		strconv.FormatInt(
			devices[0].Id,
			10,
		),
	)

	rec :=
		httptest.NewRecorder()

	h.HandleEditDevice(
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

func TestHandleEditDeviceInvalidID(t *testing.T) {

	h := testHandler(t)

	req :=
		httptest.NewRequest(
			"GET",
			"/device/edit",
			nil,
		)

	req.SetPathValue(
		"id",
		"abc",
	)

	rec :=
		httptest.NewRecorder()

	h.HandleEditDevice(
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

func TestHandleSaveDevice(t *testing.T) {

	h := testHandler(t)

	h.storage.AddDevice(
		models.DeviceInput{
			Name: "Old",
		},
	)

	devices, _ :=
		h.storage.GetDevices()

	form := url.Values{}

	form.Set(
		"name",
		"New",
	)

	form.Set(
		"purchaseDate",
		"2024-01-01",
	)

	form.Set(
		"purchasePrice",
		"200",
	)

	req :=
		httptest.NewRequest(
			"POST",
			"/device/save",
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
			devices[0].Id,
			10,
		),
	)

	rec :=
		httptest.NewRecorder()

	h.HandleSaveDevice(
		rec,
		req,
	)

	if rec.Code != http.StatusFound {
		t.Fatalf(
			"got %d",
			rec.Code,
		)
	}

	updated, _ :=
		h.storage.GetDevices()

	if updated[0].Name != "New" {
		t.Fatal(
			"device not updated",
		)
	}
}

func TestHandleSaveDeviceInvalidID(t *testing.T) {

	h := testHandler(t)

	req :=
		httptest.NewRequest(
			"POST",
			"/device/save",
			nil,
		)

	req.SetPathValue(
		"id",
		"abc",
	)

	rec :=
		httptest.NewRecorder()

	h.HandleSaveDevice(
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

func TestHandleDeleteDevice(t *testing.T) {

	h := testHandler(t)

	h.storage.AddDevice(
		models.DeviceInput{
			Name: "Mower",
		},
	)

	devices, _ :=
		h.storage.GetDevices()

	req :=
		httptest.NewRequest(
			"GET",
			"/device/delete",
			nil,
		)

	req.SetPathValue(
		"id",
		strconv.FormatInt(
			devices[0].Id,
			10,
		),
	)

	rec :=
		httptest.NewRecorder()

	h.HandleDeleteDevice(
		rec,
		req,
	)

	if rec.Code != http.StatusFound {
		t.Fatalf(
			"got %d",
			rec.Code,
		)
	}

	result, _ :=
		h.storage.GetDevices()

	if len(result) != 0 {
		t.Fatal(
			"device not deleted",
		)
	}
}

func TestHandleDeleteDeviceInvalidID(t *testing.T) {

	h := testHandler(t)

	req :=
		httptest.NewRequest(
			"GET",
			"/device/delete",
			nil,
		)

	req.SetPathValue(
		"id",
		"abc",
	)

	rec :=
		httptest.NewRecorder()

	h.HandleDeleteDevice(
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

func TestHandleDeviceHistory(t *testing.T) {

	h := testHandler(t)

	h.storage.AddDevice(
		models.DeviceInput{
			Name: "Mower",
		},
	)

	devices, _ :=
		h.storage.GetDevices()

	req :=
		httptest.NewRequest(
			"GET",
			"/device/history",
			nil,
		)

	req.SetPathValue(
		"id",
		strconv.FormatInt(
			devices[0].Id,
			10,
		),
	)

	rec :=
		httptest.NewRecorder()

	h.HandleDeviceHistory(
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

func TestHandleDeviceHistoryInvalidID(t *testing.T) {

	h := testHandler(t)

	req :=
		httptest.NewRequest(
			"GET",
			"/device/history",
			nil,
		)

	req.SetPathValue(
		"id",
		"abc",
	)

	rec :=
		httptest.NewRecorder()

	h.HandleDeviceHistory(
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
