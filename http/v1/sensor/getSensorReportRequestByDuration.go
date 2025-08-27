package sensor

import "time"

type getSensorReportRequestByIDandDuration struct {
	ID1   string    `json:"id1" validate:"required"`
	ID2   int       `json:"id2" validate:"required"`
	Start time.Time `json:"start" validate:"required"`
	End   time.Time `json:"end" validate:"required"`
}
