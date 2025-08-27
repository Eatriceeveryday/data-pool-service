package sensor

type EditSensorReportRequestById struct {
	ID1   string  `json:"id1" validate:"required"`
	ID2   int     `json:"id2" validate:"required"`
	Value float32 `json:"value" validate:"required"`
}
