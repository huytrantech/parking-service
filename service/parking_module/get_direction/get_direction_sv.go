package get_direction

import (
	"context"
	"encoding/json"
	"parking-service/constants"
	"parking-service/core"
	"parking-service/core/utils"
	"parking-service/model/api_model"
	"parking-service/model/database_model"
	"parking-service/provider/redis_provider"
	"parking-service/proxy/google_map_proxy"
	"parking-service/repository"
	"time"
)

type IGetDirectionService interface {
	GetDirection(ctx context.Context,
		request api_model.GetDirectionRequestDto) (
		response api_model.GetDirectionResponseDto, errRes error)
}

type getDirectionService struct {
	IGoogleMapProxy    google_map_proxy.IGoogleMapProxy
	IRedisProvider     redis_provider.IRedisProvider
	IParkingRepository repository.IParkingRepository
}

func NewGetDirectionService(IGoogleMapProxy google_map_proxy.IGoogleMapProxy,
	IRedisProvider redis_provider.IRedisProvider,
	IParkingRepository repository.IParkingRepository) IGetDirectionService {
	return &getDirectionService{
		IGoogleMapProxy:    IGoogleMapProxy,
		IRedisProvider:     IRedisProvider,
		IParkingRepository: IParkingRepository,
	}
}

func (sv *getDirectionService) GetDirection(ctx context.Context,
	request api_model.GetDirectionRequestDto) (
	response api_model.GetDirectionResponseDto, errRes error) {

	response.Points = make([]api_model.PointLocationApiDto, 0)
	dataPointGoogle := make([]google_map_proxy.PointLocation, 0)
	key := utils.GetKeyRedis(request.Origin, request.ParkingId)
	valueRedis, err := sv.IRedisProvider.GetKey(ctx, key)
	if len(valueRedis) > 0 && err == nil {
		json.Unmarshal([]byte(valueRedis), &response)
		if len(response.Points) > 0 {
			return
		}
	}

	parking, err := sv.IParkingRepository.FindOne(ctx, database_model.ParkingQueryModel{
		Filter: database_model.ParkingFilterModel{PublicId: request.ParkingId},
		Fields: []string{"address"},
	})
	if err != nil {
		errRes = core.NewInternalError(err, constants.ERR_100001)
		return
	}

	if parking == nil {
		errRes = core.NewBadRequestErrorMessage("Not found parking")
		return
	}

	dataPointGoogle, err = sv.IGoogleMapProxy.GetDirection(ctx, google_map_proxy.PointLocation{
		Lat: request.Origin.Lat,
		Lng: request.Origin.Lng,
	}, google_map_proxy.PointLocation{
		Lat: parking.Address.Lat,
		Lng: parking.Address.Lng,
	})
	if err != nil {
		errRes = err
		return
	}
	for _, value := range dataPointGoogle {
		response.Points = append(response.Points, api_model.PointLocationApiDto{
			Lat: value.Lat,
			Lng: value.Lng,
		})
	}

	_ = sv.IRedisProvider.SetKey(ctx, key, response, time.Minute*60)
	return
}
