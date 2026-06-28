package handler

import (
	"github.com/lukas-arnold/garden-equipment-log/internal/models"
	"github.com/lukas-arnold/garden-equipment-log/internal/utils"
)

type bottleOperationRow struct {
	models.BottleOperation
	RestGas float64
	UsedGas float64
}

type deviceOperationRow struct {
	models.DeviceOperation
	Time float64
}

func buildBottleHistoryRows(
	bottle models.Bottle,
) []bottleOperationRow {

	rows := make(
		[]bottleOperationRow,
		0,
		len(bottle.OperationHistory),
	)

	for _, op := range bottle.OperationHistory {

		usedGas := bottle.InitialWeight - op.Weight

		rows = append(
			rows,
			bottleOperationRow{
				BottleOperation: op,
				RestGas:         bottle.FillingWeight - usedGas,
				UsedGas:         usedGas,
			},
		)
	}

	return rows
}

func buildDeviceHistoryRows(
	device models.Device,
) []deviceOperationRow {
	rows := make(
		[]deviceOperationRow,
		0,
		len(device.OperationHistory),
	)

	for _, op := range device.OperationHistory {
		var minutes float64

		if op.StartTime != "" &&
			op.EndTime != "" {

			start, err1 :=
				utils.ParseDateTime(op.StartTime)

			end, err2 :=
				utils.ParseDateTime(op.EndTime)

			if err1 == nil &&
				err2 == nil {

				minutes =
					end.Sub(start).Minutes()

				if minutes < 0 {
					minutes = 0
				}
			}
		}

		rows = append(
			rows,
			deviceOperationRow{
				DeviceOperation: op,
				Time:            minutes,
			},
		)
	}

	return rows
}
