package detail_parking

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"parking-service/core/enums/parking_status_enum"
	"parking-service/core/enums/parking_type_car_enum"
	"parking-service/core/enums/parking_types_car_status_enum"
	"parking-service/core/utils"
	"parking-service/model/api_model"
	"parking-service/model/database_model"
	"parking-service/provider/location_provider"
	"parking-service/repository"
	"time"
)

type IDetailParkingService interface {
	GetDetailParking(ctx context.Context, parkingId string) (
		detailParking *api_model.DetailParkingResponseDto, err error)
}
type detailParkingService struct {
	IParkingRepository     repository.IParkingRepository
	IParkingSlotRepository repository.IParkingSlotRepository
	ILocationProvider      location_provider.ILocationProvider
}

func NewDetailParkingService(
	IParkingRepository repository.IParkingRepository,
	IParkingSlotRepository repository.IParkingSlotRepository,
	ILocationProvider location_provider.ILocationProvider) IDetailParkingService {
	return &detailParkingService{
		IParkingRepository:     IParkingRepository,
		IParkingSlotRepository: IParkingSlotRepository,
		ILocationProvider:      ILocationProvider}
}

func (sv *detailParkingService) GetDetailParking(ctx context.Context, parkingId string) (
	detailParking *api_model.DetailParkingResponseDto, err error) {

	var dataParking *database_model.ParkingModel
	var dataParkingSlot []database_model.ParkingSlotModel

	errG, ctxG := errgroup.WithContext(ctx)
	errG.Go(
		func() error {
			var errGroup error
			dataParking, errGroup = sv.IParkingRepository.FindOne(ctxG, database_model.ParkingQueryModel{
				Filter: database_model.ParkingFilterModel{
					PublicId: parkingId,
				},
				Fields: []string{"id", "address", "parking_name", "parking_phone",
					"status", "open_at", "close_at", "owner_name", "owner_phone", "public_id"},
			})
			if errGroup != nil {
				errGroup = errors.New(fmt.Sprintf("IParkingRepository.FindOne With Error %s", errGroup.Error()))
				return err
			}

			if dataParking == nil {
				return errors.New("IParkingRepository.FindOne not found")
			}
			return nil
		})

	//errG.Go(func() error {
	//	var errGroup error
	//	dataParkingSlot, errGroup = sv.IParkingSlotRepository.
	//		FindManyParkingSlotWithParkingById(ctxG, []string{"id", "parking_type", "status", "total_slot", "current_slot", "price"}, parkingId)
	//
	//	if errGroup != nil {
	//		errGroup = errors.New(fmt.Sprintf("IParkingRepository.FindOne With Error %s" , errGroup.Error()))
	//		return err
	//	}
	//
	//	return nil
	//})

	if err = errG.Wait(); err != nil {
		return nil, err
	}

	detailParking = &api_model.DetailParkingResponseDto{
		ParkingId:     dataParking.PublicId,
		ParkingName:   dataParking.ParkingName,
		ParkingPhone:  dataParking.ParkingPhone,
		Status:        dataParking.Status,
		StatusDisplay: parking_status_enum.GetEnumFromData(dataParking.Status).Display,
		FullAddress:   utils.GetAddressString(dataParking.Address, sv.ILocationProvider.GetMapData()),
		Address:       dataParking.Address.Address,
		Location: api_model.PointLocationApiDto{
			Lat: dataParking.Address.Lat,
			Lng: dataParking.Address.Lng,
		},
		CityId:     dataParking.Address.CityId,
		DistrictId: dataParking.Address.DistrictId,
		WardId:     dataParking.Address.WardId,
		OwnerPhone: dataParking.OwnerPhone,
		OwnerName:  dataParking.OwnerName,
	}
	if dataParking.CloseAt == nil {
		closeAt := time.Date(0, 0, 0, 23, 59, 59, 0, time.Local)
		dataParking.CloseAt = &closeAt
	}
	if dataParking.OpenAt == nil {
		openAt := time.Date(0, 0, 0, 0, 0, 0, 0, time.Local)
		dataParking.OpenAt = &openAt
	}
	detailParking.OpenAtHour = dataParking.OpenAt.Hour()
	detailParking.OpenAtMinute = dataParking.OpenAt.Minute()
	detailParking.CloseAtHour = dataParking.CloseAt.Hour()
	detailParking.CloseAtMinute = dataParking.CloseAt.Minute()
	detailParking.ParkingTypes = make([]api_model.ParkingTypesResponseDto, 0)

	for _, v := range dataParkingSlot {
		detailParking.ParkingTypes = append(detailParking.ParkingTypes, api_model.ParkingTypesResponseDto{
			Id:            v.Id,
			Type:          v.ParkingType,
			Status:        v.Status,
			Logo:          parking_type_car_enum.GetLogoFromData(v.ParkingType),
			StatusDisplay: parking_types_car_status_enum.GetDisplayFromData(v.Status),
			TypeDisplay:   parking_type_car_enum.GetTypeNameFromData(v.ParkingType),
			TotalSlot:     v.TotalSlot,
			CurrentSlot:   v.CurrentSlot,
			Price:         v.Price,
		})
	}

	return
}
