package reopen_parking

import (
	"context"
	"parking-service/constants"
	"parking-service/core"
	"parking-service/core/enums/parking_status_enum"
	"parking-service/core/utils"
	"parking-service/model/database_model"
	"parking-service/repository"
)

type IReopenParkingService interface {
	ReopenParking(ctx context.Context, parkingId string, username string) error
}

type reopenParkingService struct {
	IParkingRepository repository.IParkingRepository
}

func NewReopenParkingService(IParkingRepository repository.IParkingRepository) IReopenParkingService {
	return &reopenParkingService{IParkingRepository: IParkingRepository}
}

func (sv *reopenParkingService) ReopenParking(ctx context.Context, parkingId string, username string) error {

	if len(parkingId) == 0 {
		return core.NewBadRequestErrorMessage("Bad Request")
	}

	dataParking, err := sv.IParkingRepository.FindOne(ctx, database_model.ParkingQueryModel{
		Fields: []string{"status", "public_id"},
		Filter: database_model.ParkingFilterModel{PublicId: parkingId}})
	if err != nil {
		return core.NewInternalError(err, constants.ERR_100001)
	}

	if dataParking == nil {
		return core.NewBadRequestErrorMessage("Not found parking")
	}

	if !utils.CanReopen(*dataParking) {
		return core.NewBadRequestErrorMessage("Cannot reopen parking")
	}

	status := parking_status_enum.Pending().Data
	err = sv.IParkingRepository.UpdateOne(ctx, database_model.ParkingQueryModel{
		Filter:     database_model.ParkingFilterModel{PublicId: parkingId},
		UpdateUser: username,
		Update:     database_model.ParkingUpdateModel{Status: &status},
	})
	if err != nil {
		return core.NewInternalError(err, constants.ERR_100001)
	}

	return nil
}
