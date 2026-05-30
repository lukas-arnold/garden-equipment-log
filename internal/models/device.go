package models

type DeviceInput struct {
	Name          string  `json:"Name"`
	PurchaseDate  string  `json:"PurchaseDate"`
	PurchasePrice float64 `json:"PurchasePrice"`
}

type Device struct {
	Id int64 `json:"Id"`
	DeviceInput
	OperationHistory    []DeviceOperation `json:"OperationHistory"`
	LastUsageDate       string            `json:"-"`
	TotalOperations     int               `json:"-"`
	TotalOperationHours float64           `json:"-"`
	PricePerHour        float64           `json:"-"`
}

type DeviceOperationInput struct {
	StartTime string `json:"StartTime"`
	EndTime   string `json:"EndTime"`
	Note      string `json:"Note"`
}

type DeviceOperation struct {
	Id int64 `json:"Id"`
	DeviceOperationInput
}
