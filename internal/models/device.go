package models

type DeviceInput struct {
	Name          string  `json:"name"`
	PurchaseDate  string  `json:"purchaseDate"`
	PurchasePrice float64 `json:"purchasePrice"`
}

type Device struct {
	Id int64 `json:"id"`
	DeviceInput
	OperationHistory   []DeviceOperation `json:"operationHistory"`
	LastUsageDate      string
	TotalOperations    int
	TotalOperationTime float64
	PricePerHour       float64
}

type DeviceOperationInput struct {
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Note      string `json:"note"`
}

type DeviceOperation struct {
	Id int64 `json:"id"`
	DeviceOperationInput
}
