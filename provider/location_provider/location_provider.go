package location_provider

import (
	"encoding/json"
	"io/ioutil"
)

type WardModel struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type DistrictModel struct {
	Id    int         `json:"id"`
	Name  string      `json:"name"`
	Wards []WardModel `json:"wards"`
}

type CityModel struct {
	Id        int             `json:"id"`
	Name      string          `json:"name"`
	Districts []DistrictModel `json:"districts"`
}

type MapWardDto struct {
	WardName string
}

type MapDistrictDto struct {
	DistrictName string
	WardMap map[int]MapWardDto
}

type MapCityDto struct {
	CityName string
	DistrictMap map[int]MapDistrictDto
}

type ILocationProvider interface {
	GetLocationData() []CityModel
	GetMapData() map[int]MapCityDto
}

type locationProvider struct {
	locationData []CityModel
	mapData map[int]MapCityDto
}

func NewLocationProvider() ILocationProvider {
	file, _ := ioutil.ReadFile("region.json")

	var data []CityModel

	_ = json.Unmarshal(file, &data)
	cityMap := make(map[int]MapCityDto)
	for _ , itemCity := range data{
		districtMap := make(map[int]MapDistrictDto)
		for _ ,itemDistrict := range itemCity.Districts{
			wardMap := make(map[int]MapWardDto)
			for _ , itemWard := range itemDistrict.Wards {
				wardMap[itemWard.Id] =MapWardDto{WardName: itemWard.Name}
			}
			districtMap[itemDistrict.Id] = MapDistrictDto{
				DistrictName: itemDistrict.Name,
				WardMap:      wardMap,
			}
		}
		cityMap[itemCity.Id] = MapCityDto{
			CityName:    itemCity.Name,
			DistrictMap: districtMap,
		}
	}

	return &locationProvider{locationData: data,mapData: cityMap}
}

func (provider *locationProvider) GetLocationData() []CityModel {
	return provider.locationData
}

func (provider *locationProvider) GetMapData() map[int]MapCityDto {
	return provider.mapData
}
