package models

type BottleInput struct {
	PurchaseDate  string  `json:"PurchaseDate"`
	PurchasePrice float64 `json:"PurchasePrice"`
	InitialWeight float64 `json:"InitialWeight"`
	FillingWeight float64 `json:"FillingWeight"`
}

type Bottle struct {
	Id int64 `json:"Id"`
	BottleInput
	OperationHistory []BottleOperation `json:"OperationHistory"`
	TotalOperations  int               `json:"-"`
	RestGas          float64           `json:"-"`
	UsedGas          float64           `json:"-"`
}

type BottleOperationInput struct {
	Date   string  `json:"Date"`
	Weight float64 `json:"Weight"`
}

type BottleOperation struct {
	Id int64 `json:"Id"`
	BottleOperationInput
}
