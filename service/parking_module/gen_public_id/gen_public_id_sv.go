package gen_public_id

import (
	"context"
	"parking-service/constants"
	"parking-service/core"
	"parking-service/core/utils"
	"parking-service/model/database_model"
	"parking-service/repository"
)

type IGenPublicIdService interface {
	GenPublicId(ctx context.Context) (string, error)
}

type genPublicIdService struct {
	iParkingRepository repository.IParkingRepository
}

func (sv *genPublicIdService) GenPublicId(ctx context.Context) (publicId string, errRes error) {
	counterGen := 5
	for true {
		if counterGen == 0 {
			errRes = core.NewBadRequestErrorMessage("gen public id fail")
			return
		}
		publicId = utils.GenParkingId()
		parkingData, err := sv.iParkingRepository.FindOne(ctx, database_model.ParkingQueryModel{
			Fields: []string{"public_id"}, Filter: database_model.ParkingFilterModel{
				PublicId: publicId,
			},
		})
		if err != nil {
			errRes = core.NewInternalError(err, constants.ERR_100001)
			return
		}
		if parkingData == nil {
			return
		}
		counterGen -= 1
	}
	return
}

func NewGenPublicIdService(iParkingRepository repository.IParkingRepository) IGenPublicIdService {
	return &genPublicIdService{iParkingRepository: iParkingRepository}
}
