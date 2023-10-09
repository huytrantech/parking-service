package import_parking_csv

import (
	"context"
	"encoding/json"
	"github.com/spf13/cast"
	"parking-service/model/api_model"
	"parking-service/model/database_model"
	"parking-service/repository"
	"parking-service/service/parking_module/gen_public_id"
	"time"
)

type IImportParkingCSVService interface {
	ImportParkingCSV(ctx context.Context, dataCsv [][]string, userName string) (response api_model.ImportParkingResponseDto)
}

type importParkingCsvService struct {
	iParkingRepository  repository.IParkingRepository
	iGenPublicIdService gen_public_id.IGenPublicIdService
}

func (sv *importParkingCsvService) ImportParkingCSV(ctx context.Context, dataCsv [][]string, userName string) (response api_model.ImportParkingResponseDto) {

	mapCol := map[string]int{
		"OwnerName":    0,
		"ParkingName":  1,
		"OwnerPhone":   2,
		"ParkingPhone": 3,
		"Address":      4,
		"Status":       5,
		"OpenAt":       6,
		"CloseAt":      7,
	}
	now := time.Now()
	dataCsv = dataCsv[1:]
	for _, dataRows := range dataCsv {
		publicId, err := sv.iGenPublicIdService.GenPublicId(ctx)
		if err != nil {
			response.TotalFail += 1
			continue
		}
		var addressData database_model.ParkingAddressModel
		err = json.Unmarshal([]byte(dataRows[mapCol["Address"]]), &addressData)
		if err != nil {
			response.TotalFail += 1
			continue
		}
		model := database_model.ParkingModel{
			BaseEntity: database_model.BaseEntity{
				CreatedName: userName,
				CreatedDate: now,
			},
			PublicId:     publicId,
			Address:      addressData,
			Status:       cast.ToInt(dataRows[mapCol["Status"]]),
			OwnerPhone:   dataRows[mapCol["OwnerPhone"]],
			ParkingName:  dataRows[mapCol["ParkingName"]],
			OwnerName:    dataRows[mapCol["OwnerName"]],
			ParkingPhone: dataRows[mapCol["ParkingPhone"]],
		}

		openAtTime, err := time.Parse(time.TimeOnly, dataRows[mapCol["OpenAt"]])
		if err == nil && openAtTime.IsZero() == false {
			model.OpenAt = &openAtTime
		}

		closeAtTime, err := time.Parse(time.TimeOnly, dataRows[mapCol["CloseAt"]])
		if err == nil && closeAtTime.IsZero() == false {
			model.CloseAt = &closeAtTime
		}

		_, err = sv.iParkingRepository.CreateParking(ctx, model)
		if err != nil {
			response.TotalFail += 1
		} else {
			response.TotalSuccess += 1
		}
	}

	return
}

func NewImportParkingCSVService(iParkingRepository repository.IParkingRepository,
	iGenPublicIdService gen_public_id.IGenPublicIdService) IImportParkingCSVService {
	return &importParkingCsvService{iParkingRepository: iParkingRepository, iGenPublicIdService: iGenPublicIdService}
}
