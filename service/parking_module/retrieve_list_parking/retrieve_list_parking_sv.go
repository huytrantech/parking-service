package retrieve_list_parking

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"parking-service/constants"
	"parking-service/core"
	"parking-service/core/enums/parking_status_enum"
	"parking-service/core/utils"
	"parking-service/model/api_model"
	"parking-service/model/database_model"
	"parking-service/repository"
)

type IRetrieveListParkingService interface {
	RetrieveListParkingDatabase(ctx context.Context, request api_model.RetrieveListParkingRequestDto) (
		response *api_model.RetrieveListParkingResponseBaseDto, err error)
}

type retrieveListParkingService struct {
	IParkingRepository repository.IParkingRepository
}

func NewRetrieveListParkingService(IParkingRepository repository.IParkingRepository) IRetrieveListParkingService {
	return &retrieveListParkingService{IParkingRepository: IParkingRepository}
}

func (sv *retrieveListParkingService) RetrieveListParkingDatabase(ctx context.Context, request api_model.RetrieveListParkingRequestDto) (
	response *api_model.RetrieveListParkingResponseBaseDto, err error) {

	if request.PageIndex <= 0 {
		request.PageIndex = 1
	}
	if request.PageLimit <= 0 {
		request.PageLimit = 10
	}

	queryDto := database_model.ParkingQueryModel{
		Limit:  request.PageLimit,
		Offset: int64((request.PageIndex - 1) * request.PageLimit),
		Filter: database_model.ParkingFilterModel{},
	}

	var dataParking []database_model.ParkingModel
	var total int

	errG, ctxG := errgroup.WithContext(ctx)

	errG.Go(func() error {
		var errGroup error
		dataParking, errGroup = sv.IParkingRepository.FindMany(ctxG, queryDto)
		if errGroup != nil {
			return errors.New(fmt.Sprintf("IParkingRepository.FindMany with error %s", errGroup.Error()))
		}
		return nil
	})

	errG.Go(func() error {
		var errGroup error
		total, errGroup = sv.IParkingRepository.CountParking(ctxG, queryDto)
		if errGroup != nil {
			return errors.New(fmt.Sprintf("IParkingRepository.CountParking with error %s", errGroup.Error()))
		}
		return nil
	})

	if err = errG.Wait(); err != nil {
		err = core.NewInternalError(err, constants.ERR_100001)
		return
	}

	result := make([]api_model.RetrieveListParkingResponseDto, len(dataParking))
	for index, parking := range dataParking {
		result[index] = api_model.RetrieveListParkingResponseDto{
			PublicId:      parking.PublicId,
			ParkingName:   parking.ParkingName,
			ParkingPhone:  parking.ParkingPhone,
			Status:        parking.Status,
			StatusDisplay: parking_status_enum.GetEnumFromData(parking.Status).Display,
			Roles: api_model.RetrieveListParkingActionRolesDto{
				CanApprove: utils.CanApprove(parking),
				CanRemove:  utils.CanRemove(parking),
				CanDenied:  utils.CanDenied(parking),
				CanBlock:   utils.CanBlock(parking),
				CanClose:   utils.CanClose(parking),
				CanReopen:  utils.CanReopen(parking),
			},
		}
	}
	response = new(api_model.RetrieveListParkingResponseBaseDto)
	response.PageLimit = request.PageLimit
	response.PageIndex = request.PageIndex
	response.IsLastPage = len(dataParking)+int((request.PageIndex-1)*request.PageLimit) == total
	response.Total = total
	response.Parks = result
	return
}
