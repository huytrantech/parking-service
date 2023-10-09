package elastic_search

type ParkingModel struct {
	ParkingId    int            `json:"parking_id"`
	PublicId     string         `json:"public_id"`
	Name         string         `json:"name"`
	CityId       int            `json:"city_id"`
	DistrictId   int            `json:"district_id"`
	WardId       int            `json:"ward_id"`
	Location     Location       `json:"location"`
	Status       int            `json:"status"`
	Distance     float64        `json:"-"`
	ParkingTypes []ParkingTypes `json:"parking_types"`
	ParkingPhone string         `json:"parking_phone"`
	Address      string         `json:"address"`
	OpenAt       int            `json:"open_at"`
	CloseAt      int            `json:"close_at"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type ParkingTypes struct {
	Type   int     `json:"type"`
	Status int     `json:"status"`
	Price  float64 `json:"price"`
}
