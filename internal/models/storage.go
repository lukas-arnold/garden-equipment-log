package models

type EquipmentStorage struct {
	Devices []Device `json:"devices"`
	Bottles []Bottle `json:"bottles"`
}

type ChartDataset struct {
	Label string
	Data  []float64
	Unit  string
}

type ChartModel struct {
	Type   string
	Labels []string
	Sets   []ChartDataset
}
