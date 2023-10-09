package create_parking

import (
	"context"
	"github.com/spf13/cast"
	"parking-service/constants"
	"parking-service/core"
	"parking-service/core/enums/parking_status_enum"
	"parking-service/model/api_model"
	"parking-service/model/database_model"
	"parking-service/repository"
	"parking-service/service/parking_module/gen_public_id"
	"strings"
	"time"
)

type ICreateParkingService interface {
	CreateParking(ctx context.Context, request api_model.CreateParkingRequestDto) (
		publicId string, errRes error)
}
type createParkingService struct {
	IParkingRepository     repository.IParkingRepository
	IParkingSlotRepository repository.IParkingSlotRepository
	IGenPublicIdService    gen_public_id.IGenPublicIdService
}

func NewCreateParkingService(
	IParkingRepository repository.IParkingRepository,
	IParkingSlotRepository repository.IParkingSlotRepository,
	IGenPublicIdService gen_public_id.IGenPublicIdService) ICreateParkingService {
	return &createParkingService{
		IParkingRepository:     IParkingRepository,
		IParkingSlotRepository: IParkingSlotRepository,
		IGenPublicIdService:    IGenPublicIdService,
	}
}

func (sv *createParkingService) CreateParking(ctx context.Context, request api_model.CreateParkingRequestDto) (
	publicId string, errRes error) {

	if err := request.Invalid(); err != nil {
		errRes = core.NewBadRequestErrorMessage(err.Error())
		return
	}
	now := time.Now()
	parkingTypes := make([]string, 0)
	for _, value := range request.ParkingTypes {
		parkingTypes = append(parkingTypes, cast.ToString(value.ParkingType))
	}
	publicId, errRes = sv.IGenPublicIdService.GenPublicId(ctx)
	if errRes != nil {
		return
	}

	modelAddParking := database_model.ParkingModel{
		OwnerName: request.OwnerName,
		Address: database_model.ParkingAddressModel{
			CityId:     request.CityId,
			DistrictId: request.DistrictId,
			WardId:     request.WardId,
			Address:    request.Address,
			Lat:        request.Lat,
			Lng:        request.Lng,
		},
		PublicId:     publicId,
		Status:       parking_status_enum.Pending().Data,
		OwnerPhone:   request.OwnerPhone,
		ParkingName:  request.ParkingName,
		ParkingPhone: request.ParkingPhone,
	}

	modelAddParking.CreatedDate = now
	modelAddParking.CreatedName = request.Username

	openAt := time.Date(0, 0, 0, request.OpenAtHour, request.OpenAtMinute, 0, 0, time.Local)
	closeAt := time.Date(0, 0, 0, request.CloseAtHour, request.CloseAtMinute, 59, 0, time.Local)
	modelAddParking.OpenAt = &openAt
	modelAddParking.CloseAt = &closeAt

	if len(request.Images) > 0 {
		images := strings.Join(request.Images, ",")
		modelAddParking.Images = &images
	}
	if len(parkingTypes) > 0 {
		parkingTypesStr := strings.Join(parkingTypes, ",")
		modelAddParking.ParkingTypes = &parkingTypesStr
	}
	_, err := sv.IParkingRepository.CreateParking(ctx, modelAddParking)
	if err != nil {
		errRes = core.NewInternalError(err, constants.ERR_100001)
		return
	}

	//if len(request.ParkingTypes) > 0 {
	//	arrParkingSlotModel := make([]database_model.ParkingSlotModel, 0)
	//	for _, parkingType := range request.ParkingTypes {
	//		modelAddSlot := database_model.ParkingSlotModel{
	//			ParkingId:   parkingID,
	//			ParkingType: parkingType.ParkingType,
	//			Price:       parkingType.Price,
	//			TotalSlot:   parkingType.TotalSlot,
	//			Status:      parking_types_car_status_enum.Available().Data,
	//		}
	//		arrParkingSlotModel = append(arrParkingSlotModel, modelAddSlot)
	//	}
	//
	//	err = sv.IParkingSlotRepository.InsertManyParkingSlot(ctx, arrParkingSlotModel)
	//	if err != nil {
	//		errRes = core.NewInternalError(err, constants.ERR_100001)
	//		return
	//	}
	//}

	return
}
