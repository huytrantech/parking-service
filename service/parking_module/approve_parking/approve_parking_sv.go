package approve_parking

import (
	"context"
	"parking-service/constants"
	"parking-service/core"
	"parking-service/core/enums/parking_status_enum"
	"parking-service/core/utils"
	"parking-service/model/database_model"
	"parking-service/provider/location_provider"
	"parking-service/proxy/elastic_search_proxy"
	"parking-service/repository"
	"parking-service/service/parking_module/sync_on_elastic"
)

type IApproveParkingService interface {
	ApproveParking(ctx context.Context, parkingId string, username string) (
		errRes error)
}

type approveParkingService struct {
	IParkingRepository    repository.IParkingRepository
	IElasticSearchProxy   elastic_search_proxy.IElasticSearchProxy
	ISyncOnElasticService sync_on_elastic.ISyncOnElasticService
	ILocationProvider     location_provider.ILocationProvider
}

func NewApproveParkingService(
	IParkingRepository repository.IParkingRepository,
	IElasticSearchProxy elastic_search_proxy.IElasticSearchProxy,
	ISyncOnElasticService sync_on_elastic.ISyncOnElasticService,
	ILocationProvider location_provider.ILocationProvider) IApproveParkingService {
	return &approveParkingService{
		IParkingRepository:    IParkingRepository,
		IElasticSearchProxy:   IElasticSearchProxy,
		ISyncOnElasticService: ISyncOnElasticService,
		ILocationProvider:     ILocationProvider,
	}
}

func (sv *approveParkingService) ApproveParking(ctx context.Context, parkingId string, username string) (
	errRes error) {

	if len(parkingId) == 0 {
		errRes = core.NewBadRequestErrorMessage("Bad Request")
		return
	}

	dataParking, err := sv.IParkingRepository.FindOne(ctx, database_model.ParkingQueryModel{
		Fields: []string{"public_id", "status"},
		Filter: database_model.ParkingFilterModel{PublicId: parkingId}})
	if err != nil {
		errRes = core.NewInternalError(err, constants.ERR_100001)
		return
	}

	if dataParking == nil {
		errRes = core.NewBadRequestErrorMessage("Not found parking")
		return
	}

	if !utils.CanApprove(*dataParking) {
		errRes = core.NewBadRequestErrorMessage("Cannot approved parking")
		return
	}

	status := parking_status_enum.Active().Data
	err = sv.IParkingRepository.UpdateOne(ctx, database_model.ParkingQueryModel{
		Filter:     database_model.ParkingFilterModel{PublicId: parkingId},
		UpdateUser: username,
		Update:     database_model.ParkingUpdateModel{Status: &status}})
	if err != nil {
		errRes = core.NewInternalError(err, constants.ERR_100001)
		return
	}

	//err = sv.ISyncOnElasticService.SyncOnElastic(ctx, parkingId)
	//
	//if err != nil {
	//	errRes = core.NewInternalError(err, constants.ERR_100001)
	//	return
	//}

	return
}
