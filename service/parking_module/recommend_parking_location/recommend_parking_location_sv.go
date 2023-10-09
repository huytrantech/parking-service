package recommend_parking_location

import (
	"context"
	"parking-service/constants"
	"parking-service/core"
	"parking-service/model/api_model"
	"parking-service/proxy/elastic_search_proxy"
	"parking-service/repository"
)

type IRecommendParkingLocationService interface {
	RecommendParkingLocation(ctx context.Context, request api_model.RecommendParkingLocationRequestDto) (
		data api_model.RecommendParkingLocationBaseResponseDto, err error)
}

type recommendParkingLocationService struct {
	IParkingRepository  repository.IParkingRepository
	IElasticSearchProxy elastic_search_proxy.IElasticSearchProxy
}

func NewRecommendParkingLocationService(IParkingRepository repository.IParkingRepository,
	IElasticSearchProxy elastic_search_proxy.IElasticSearchProxy) IRecommendParkingLocationService {
	return &recommendParkingLocationService{
		IParkingRepository:  IParkingRepository,
		IElasticSearchProxy: IElasticSearchProxy,
	}
}

func (sv *recommendParkingLocationService) RecommendParkingLocation(ctx context.Context, request api_model.RecommendParkingLocationRequestDto) (
	data api_model.RecommendParkingLocationBaseResponseDto, err error) {

	dataEs, errQuery := sv.IElasticSearchProxy.GetDataElasticSearchCustomQuery(ctx, request.ConvertQueryES())
	if errQuery != nil {
		err = core.NewInternalError(errQuery, constants.ERR_100001)
		return
	}

	for _, value := range dataEs {
		data.Parking = append(data.Parking, api_model.RecommendParkingLocationResponseDto{
			ParkingName: value.Name,
			ParkingId:   value.PublicId,
			Location: api_model.PointLocationApiDto{
				Lat: value.Location.Lat,
				Lng: value.Location.Lon,
			},
		})
	}

	return
}
