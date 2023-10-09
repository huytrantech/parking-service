package utils

import (
	"fmt"
	"math/rand"
	"parking-service/core/enums/parking_status_enum"
	"parking-service/model/api_model"
	"parking-service/model/database_model"
	"parking-service/provider/location_provider"
	"strings"
)

type ParkingRoles struct {
	CanApproved bool
	CanClose    bool
}

func GetParkingRoles(model database_model.ParkingModel) ParkingRoles {
	var roles ParkingRoles

	roles.CanApproved = CanApprove(model)
	roles.CanClose = CanClose(model)
	return roles
}

func CanApprove(model database_model.ParkingModel) bool {

	if model.Status == parking_status_enum.Pending().Data {
		return true
	}
	return false
}

func CanDenied(model database_model.ParkingModel) bool {

	if model.Status == parking_status_enum.Pending().Data {
		return true
	}
	return false
}

func CanClose(model database_model.ParkingModel) bool {
	return true
}

func CanBlock(model database_model.ParkingModel) bool {
	return true
}

func CanRemove(model database_model.ParkingModel) bool {
	if model.Status != parking_status_enum.Active().Data {
		return true
	}
	return false
}

func CanReopen(model database_model.ParkingModel) bool {
	if model.Status == parking_status_enum.Closed().Data ||
		model.Status == parking_status_enum.Block().Data {
		return true
	}
	return false
}

func GetKeyRedis(origin api_model.PointLocationApiDto, destination string) string {
	key := fmt.Sprintf("direction_%f,%f_%s",
		origin.Lat, origin.Lng, destination)
	return key
}

func GetAddressString(addressModel database_model.ParkingAddressModel, mapData map[int]location_provider.MapCityDto) string {
	city := mapData[addressModel.CityId]
	district := city.DistrictMap[addressModel.DistrictId]
	ward := district.WardMap[addressModel.WardId]
	addressStr := fmt.Sprintf("%s, %s, %s, %s",
		addressModel.Address,
		ward.WardName,
		district.DistrictName,
		city.CityName)

	return addressStr
}

const characters = "1234567890QWERTYUIOPASDFGHJKLZXCVBNM"

func GenParkingId() (parkingId string) {
	prefix := "PK"
	listCharacters := strings.Split(characters, "")
	listIds := make([]string, 0)
	for i := 0; i < 10; i++ {
		randIndex := rand.Intn(len(listCharacters))
		listIds = append(listIds, listCharacters[randIndex])
	}
	parkingId = fmt.Sprintf("%s_%s", prefix, strings.Join(listIds, ""))
	return
}
