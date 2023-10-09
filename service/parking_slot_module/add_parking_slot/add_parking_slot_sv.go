package add_parking_slot

import (
	"context"
	"parking-service/constants"
	"parking-service/core"
	"parking-service/model/api_model"
	"parking-service/model/database_model"
	"parking-service/repository"
)

type IAddParkingSlotService interface {
	AddParkingSlot(ctx context.Context, request api_model.AddParkingSlotRequest) error
}

type addParkingSlotService struct {
	IParkingSlotRepository repository.IParkingSlotRepository
}

func NewAddParkingSlotService(IParkingSlotRepository repository.IParkingSlotRepository) IAddParkingSlotService {
	return &addParkingSlotService{IParkingSlotRepository: IParkingSlotRepository}
}

func(sv *addParkingSlotService) AddParkingSlot(ctx context.Context, request api_model.AddParkingSlotRequest) error {

	if err := request.Invalid(); err != nil {
		return core.NewBadRequestError(err)
	}

	selectOption := core.SelectOption{
		Projection:    []string{"parking_type"},
		Table:         "parking_slot",
		Filter: map[string]interface{}{
			"parking_type": request.Type,
			"parking_id": request.ParkingId,
		},
	}
	selectOption.BuildQuery()
	queryStr , dataParam := selectOption.GetFinalQuery()
	dataParkingSlot , err := sv.IParkingSlotRepository.FindOneParkingSlot(ctx , queryStr , dataParam)
	if err != nil {
		return core.NewInternalError(err , constants.ERR_100001)
	}

	if dataParkingSlot != nil {
		return core.NewBadRequestErrorMessage("Parking Slot is exist")
	}

	_,err = sv.IParkingSlotRepository.InsertOneParkingSlot(ctx , database_model.ParkingSlotModel{
		ParkingId:   request.ParkingId,
		ParkingType: request.Type,
		Price:       request.Price,
		TotalSlot:   &request.TotalSlot,
	})
	if err != nil {
		return core.NewInternalError(err , constants.ERR_100001)
	}
	return nil
}