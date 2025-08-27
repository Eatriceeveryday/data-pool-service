package sensor

type DeleteSensorReportByIDRequest struct {
	ID1 string `json:"id1" validate:"required"`
	ID2 int    `json:"id2" validate:"required"`
}
