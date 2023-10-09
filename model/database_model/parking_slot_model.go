package database_model

type ParkingSlotModel struct {
	Id          int         `db:"id"`
	ParkingId   int         `db:"parking_id"`
	ParkingType int         `db:"parking_type"`
	Price       float64     `db:"price"`
	TotalSlot   *int64      `db:"total_slot"`
	CurrentSlot *int64      `db:"current_slot"`
	Status      int         `db:"status"`
	Metadata    interface{} `db:"metadata"`
}
