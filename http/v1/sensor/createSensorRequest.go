package sensor

type CreateSensorRequest struct {
	SensorType string `json:"sensorType" validate:"required"`
	ID1        string `json:"id1" validate:"required"`
	ID2        int    `json:"id2" validate:"required"`
}
