package update_parking

import (
	"context"
	"parking-service/constants"
	"parking-service/core"
	"parking-service/model/api_model"
	"parking-service/model/database_model"
	"parking-service/proxy/elastic_search_proxy"
	"parking-service/repository"
)

type IUpdateParkingService interface {
	UpdateParking(ctx context.Context,
		request api_model.UpdateParkingRequestDto, parkingId string) error
}

type updateParkingService struct {
	IParkingRepository  repository.IParkingRepository
	IElasticSearchProxy elastic_search_proxy.IElasticSearchProxy
}

func NewUpdateParkingService(IParkingRepository repository.IParkingRepository,
	IElasticSearchProxy elastic_search_proxy.IElasticSearchProxy) IUpdateParkingService {
	return &updateParkingService{
		IElasticSearchProxy: IElasticSearchProxy,
		IParkingRepository:  IParkingRepository,
	}
}

func (sv *updateParkingService) UpdateParking(ctx context.Context,
	request api_model.UpdateParkingRequestDto, parkingId string) error {

	if len(parkingId) == 0 {
		return core.NewBadRequestErrorMessage("Thông tin không hợp lệ")
	}

	parkingModel, err := sv.IParkingRepository.FindOne(ctx, database_model.ParkingQueryModel{
		Fields: []string{"public_id"}, Filter: database_model.ParkingFilterModel{PublicId: parkingId},
	})
	if err != nil {
		return core.NewInternalError(err, constants.ERR_100001)
	}

	if parkingModel == nil {
		return core.NewBadRequestErrorMessage("Không tìm thấy thông tin parking")
	}

	updateDto := database_model.ParkingUpdateModel{}

	if len(request.ParkingName) > 0 {
		updateDto.ParkingName = &request.ParkingName
	}

	if len(request.ParkingPhone) > 0 {
		updateDto.ParkingPhone = &request.ParkingPhone
	}

	if len(request.OwnerName) > 0 {
		updateDto.OwnerName = &request.OwnerName
	}

	if len(request.OwnerPhone) > 0 {
		updateDto.OwnerPhone = &request.OwnerPhone
	}

	if len(request.Address) > 0 {
		updateDto.Address = &request.Address
	}
	if request.Lng > 0 {
		updateDto.Lng = &request.Lng
	}

	if request.Lat > 0 {
		updateDto.Lat = &request.Lat
	}

	if request.CityId > 0 {
		updateDto.CityId = &request.CityId
	}

	if request.DistrictId > 0 {
		updateDto.DistrictId = &request.DistrictId
	}

	if request.WardId > 0 {
		updateDto.WardId = &request.WardId
	}

	if request.OpenAt != nil {
		updateDto.OpenAt = request.OpenAt
	}
	if request.CloseAt != nil {
		updateDto.CloseAt = request.CloseAt
	}

	err = sv.IParkingRepository.UpdateOne(ctx, database_model.ParkingQueryModel{
		UpdateUser: request.Username,
		Filter:     database_model.ParkingFilterModel{PublicId: parkingId},
		Update:     updateDto,
	})
	if err != nil {
		return core.NewInternalError(err, constants.ERR_100001)
	}
	return nil
}
