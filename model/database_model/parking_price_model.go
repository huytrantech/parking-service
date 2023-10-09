package database_model

type ParkingPriceModel struct {
	ParkingId int     `db:"parking_id"`
	Name      string  `db:"name"`
	Slots     int     `db:"slots"`
	Price     float64 `db:"price"`
	Id        int     `db:"id"`
}
