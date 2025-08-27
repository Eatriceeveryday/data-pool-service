package sensor

import "time"

type getSensorReportRequestByDuration struct {
	Start time.Time `json:"start" validate:"required"`
	End   time.Time `json:"end" validate:"required"`
}
