package close_parking

import (
	"context"
	"parking-service/constants"
	"parking-service/core"
	"parking-service/core/enums/parking_status_enum"
	"parking-service/core/utils"
	"parking-service/model/database_model"
	"parking-service/proxy/elastic_search_proxy"
	"parking-service/repository"
)

type ICloseParkingService interface {
	CloseParking(ctx context.Context, parkingId string, username string) error
}

type closeParkingService struct {
	IParkingRepository  repository.IParkingRepository
	IElasticSearchProxy elastic_search_proxy.IElasticSearchProxy
}

func NewCloseParkingService(IParkingRepository repository.IParkingRepository,
	IElasticSearchProxy elastic_search_proxy.IElasticSearchProxy) ICloseParkingService {
	return &closeParkingService{
		IParkingRepository:  IParkingRepository,
		IElasticSearchProxy: IElasticSearchProxy,
	}
}

func (sv *closeParkingService) CloseParking(ctx context.Context, parkingId string, username string) error {

	if len(parkingId) <= 0 {
		return core.NewBadRequestErrorMessage("thông tin không hợp lệ")
	}

	parking, err := sv.IParkingRepository.FindOne(ctx, database_model.ParkingQueryModel{
		Fields: []string{"status", "public_id"},
		Filter: database_model.ParkingFilterModel{PublicId: parkingId},
	})
	if err != nil {
		return core.NewInternalError(err, constants.ERR_100001)
	}

	if parking == nil {
		return core.NewBadRequestErrorMessage("Không tìm thấy thông tin parking")
	}

	if utils.CanClose(*parking) == false {
		return core.NewBadRequestErrorMessage("Không thể close parking này")
	}

	status := parking_status_enum.Closed().Data
	err = sv.IParkingRepository.UpdateOne(ctx, database_model.ParkingQueryModel{
		UpdateUser: username,
		Update: database_model.ParkingUpdateModel{
			Status: &status,
		},
		Filter: database_model.ParkingFilterModel{PublicId: parkingId},
	})
	if err != nil {
		return core.NewInternalError(err, constants.ERR_100001)
	}

	//err = sv.IElasticSearchProxy.RemoveParkingByQuery(ctx, map[string]interface{}{
	//	"query": map[string]interface{}{
	//		"term": map[string]interface{}{
	//			"parking_id": parkingId,
	//		},
	//	},
	//})
	if err != nil {
		return core.NewInternalError(err, constants.ERR_100001)
	}
	return nil
}
