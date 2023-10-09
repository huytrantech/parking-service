package api_model

type ListCityResponseDto struct {
	CityId   int    `json:"city_id"`
	CityName string `json:"city_name"`
}

type ListDistrictResponseDto struct {
	ListCityResponseDto
	DistrictId     int    `json:"district_id"`
	DistrictName   string `json:"district_name"`
}

type ListWardResponseDto struct {
	ListDistrictResponseDto
	WardId         int    `json:"ward_id"`
	WardName       string `json:"ward_name"`
}
