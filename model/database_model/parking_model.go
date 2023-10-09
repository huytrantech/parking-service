package database_model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"parking-service/core/utils/time_utils"
	"parking-service/model/proxy_model/elastic_search"
	"strings"
	"time"
)

type ParkingModel struct {
	BaseEntity
	PublicId     string              `db:"public_id"`
	Address      ParkingAddressModel `db:"address"`
	Status       int                 `db:"status"`
	OwnerPhone   string              `db:"owner_phone"`
	ParkingName  string              `db:"parking_name"`
	OwnerName    string              `db:"owner_name"`
	ParkingPhone string              `db:"parking_phone"`
	OpenAt       *time.Time          `db:"open_at"`
	CloseAt      *time.Time          `db:"close_at"`
	Images       *string             `db:"images"`
	ParkingTypes *string             `db:"parking_types"`
}

func (model ParkingModel) ConvertToModelES(modelSlot []ParkingSlotModel) elastic_search.ParkingModel {

	mapParkingSlot := make(map[int]ParkingSlotModel)
	for _, value := range modelSlot {
		mapParkingSlot[value.ParkingType] = value
	}

	modelSyncES := elastic_search.ParkingModel{
		ParkingId: model.Id,
		PublicId:  model.PublicId,
		Location: elastic_search.Location{
			Lat: model.Address.Lat,
			Lon: model.Address.Lng,
		},
		Name:         model.ParkingName,
		CityId:       model.Address.CityId,
		DistrictId:   model.Address.DistrictId,
		WardId:       model.Address.WardId,
		Status:       model.Status,
		ParkingPhone: model.ParkingPhone,
		OpenAt:       time_utils.ConvertFromTimeToSecond(model.OpenAt),
		CloseAt:      time_utils.ConvertFromTimeToSecond(model.CloseAt),
	}
	for _, v := range modelSlot {
		modelSyncES.ParkingTypes = append(modelSyncES.ParkingTypes, elastic_search.ParkingTypes{
			Type:   v.ParkingType,
			Status: v.Status,
			Price:  v.Price,
		})
	}
	return modelSyncES
}

type ParkingAddressModel struct {
	CityId     int     `json:"city_id"`
	DistrictId int     `json:"district_id"`
	WardId     int     `json:"ward_id"`
	Address    string  `json:"address"`
	Lat        float64 `json:"lat"`
	Lng        float64 `json:"lng"`
}

func (a ParkingAddressModel) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Make the Attrs struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (a *ParkingAddressModel) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

type ParkingQueryModel struct {
	Limit      int16
	Offset     int64
	Filter     ParkingFilterModel
	Update     ParkingUpdateModel
	Fields     []string
	UpdateUser string
}

type ParkingFilterModel struct {
	PublicId string
	Status   int
	Id       int
}

type ParkingUpdateModel struct {
	PublicId     string
	ParkingName  *string
	ParkingPhone *string
	OwnerName    *string
	OwnerPhone   *string
	Address      *string
	Lng          *float64
	Lat          *float64
	CityId       *int
	DistrictId   *int
	WardId       *int
	OpenAt       *time.Time
	CloseAt      *time.Time
	Status       *int
	Images       *string
	ParkingTypes *string
}

func (dto ParkingQueryModel) ToFilter() (conditions string, filterParams []interface{}) {
	conditionArray := make([]string, 0)
	filterParams = make([]interface{}, 0)
	counter := 1
	if len(dto.Filter.PublicId) > 0 {
		conditionArray = append(conditionArray, fmt.Sprintf("public_id = $%d", counter))
		filterParams = append(filterParams, dto.Filter.PublicId)
		counter += 1
	}
	if dto.Filter.Status > 0 {
		conditionArray = append(conditionArray, fmt.Sprintf("status = $%d", counter))
		filterParams = append(filterParams, dto.Filter.Status)
		counter += 1
	}
	if dto.Filter.Id > 0 {
		conditionArray = append(conditionArray, fmt.Sprintf("id = $%d", counter))
		filterParams = append(filterParams, dto.Filter.Id)
		counter += 1
	}
	conditions = strings.Join(conditionArray, ",")
	return
}

func (dto ParkingQueryModel) ToFilterUpdate() (conditions string, updateQuery string, filterParams []interface{}) {
	conditionArray := make([]string, 0)
	updateArray := make([]string, 0)
	filterParams = make([]interface{}, 0)
	counter := 1
	//update query
	if dto.Update.ParkingPhone != nil {
		updateArray = append(updateArray, fmt.Sprintf("parking_phone = $%d", counter))
		filterParams = append(filterParams, dto.Update.ParkingPhone)
		counter += 1
	}

	if dto.Update.ParkingName != nil {
		updateArray = append(updateArray, fmt.Sprintf("parking_name = $%d", counter))
		filterParams = append(filterParams, dto.Update.ParkingName)
		counter += 1
	}

	if dto.Update.OwnerPhone != nil {
		updateArray = append(updateArray, fmt.Sprintf("owner_phone = $%d", counter))
		filterParams = append(filterParams, dto.Update.OwnerPhone)
		counter += 1
	}

	if dto.Update.OwnerName != nil {
		updateArray = append(updateArray, fmt.Sprintf("owner_name = $%d", counter))
		filterParams = append(filterParams, dto.Update.OwnerName)
		counter += 1
	}

	if dto.Update.Address != nil {
		updateArray = append(updateArray, fmt.Sprintf("address['address'] = to_jsonb('%s'::text)", *dto.Update.Address))
	}
	if dto.Update.CityId != nil {
		updateArray = append(updateArray, fmt.Sprintf("address['city_id'] = to_jsonb(%d)", *dto.Update.CityId))
	}
	if dto.Update.DistrictId != nil {
		updateArray = append(updateArray, fmt.Sprintf("address['district_id'] = to_jsonb(%d)", *dto.Update.DistrictId))
	}
	if dto.Update.WardId != nil {
		updateArray = append(updateArray, fmt.Sprintf("address['ward_id'] = to_jsonb(%d)", *dto.Update.WardId))
	}
	if dto.Update.Lng != nil {
		updateArray = append(updateArray, fmt.Sprintf("address['lng'] = to_jsonb(%f)", *dto.Update.Lng))
	}
	if dto.Update.Lat != nil {
		updateArray = append(updateArray, fmt.Sprintf("address['lat'] = to_jsonb(%f)", *dto.Update.Lat))
	}
	if dto.Update.OpenAt != nil {
		updateArray = append(updateArray, fmt.Sprintf("open_at = $%d", counter))
		filterParams = append(filterParams, dto.Update.OpenAt)
		counter += 1
	}
	if dto.Update.CloseAt != nil {
		updateArray = append(updateArray, fmt.Sprintf("close_at = $%d", counter))
		filterParams = append(filterParams, dto.Update.CloseAt)
		counter += 1
	}
	if dto.Update.Status != nil {
		updateArray = append(updateArray, fmt.Sprintf("status = $%d", counter))
		filterParams = append(filterParams, dto.Update.Status)
		counter += 1
	}
	if dto.Update.Images != nil {
		updateArray = append(updateArray, fmt.Sprintf("images = $%d", counter))
		filterParams = append(filterParams, dto.Update.Images)
		counter += 1
	}
	if dto.Update.ParkingTypes != nil {
		updateArray = append(updateArray, fmt.Sprintf("parking_types = $%d", counter))
		filterParams = append(filterParams, dto.Update.ParkingTypes)
		counter += 1
	}

	if len(updateArray) == 0 {
		return
	}
	if len(dto.UpdateUser) == 0 {
		dto.UpdateUser = "system"
	}

	updateArray = append(updateArray, fmt.Sprintf("updated_name = $%d", counter))
	filterParams = append(filterParams, dto.UpdateUser)
	counter += 1
	updateArray = append(updateArray, fmt.Sprintf("updated_date = $%d", counter))
	filterParams = append(filterParams, time.Now())
	counter += 1
	if len(dto.Filter.PublicId) > 0 {
		conditionArray = append(conditionArray, fmt.Sprintf("public_id = $%d", counter))
		filterParams = append(filterParams, dto.Filter.PublicId)
		counter += 1
	}

	if dto.Filter.Status > 0 {
		conditionArray = append(conditionArray, fmt.Sprintf("status = $%d", counter))
		filterParams = append(filterParams, dto.Filter.Status)
		counter += 1
	}

	if dto.Filter.Id > 0 {
		conditionArray = append(conditionArray, fmt.Sprintf("id = $%d", counter))
		filterParams = append(filterParams, dto.Filter.Id)
		counter += 1
	}
	if len(conditionArray) == 0 {
		return
	}

	conditions = strings.Join(conditionArray, ",")
	updateQuery = strings.Join(updateArray, ",")

	return
}
