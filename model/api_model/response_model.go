package api_model

type GetCircleParkingLocationModelResponseDto struct {
	ParkingId     int                       `json:"parking_id"`
	PublicId      string                    `json:"public_id"`
	Name          string                    `json:"name"`
	Status        int                       `json:"status"`
	StatusDisplay string                    `json:"status_display"`
	Address       string                    `json:"address"`
	ParkingPhone  string                    `json:"parking_phone"`
	Distance      float64                   `json:"distance"`
	ParkingTypes  []ParkingTypesResponseDto `json:"parking_types"`
	Location      PointLocationApiDto       `json:"location"`
}

type PointLocationApiDto struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
type GetDirectionResponseDto struct {
	Points []PointLocationApiDto `json:"points"`
}

type RecommendParkingLocationResponseDto struct {
	ParkingName string              `json:"parking_name"`
	ParkingId   string              `json:"parking_id"`
	Location    PointLocationApiDto `json:"location"`
}

type RecommendParkingLocationBaseResponseDto struct {
	Parking []RecommendParkingLocationResponseDto `json:"parking"`
}

type DetailParkingResponseDto struct {
	ParkingId     string                    `json:"parking_id"`
	ParkingName   string                    `json:"parking_name"`
	ParkingPhone  string                    `json:"parking_phone"`
	StatusDisplay string                    `json:"status_display"`
	Status        int                       `json:"status"`
	Address       string                    `json:"address"`
	FullAddress   string                    `json:"full_address"`
	Location      PointLocationApiDto       `json:"location"`
	ParkingTypes  []ParkingTypesResponseDto `json:"parking_types"`
	OpenAtHour    int                       `json:"open_at_hour"`
	CloseAtHour   int                       `json:"close_at_hour"`
	OpenAtMinute  int                       `json:"open_at_minute"`
	CloseAtMinute int                       `json:"close_at_minute"`
	CityId        int                       `json:"city_id"`
	DistrictId    int                       `json:"district_id"`
	WardId        int                       `json:"ward_id"`
	OwnerName     string                    `json:"owner_name"`
	OwnerPhone    string                    `json:"owner_phone"`
}

type ParkingTypesResponseDto struct {
	Id            int     `json:"id"`
	Type          int     `json:"type"`
	Status        int     `json:"status"`
	Logo          string  `json:"logo"`
	StatusDisplay string  `json:"status_display"`
	TypeDisplay   string  `json:"type_display"`
	TotalSlot     *int64  `json:"total_slot"`
	CurrentSlot   *int64  `json:"current_slot"`
	Price         float64 `json:"price"`
}

type RetrieveListParkingResponseDto struct {
	PublicId      string                            `json:"public_id"`
	ParkingName   string                            `json:"parking_name"`
	ParkingPhone  string                            `json:"parking_phone"`
	Status        int                               `json:"status"`
	StatusDisplay string                            `json:"status_display"`
	Roles         RetrieveListParkingActionRolesDto `json:"roles"`
}

type RetrieveListParkingResponseBaseDto struct {
	Parks      []RetrieveListParkingResponseDto `json:"parks"`
	Total      int                              `json:"total"`
	PageIndex  int16                            `json:"page_index"`
	PageLimit  int16                            `json:"page_limit"`
	IsLastPage bool                             `json:"is_last_page"`
}

type RetrieveListParkingActionRolesDto struct {
	CanApprove bool `json:"can_approve"`
	CanRemove  bool `json:"can_remove"`
	CanDenied  bool `json:"can_denied"`
	CanBlock   bool `json:"can_block"`
	CanClose   bool `json:"can_close"`
	CanReopen  bool `json:"can_reopen"`
}

type GetPlaceHolderResponseDto struct {
	Place   string `json:"place"`
	PlaceId string `json:"place_id"`
}

type DetailPlaceResponseDto struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type ImportParkingResponseDto struct {
	TotalSuccess int64 `json:"total_success"`
	TotalFail    int64 `json:"total_fail"`
}
