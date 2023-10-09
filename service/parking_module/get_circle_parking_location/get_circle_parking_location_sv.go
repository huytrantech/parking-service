package get_circle_parking_location

import (
	"context"
	"parking-service/constants"
	"parking-service/core"
	"parking-service/core/enums/parking_status_enum"
	"parking-service/core/enums/parking_type_car_enum"
	"parking-service/core/enums/parking_types_car_status_enum"
	"parking-service/model/api_model"
	"parking-service/model/proxy_model/elastic_search"
	"parking-service/proxy/elastic_search_proxy"
)

type IGetCircleParkingLocationService interface {
	GetCircleParkingLocation(ctx context.Context, request api_model.GetCircleParkingLocationRequestDto) (
		parking []api_model.GetCircleParkingLocationModelResponseDto, errRes error)
}

type getCircleParkingLocationService struct {
	IElasticSearchProxy elastic_search_proxy.IElasticSearchProxy
}

func NewGetCircleParkingLocationService(IElasticSearchProxy elastic_search_proxy.IElasticSearchProxy) IGetCircleParkingLocationService {
	return &getCircleParkingLocationService{IElasticSearchProxy: IElasticSearchProxy}
}

func (sv *getCircleParkingLocationService) GetCircleParkingLocation(ctx context.Context, request api_model.GetCircleParkingLocationRequestDto) (
	parking []api_model.GetCircleParkingLocationModelResponseDto, errRes error) {

	parking = make([]api_model.GetCircleParkingLocationModelResponseDto, 0)
	var parkingList []elastic_search.ParkingModel
	var err error

	if request.Distance <= 0 {
		request.Distance = 5
	}
	if request.Distance > 100 {
		request.Distance = 100
	}

	parkingList, err = sv.IElasticSearchProxy.GetCircleParkingLocation(ctx, request)
	if err != nil {
		errRes = core.NewInternalError(err, constants.ERR_100001)
		return
	}

	for _, parkingModel := range parkingList {
		parkingTypes := make([]api_model.ParkingTypesResponseDto, 0)
		for _, types := range parkingModel.ParkingTypes {
			parkingTypes = append(parkingTypes, api_model.ParkingTypesResponseDto{
				Type:          types.Type,
				Status:        types.Status,
				Logo:          parking_type_car_enum.GetLogoFromData(types.Type),
				StatusDisplay: parking_types_car_status_enum.GetDisplayFromData(types.Status),
				TypeDisplay:   parking_type_car_enum.GetTypeNameFromData(types.Type),
			})
		}
		parking = append(parking, api_model.GetCircleParkingLocationModelResponseDto{
			ParkingId:     parkingModel.ParkingId,
			PublicId:      parkingModel.PublicId,
			Name:          parkingModel.Name,
			Status:        parkingModel.Status,
			StatusDisplay: parking_status_enum.GetEnumFromData(parkingModel.Status).Display,
			Distance:      parkingModel.Distance,
			Location: api_model.PointLocationApiDto{
				Lat: parkingModel.Location.Lat,
				Lng: parkingModel.Location.Lon,
			},
			ParkingTypes: parkingTypes,
			Address:      parkingModel.Address,
			ParkingPhone: parkingModel.ParkingPhone,
		})
	}

	return
}
