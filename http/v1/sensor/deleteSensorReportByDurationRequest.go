package sensor

import "time"

type DeleteSensorReportByDurationRequest struct {
	Start time.Time `json:"start" validate:"required"`
	End   time.Time `json:"end" validate:"required"`
}
