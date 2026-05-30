package models

type EquipmentStorage struct {
	Devices []Device `json:"Devices"`
	Bottles []Bottle `json:"Bottles"`
}

type DevicesForChart struct {
	Devices []Device
	Labels  []string
	Hours   [][]float64
}

type BottlesForChart struct {
	Bottles     []Bottle
	Names       []string
	LastWeights []float64
	Dates       []string
	Weights     [][]float64
}

type DeviceOperationTimesForChart struct {
	Device    Device
	Dates     []string
	Durations []float64
}

type BottleWeightHistoryForChart struct {
	Bottle  Bottle
	Dates   []string
	Weights []float64
}
