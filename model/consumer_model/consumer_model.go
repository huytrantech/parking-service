package consumer_model

type SyncEsJobRequest struct {
	ParkingId int    `json:"parking_id"`
	Source    string `json:"source"`
}