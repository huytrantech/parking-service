package update_parking_slot

import (
	"context"
	"parking-service/constants"
	"parking-service/core"
	"parking-service/model/api_model"
	"parking-service/repository"
)

type IUpdateParkingSlotService interface {
	UpdateParkingSlot(ctx context.Context ,
		request api_model.UpdateParkingSlotRequest) error
}

type updateParkingSlotService struct {
	IParkingSlotRepository repository.IParkingSlotRepository
}

func NewUpdateParkingSlotService(IParkingSlotRepository repository.IParkingSlotRepository) IUpdateParkingSlotService {
	return &updateParkingSlotService{IParkingSlotRepository: IParkingSlotRepository}
}

func (sv *updateParkingSlotService) UpdateParkingSlot(ctx context.Context ,
	request api_model.UpdateParkingSlotRequest) error {

	selectOption := core.SelectOption{
		Projection:    []string{"id"},
		Table:         "parking_slot",
		Filter: map[string]interface{}{
			"id": request.ParkingSlotId,
			"parking_id": request.ParkingId,
		},
	}
	selectOption.BuildQuery()
	queryStr , dataParam := selectOption.GetFinalQuery()
	dataParkingSlot , err := sv.IParkingSlotRepository.FindOneParkingSlot(ctx , queryStr , dataParam)
	if err != nil {
		return core.NewInternalError(err , constants.ERR_100001)
	}
	if dataParkingSlot == nil {
		return core.NewBadRequestErrorMessage("Not found parking slot")
	}

	updateOption := core.UpdateOption{
		Filter: map[string]interface{}{
			"id": request.ParkingSlotId,
		},
		Updated: map[string]interface{}{
			"total_slot":request.TotalSlot,
			"price": request.Price,
		},
		Table: "parking_slot",
	}

	updateOption.BuildQuery()
	queryStr, dataParam = updateOption.GetFinalQuery()
	err = sv.IParkingSlotRepository.UpdateParkingSlot(ctx , queryStr , dataParam)
	if err != nil {
		return core.NewInternalError(err , constants.ERR_100001)
	}
	return nil
}