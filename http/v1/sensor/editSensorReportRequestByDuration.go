package sensor

import "time"

type EditSensorReportRequestByDuration struct {
	Start time.Time `json:"start" validate:"required"`
	End   time.Time `json:"end" validate:"required"`
	Value float32   `json:"value" validate:"required"`
}
