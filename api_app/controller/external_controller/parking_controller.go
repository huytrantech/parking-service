package external_controller

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
	"parking-service/constants"
	"parking-service/core"
	"parking-service/model/api_model"
	"parking-service/model/proxy_model/goong"
	"parking-service/service/parking_module/create_parking"
	"parking-service/service/parking_module/detail_parking"
	"parking-service/service/parking_module/get_circle_parking_location"
	"parking-service/service/parking_module/get_direction"
	"parking-service/service/parking_module/recommend_parking_location"
	"parking-service/service/place_module/detail_place"
	"parking-service/service/place_module/get_placeholder"
	"strings"
)

type IExternalParkingController interface {
	GetCircleParkingLocationController(c echo.Context) error
	GetDirectionController(c echo.Context) error
	RecommendParkingLocationController(c echo.Context) error
	GetDetailParkingController(c echo.Context) error
	GetPlaceHolder(c echo.Context) error
	GetDetailPlaceHolder(c echo.Context) error
	AddParking(c echo.Context) error
}

type externalParkingController struct {
	IGetCircleParkingLocationService get_circle_parking_location.IGetCircleParkingLocationService
	IGetDirectionService             get_direction.IGetDirectionService
	IRecommendParkingLocationService recommend_parking_location.IRecommendParkingLocationService
	IDetailParkingService            detail_parking.IDetailParkingService
	IGetPlaceHolderService           get_placeholder.IGetPlaceHolderService
	IDetailPlaceService              detail_place.IDetailPlaceService
	ICreateParkingService            create_parking.ICreateParkingService
}

func NewExternalParkingController(
	IGetCircleParkingLocationService get_circle_parking_location.IGetCircleParkingLocationService,
	IGetDirectionService get_direction.IGetDirectionService,
	ICreateParkingService create_parking.ICreateParkingService,
	IRecommendParkingLocationService recommend_parking_location.IRecommendParkingLocationService,
	IGetPlaceHolderService get_placeholder.IGetPlaceHolderService,
	IDetailPlaceService detail_place.IDetailPlaceService,
	IDetailParkingService detail_parking.IDetailParkingService) IExternalParkingController {
	return &externalParkingController{
		IGetCircleParkingLocationService: IGetCircleParkingLocationService,
		IGetDirectionService:             IGetDirectionService,
		IRecommendParkingLocationService: IRecommendParkingLocationService,
		IGetPlaceHolderService:           IGetPlaceHolderService,
		IDetailPlaceService:              IDetailPlaceService,
		ICreateParkingService:            ICreateParkingService,
		IDetailParkingService:            IDetailParkingService}
}

// GetCircleParkingLocationController godoc
// @Summary Top nearest parking in circle radius
// @Description Show the top nearest parking in circle radius of specific location
// @Tags external
// @Security ApiKeyAuth
// @Param lat query float64 false "lat"
// @Param lng query float64 false "lng"
// @Param distance query int false "distance"
// @Param text_search query string false "text_search"
// @Param city query int false "city"
// @Param district query int false "district"
// @param Authorization header string true "Authorization"
// @Accept json
// @Produce json
// @Success 200 {object} api_model.BaseResponseAPIModel{data=[]api_model.GetCircleParkingLocationModelResponseDto}
// @Router /api/external/v1/circle-location [get]
func (ctr *externalParkingController) GetCircleParkingLocationController(c echo.Context) error {
	request := api_model.GetCircleParkingLocationRequestDto{
		Lat:        cast.ToFloat64(c.QueryParam("lat")),
		Lng:        cast.ToFloat64(c.QueryParam("lng")),
		Distance:   cast.ToInt(c.QueryParam("distance")),
		TextSearch: c.QueryParam("text_search"),
		CityId:     cast.ToInt(c.QueryParam("city")),
		DistrictId: cast.ToInt(c.QueryParam("district")),
	}
	parkingTypes := c.QueryParam("parking_types")
	if len(parkingTypes) > 0 {
		parkingTypesInt := make([]int, len(strings.Split(parkingTypes, ",")))
		for index, value := range strings.Split(parkingTypes, ",") {
			if cast.ToInt(value) > 0 {
				parkingTypesInt[index] = cast.ToInt(value)
			}
		}
		request.ParkingTypes = parkingTypesInt
	}
	data, svErr := ctr.IGetCircleParkingLocationService.GetCircleParkingLocation(context.Background(), request)
	if svErr != nil {
		return api_model.FailResponse(c, svErr)
	}
	return api_model.SuccessResponse(c, data)
}

// RecommendParkingLocationController godoc
// @Summary Show the recommend parking
// @Description  Show the recommend parking
// @Tags external
// @Security ApiKeyAuth
// @Param page_index query int false "page_index"
// @Param page_limit query int false "page_limit"
// @Param city_id query int false "city_id"
// @Param lat query int false "lat"
// @Param lng query int false "lng"
// @param Authorization header string true "Authorization"
// @Accept json
// @Produce json
// @Success 200 {object} api_model.BaseResponseAPIModel{data=[]api_model.RecommendParkingLocationBaseResponseDto}
// @Router /api/external/v1/recommend [get]
func (ctr *externalParkingController) RecommendParkingLocationController(c echo.Context) error {
	request := api_model.RecommendParkingLocationRequestDto{
		PageIndex: cast.ToInt16(c.QueryParam("page_index")),
		PageLimit: cast.ToInt16(c.QueryParam("page_limit")),
		CityId:    cast.ToInt(c.QueryParam("city_id")),
		Lat:       cast.ToFloat64(c.QueryParam("lat")),
		Lng:       cast.ToFloat64(c.QueryParam("lng")),
	}

	data, svErr := ctr.IRecommendParkingLocationService.RecommendParkingLocation(context.Background(), request)
	if svErr != nil {
		return api_model.FailResponse(c, svErr)
	}
	return api_model.SuccessResponse(c, data)
}

// GetDirectionController godoc
// @Summary Show the direction from specific location to parking
// @Description  Show the direction from specific location to parking
// @Tags external
// @Security ApiKeyAuth
// @Param origin query string false "origin"
// @Param parking_id query string false "parking_id"
// @param Authorization header string true "Authorization"
// @Accept json
// @Produce json
// @Success 200 {object} api_model.BaseResponseAPIModel{data=api_model.GetDirectionResponseDto}
// @Router /api/external/v1/directions [get]
func (ctr *externalParkingController) GetDirectionController(c echo.Context) error {
	originStr := c.QueryParam("origin")
	parkingId := c.QueryParam("parking_id")
	originArr := strings.Split(originStr, ",")
	if len(originArr) != 2 || len(parkingId) == 0 {
		return api_model.FailResponse(c, core.NewBadRequestErrorMessage("Bad Request"))
	}
	request := api_model.GetDirectionRequestDto{
		Origin: api_model.PointLocationApiDto{
			Lat: cast.ToFloat64(originArr[0]),
			Lng: cast.ToFloat64(originArr[1]),
		},
		ParkingId: parkingId,
	}
	data, svErr := ctr.IGetDirectionService.GetDirection(context.Background(), request)
	if svErr != nil {
		return api_model.FailResponse(c, svErr)
	}
	return api_model.SuccessResponse(c, data)
}

// GetDetailParkingController godoc
// @Summary Parking detail information
// @Description  Show  the parking detail information
// @Tags external
// @Security ApiKeyAuth
// @Param parking_id path string true "parking_id"
// @param Authorization header string true "Authorization"
// @Accept json
// @Produce json
// @Success 200 {object} api_model.BaseResponseAPIModel{data=api_model.DetailParkingResponseDto}
// @Router /api/external/v1/{parking_id} [get]
func (ctr *externalParkingController) GetDetailParkingController(c echo.Context) error {
	parkingId := c.Param("parking_id")
	data, svErr := ctr.IDetailParkingService.GetDetailParking(context.Background(), parkingId)
	if svErr != nil {
		return api_model.FailResponse(c, svErr)
	}
	return api_model.SuccessResponse(c, data)
}

// GetPlaceHolder godoc
// @Summary GetPlaceHolder
// @Description  GetPlaceHolder
// @Tags external
// @Security ApiKeyAuth
// @Param input query string true "input"
// @Param limit query int false "limit"
// @Param radius query int false "radius"
// @param Authorization header string true "Authorization"
// @Accept json
// @Produce json
// @Success 200 {object} api_model.BaseResponseAPIModel{data=[]api_model.GetPlaceHolderResponseDto}
// @Router /api/external/v1/placeholder [get]
func (ctr *externalParkingController) GetPlaceHolder(c echo.Context) error {
	request := goong.GetPlaceHolderRequestDto{
		Input:  c.QueryParam("input"),
		Limit:  cast.ToInt(c.QueryParam("limit")),
		Radius: cast.ToInt(c.QueryParam("radius")),
	}
	data, svErr := ctr.IGetPlaceHolderService.GetPlaceHolder(context.Background(), request)
	if svErr != nil {
		return api_model.FailResponse(c, svErr)
	}
	return api_model.SuccessResponse(c, data)
}

// GetDetailPlaceHolder godoc
// @Summary GetDetailPlaceHolder
// @Description  GetDetailPlaceHolder
// @Tags external
// @Security ApiKeyAuth
// @Param place_id path string true "place_id"
// @param Authorization header string true "Authorization"
// @Accept json
// @Produce json
// @Success 200 {object} api_model.BaseResponseAPIModel{data=[]api_model.DetailPlaceResponseDto}
// @Router /api/external/v1/placeholder/{place_id}/detail [get]
func (ctr *externalParkingController) GetDetailPlaceHolder(c echo.Context) error {
	placeId := c.Param("place_id")
	data, svErr := ctr.IDetailPlaceService.GetDetailPlaceByPlaceId(context.Background(), placeId)
	if svErr != nil {
		return api_model.FailResponse(c, svErr)
	}
	return api_model.SuccessResponse(c, data)
}

// AddParking godoc
// @Summary AddParking
// @Description  AddParking
// @Tags external
// @Security ApiKeyAuth
// @param model body api_model.CreateParkingRequestDto true "model"
// @param Authorization header string true "Authorization"
// @Accept json
// @Produce json
// @Success 200 {object} api_model.BaseResponseAPIModel{data=boolean}
// @Router /api/external/v1/add [get]
func (ctr *externalParkingController) AddParking(c echo.Context) error {
	var request api_model.CreateParkingRequestDto
	if err := c.Bind(&request); err != nil {
		return api_model.FailResponse(c, err)
	}
	request.Username = cast.ToString(c.Get(constants.HEADER_USERNAME))
	response, svErr := ctr.ICreateParkingService.CreateParking(context.Background(), request)
	if svErr != nil {
		return api_model.FailResponse(c, svErr)
	}
	return api_model.SuccessResponse(c, response)
}
