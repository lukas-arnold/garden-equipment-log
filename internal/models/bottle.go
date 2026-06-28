package models

type BottleInput struct {
	PurchaseDate  string  `json:"purchaseDate"`
	PurchasePrice float64 `json:"purchasePrice"`
	InitialWeight float64 `json:"initialWeight"`
	FillingWeight float64 `json:"fillingWeight"`
}

type Bottle struct {
	Id int64 `json:"id"`
	BottleInput
	OperationHistory []BottleOperation `json:"operationHistory"`
	TotalOperations  int
	RestGas          float64
	UsedGas          float64
}

type BottleOperationInput struct {
	Date   string  `json:"date"`
	Weight float64 `json:"weight"`
}

type BottleOperation struct {
	Id int64 `json:"id"`
	BottleOperationInput
}
