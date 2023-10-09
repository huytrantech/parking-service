package sync_on_elastic

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"parking-service/constants"
	"parking-service/core"
	"parking-service/core/enums/parking_status_enum"
	"parking-service/core/utils"
	"parking-service/model/database_model"
	"parking-service/provider/location_provider"
	"parking-service/proxy/elastic_search_proxy"
	"parking-service/repository"
)

type ISyncOnElasticService interface {
	SyncOnElastic(ctx context.Context, parkingId int) (
		errRes error)
	SyncAllOnElastic(ctx context.Context) (
		errRes error)
}

type syncOnElasticService struct {
	IElasticSearchProxy    elastic_search_proxy.IElasticSearchProxy
	IParkingRepository     repository.IParkingRepository
	IParkingSlotRepository repository.IParkingSlotRepository
	ILocationProvider      location_provider.ILocationProvider
}

func (sv *syncOnElasticService) SyncAllOnElastic(ctx context.Context) (errRes error) {

	dataParking, err := sv.IParkingRepository.FindMany(ctx, database_model.ParkingQueryModel{
		Fields: []string{"id"},
		Filter: database_model.ParkingFilterModel{Status: parking_status_enum.Active().Data},
		Limit:  1000})
	if err != nil {
		return core.NewInternalError(err, constants.ERR_100001)
	}

	for _, parking := range dataParking {
		err = sv.SyncOnElastic(ctx, parking.Id)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
	}

	return nil
}

func NewSyncOnElasticService(
	IElasticSearchProxy elastic_search_proxy.IElasticSearchProxy,
	IParkingRepository repository.IParkingRepository,
	IParkingSlotRepository repository.IParkingSlotRepository,
	ILocationProvider location_provider.ILocationProvider) ISyncOnElasticService {
	return &syncOnElasticService{
		IElasticSearchProxy:    IElasticSearchProxy,
		IParkingRepository:     IParkingRepository,
		IParkingSlotRepository: IParkingSlotRepository,
		ILocationProvider:      ILocationProvider}
}

func (sv *syncOnElasticService) SyncOnElastic(ctx context.Context, parkingId int) (
	errRes error) {

	if parkingId == 0 {
		errRes = core.NewBadRequestErrorMessage("Bad Request")
		return
	}
	var dataParking *database_model.ParkingModel
	var dataParkingSlot []database_model.ParkingSlotModel

	errG, ctxG := errgroup.WithContext(ctx)
	errG.Go(
		func() error {
			var err error
			dataParking, err = sv.IParkingRepository.FindOne(ctxG, database_model.ParkingQueryModel{
				Filter: database_model.ParkingFilterModel{
					Id: parkingId,
				},
				Fields: []string{"id", "public_id", "address", "parking_name",
					"parking_phone", "status", "open_at", "close_at"}})
			if err != nil {
				err = core.NewInternalError(err, constants.ERR_100001)
				return err
			}

			if dataParking == nil {
				err = core.NewBadRequestErrorMessage("Not found parking")
				return err
			}
			return nil
		})

	errG.Go(func() error {
		var err error
		dataParkingSlot, err = sv.IParkingSlotRepository.
			FindManyParkingSlotWithParkingById(ctxG, []string{"parking_type", "status", "total_slot", "current_slot", "price"}, parkingId)

		if err != nil {
			err = core.NewInternalError(err, constants.ERR_100001)
			return err
		}

		return nil
	})

	if err := errG.Wait(); err != nil {
		return err
	}

	modelParkingES := dataParking.ConvertToModelES(dataParkingSlot)
	modelParkingES.Address = utils.GetAddressString(dataParking.Address, sv.ILocationProvider.GetMapData())

	err := sv.IElasticSearchProxy.CreateParking(ctx, modelParkingES)

	if err != nil {
		errRes = core.NewInternalError(err, constants.ERR_100001)
		return
	}

	return
}
