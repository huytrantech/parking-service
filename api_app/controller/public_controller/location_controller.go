package public_controller

import (
	"github.com/labstack/echo/v4"
	"parking-service/core"
	"parking-service/model/api_model"
	"parking-service/provider/location_provider"
	"sort"
	"strconv"
)

type ILocationController interface {
	GetCitiesController(c echo.Context) error
	GetDistrictsController(c echo.Context) error
	GetWardsController(c echo.Context) error
}

type locationController struct {
	ILocationProvider location_provider.ILocationProvider
}

func NewLocationController(ILocationProvider location_provider.ILocationProvider) ILocationController {
	return &locationController{ILocationProvider: ILocationProvider}
}

// GetCityListController godoc
// @Summary Get All City
// @Description Get All City
// @Tags public
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Accept json
// @Produce json
// @Success 200 {object} api_model.BaseResponseAPIModel{data=[]api_model.ListCityResponseDto}
// @Router /api/public/v1/location/cities [get]
func (ctr *locationController) GetCitiesController(c echo.Context) error {
	dataCity := ctr.ILocationProvider.GetLocationData()
	response := make([]api_model.ListCityResponseDto, len(dataCity))
	for index, value := range dataCity {
		response[index] = api_model.ListCityResponseDto{
			CityId:   value.Id,
			CityName: value.Name,
		}
	}

	return api_model.SuccessResponse(c, response)
}

// GetDistrictsController godoc
// @Summary Get All District By CityId
// @Description Get All District By CityId
// @Tags public
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Accept json
// @Produce json
// @Param city_id path int true "city_id"
// @Success 200 {object} api_model.BaseResponseAPIModel{data=[]api_model.ListDistrictResponseDto}
// @Router /api/public/v1/location/cities/{city_id}/districts [get]
func (ctr *locationController) GetDistrictsController(c echo.Context) error {
	dataMap := ctr.ILocationProvider.GetMapData()
	cityId, err := strconv.Atoi(c.Param("city_id"))
	if err != nil || cityId <= 0 {
		return api_model.FailResponse(c, core.NewBadRequestErrorMessage("Bad Request CityId"))
	}
	dataDistrictByCity := dataMap[cityId].DistrictMap
	response := make([]api_model.ListDistrictResponseDto, 0)
	cityDto := api_model.ListCityResponseDto{
		CityId:   cityId,
		CityName: dataMap[cityId].CityName,
	}
	for key, value := range dataDistrictByCity {
		item := api_model.ListDistrictResponseDto{
			DistrictId:   key,
			DistrictName: value.DistrictName,
		}
		item.ListCityResponseDto = cityDto
		response = append(response, item)
	}
	sort.SliceStable(response, func(i, j int) bool {
		return response[i].DistrictId < response[j].DistrictId
	})

	return api_model.SuccessResponse(c, response)
}

// GetWardsController godoc
// @Summary Get All Ward By CityId and DistrictId
// @Description Get All Ward By CityId and DistrictId
// @Tags public
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param city_id path int true "city_id"
// @Param district_Id path int true "district_id"
// @Accept json
// @Produce json
// @Success 200 {object} api_model.BaseResponseAPIModel{data=[]api_model.ListWardResponseDto}
// @Router /api/public/v1/location/cities/{city_id}/districts/{district_id}/wards [get]
func (ctr *locationController) GetWardsController(c echo.Context) error {
	dataMap := ctr.ILocationProvider.GetMapData()
	cityId, err := strconv.Atoi(c.Param("city_id"))
	if err != nil || cityId <= 0 {
		return api_model.FailResponse(c, core.NewBadRequestErrorMessage("Bad Request CityId"))
	}
	districtId, err := strconv.Atoi(c.Param("district_id"))
	if err != nil || districtId <= 0 {
		return api_model.FailResponse(c, core.NewBadRequestErrorMessage("Bad Request DistrictId"))
	}
	dataDistrictByCity := dataMap[cityId].DistrictMap
	cityDto := api_model.ListCityResponseDto{
		CityId:   cityId,
		CityName: dataMap[cityId].CityName,
	}
	dataWardByDistrict := dataDistrictByCity[districtId].WardMap
	districtDto := api_model.ListDistrictResponseDto{
		ListCityResponseDto: cityDto,
		DistrictId:          districtId,
		DistrictName:        dataDistrictByCity[districtId].DistrictName,
	}
	response := make([]api_model.ListWardResponseDto, 0)
	for key, value := range dataWardByDistrict {
		response = append(response, api_model.ListWardResponseDto{
			WardName:                value.WardName,
			WardId:                  key,
			ListDistrictResponseDto: districtDto,
		})
	}
	sort.SliceStable(response, func(i, j int) bool {
		return response[i].WardId < response[j].WardId
	})

	return api_model.SuccessResponse(c, response)
}
