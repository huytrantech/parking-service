package api_model

import (
	"errors"
	"fmt"
	"parking-service/core/enums/parking_type_car_enum"
	"parking-service/core/utils/text_utils"
	"time"
)

type CreateParkingRequestDto struct {
	CityId        int                            `json:"city_id"`
	DistrictId    int                            `json:"district_id"`
	WardId        int                            `json:"ward_id"`
	Address       string                         `json:"address"`
	Lat           float64                        `json:"lat"`
	Lng           float64                        `json:"lng"`
	OwnerName     string                         `json:"owner_name"`
	OwnerPhone    string                         `json:"owner_phone"`
	ParkingName   string                         `json:"parking_name"`
	ParkingPhone  string                         `json:"parking_phone"`
	Username      string                         `json:"-"`
	OpenAtHour    int                            `json:"open_at_hour"`
	OpenAtMinute  int                            `json:"open_at_minute"`
	CloseAtHour   int                            `json:"close_at_hour"`
	CloseAtMinute int                            `json:"close_at_minute"`
	Images        []string                       `json:"images"`
	ParkingTypes  []CreateParkingTypesRequestDto `json:"parking_types"`
}

type CreateParkingTypesRequestDto struct {
	ParkingType int     `json:"parking_type"`
	TotalSlot   *int64  `json:"total_slot"`
	Price       float64 `json:"price"`
}

func (request CreateParkingRequestDto) Invalid() error {

	if request.CityId <= 0 || request.DistrictId <= 0 {
		return errors.New("Thông tin location không hợp lệ")
	}

	if len(request.Address) == 0 {
		return errors.New("Thông tin địa chỉ không hợp lệ")
	}

	if len(request.ParkingName) == 0 {
		return errors.New("Thông tin liên hệ parking không hợp lệ")
	}

	if len(request.ParkingPhone) > 0 && !text_utils.CheckRegexPhoneNumber(request.ParkingPhone) {
		return errors.New("Số điện thoại liên hệ không hợp lệ")
	}

	if len(request.OwnerPhone) > 0 && !text_utils.CheckRegexPhoneNumber(request.OwnerPhone) {
		return errors.New("Số điện thoại người quản lý không hợp lệ")
	}

	for _, v := range request.ParkingTypes {
		if len(parking_type_car_enum.GetTypeNameFromData(v.ParkingType)) == 0 {
			return errors.New("Thông tin bãi đỗ xe không đúng")
		}
	}
	if request.CloseAtHour < request.OpenAtHour {
		return errors.New("Thông tin thời gian hoạt động không hợp lệ")
	}
	return nil
}

type GetCircleParkingLocationRequestDto struct {
	Lat          float64
	Lng          float64
	Distance     int
	TextSearch   string
	CityId       int
	DistrictId   int
	ParkingTypes []int
}

func (request GetCircleParkingLocationRequestDto) ConvertSearchCircleLocationRequest() map[string]interface{} {
	searchFilter := make([]interface{}, 0)
	sort := make([]interface{}, 0)
	if request.Lat != 0 && request.Lng != 0 {
		searchFilter = append(searchFilter, map[string]interface{}{
			"geo_distance": map[string]interface{}{
				"distance": fmt.Sprintf("%dkm", request.Distance),
				"location": map[string]interface{}{
					"lat": request.Lat,
					"lon": request.Lng,
				},
			},
		})
		sort = append(sort, map[string]interface{}{
			"_geo_distance": map[string]interface{}{
				"location": map[string]interface{}{
					"lat": request.Lat,
					"lon": request.Lng,
				},
				"order":           "asc",
				"unit":            "m",
				"mode":            "min",
				"distance_type":   "arc",
				"ignore_unmapped": true,
			},
		})
	}

	if len(request.ParkingTypes) > 0 {
		searchFilter = append(searchFilter, map[string]interface{}{
			"nested": map[string]interface{}{
				"path": "parking_types",
				"query": map[string]interface{}{
					"terms": map[string]interface{}{
						"parking_types.type": request.ParkingTypes,
					},
				},
			},
		})
	}

	if len(request.TextSearch) > 0 {
		shouldFilter := make([]interface{}, 0)
		shouldFilter = append(shouldFilter, map[string]interface{}{
			"match": map[string]interface{}{
				"name": request.TextSearch,
			},
		})
		shouldFilter = append(shouldFilter, map[string]interface{}{
			"term": map[string]interface{}{
				"parking_phone": request.TextSearch,
			},
		})
		shouldFilter = append(shouldFilter, map[string]interface{}{
			"match": map[string]interface{}{
				"address": request.TextSearch,
			},
		})
		searchFilter = append(searchFilter, map[string]interface{}{
			"bool": map[string]interface{}{
				"should": shouldFilter,
			},
		})
	}
	if request.CityId > 0 {
		searchFilter = append(searchFilter, map[string]interface{}{
			"term": map[string]interface{}{
				"city_id": request.CityId,
			},
		})
	}
	if request.DistrictId > 0 {
		searchFilter = append(searchFilter, map[string]interface{}{
			"term": map[string]interface{}{
				"district_id": request.DistrictId,
			},
		})
	}
	filter := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"filter": searchFilter,
			},
		},
		"sort": sort,
	}
	return filter
}

type GetDirectionRequestDto struct {
	Origin    PointLocationApiDto `json:"origin"`
	ParkingId string              `json:"parking_id"`
}

type RecommendParkingLocationRequestDto struct {
	PageIndex int16   `json:"page_index"`
	PageLimit int16   `json:"page_limit"`
	CityId    int     `json:"city_id"`
	Lat       float64 `json:"lat"`
	Lng       float64 `json:"lng"`
}

func (request RecommendParkingLocationRequestDto) ConvertQueryES() map[string]interface{} {
	searchFilter := make([]interface{}, 0)
	sort := make([]interface{}, 0)

	if request.Lat != 0 && request.Lng != 0 {
		searchFilter = append(searchFilter, map[string]interface{}{
			"geo_distance": map[string]interface{}{
				"distance": fmt.Sprintf("%dkm", 10),
				"location": map[string]interface{}{
					"lat": request.Lat,
					"lon": request.Lng,
				},
			},
		})
		sort = append(sort, map[string]interface{}{
			"_geo_distance": map[string]interface{}{
				"location": map[string]interface{}{
					"lat": request.Lat,
					"lon": request.Lng,
				},
				"order":           "asc",
				"unit":            "m",
				"mode":            "min",
				"distance_type":   "arc",
				"ignore_unmapped": true,
			},
		})
	}

	if request.CityId > 0 {
		searchFilter = append(searchFilter, map[string]interface{}{
			"term": map[string]interface{}{
				"city_id": request.CityId,
			},
		})
	}

	filter := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"filter": searchFilter,
			},
		},
		"sort": sort,
	}
	return filter
}

type RetrieveListParkingRequestDto struct {
	PageIndex int16 `json:"page_index"`
	PageLimit int16 `json:"page_limit"`
}

type UpdateParkingRequestDto struct {
	ParkingName  string     `json:"parking_name"`
	ParkingPhone string     `json:"parking_phone"`
	OwnerName    string     `json:"owner_name"`
	OwnerPhone   string     `json:"owner_phone"`
	CityId       int        `json:"city_id"`
	DistrictId   int        `json:"district_id"`
	WardId       int        `json:"ward_id"`
	Address      string     `json:"address"`
	Lat          float64    `json:"lat"`
	Lng          float64    `json:"lng"`
	Username     string     `json:"-"`
	OpenAt       *time.Time `json:"open_at"`
	CloseAt      *time.Time `json:"close_at"`
}

type GenJWTTokenExternalRequestDto struct {
	Token string `json:"token"`
}

type AddParkingSlotRequest struct {
	Type      int     `json:"type"`
	TotalSlot int64   `json:"total_slot"`
	Price     float64 `json:"price"`
	ParkingId int     `json:"-"`
}

func (request AddParkingSlotRequest) Invalid() error {
	if request.Type == 0 || request.ParkingId <= 0 {
		return errors.New("Bad Request")
	}
	return nil
}

type UpdateParkingSlotRequest struct {
	TotalSlot     int64   `json:"total_slot"`
	ParkingSlotId int     `json:"parking_slot_id"`
	Price         float64 `json:"price"`
	ParkingId     int     `json:"-"`
}
