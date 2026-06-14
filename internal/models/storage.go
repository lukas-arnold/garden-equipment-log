package models

type EquipmentStorage struct {
	Devices []Device `json:"Devices"`
	Bottles []Bottle `json:"Bottles"`
}

type ChartDataset struct {
	Label string    `json:"label"`
	Data  []float64 `json:"data"`
	Unit  string    `json:"unit"`
}

type ChartModel struct {
	Type   string         `json:"type,omitempty"`
	Labels []string       `json:"labels"`
	Sets   []ChartDataset `json:"sets"`
}
