package google_map_proxy

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"parking-service/core"
	"parking-service/provider/viper_provider"
)

type PointLocation struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
type DirectionResponseDto struct {
	GeocodedWaypoints []struct {
		GeocoderStatus string   `json:"geocoder_status"`
		PlaceId        string   `json:"place_id"`
		Types          []string `json:"types"`
	} `json:"geocoded_waypoints"`
	Routes []struct {
		Legs []struct {
			Steps []struct {
				EndLocation   PointLocation `json:"end_location"`
				StartLocation PointLocation `json:"start_location"`
			} `json:"steps"`
		} ` json:"legs"`
	} `json:"routes"`
	Status string `json:"status"`
}

type IGoogleMapProxy interface {
	GetDirection(ctx context.Context, origin PointLocation,
		destination PointLocation) (listPoints []PointLocation, err error)
}

type googleMapProxy struct {
	key    string
	domain string
	client  *http.Client
}

func NewGoogleMapProxy(IConfigProvider viper_provider.IConfigProvider) IGoogleMapProxy {
	return &googleMapProxy{
		key:    IConfigProvider.GetConfigEnv().GoogleMapKey,
		domain: IConfigProvider.GetConfigEnv().GoogleMapDomain,
		client:  &http.Client{},
	}
}

func (proxy *googleMapProxy) GetDirection(ctx context.Context, origin PointLocation,
	destination PointLocation) (listPoints []PointLocation, err error) {

	req, err := http.NewRequest("GET",
		fmt.Sprintf("%s/directions/json?origin=%f,%f&destination=%f,%f&key=%s",
			proxy.domain, origin.Lat, origin.Lng, destination.Lat, destination.Lng, proxy.key), nil)

	if err != nil {
		return
	}
	resp, err := proxy.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var dataResp DirectionResponseDto
	if err = json.Unmarshal(body, &dataResp); err != nil {
		return
	}
	if dataResp.Status == "OVER_QUERY_LIMIT" {
		err = core.NewForbidden("OVER_QUERY_LIMIT")
		return
	}
	if len(dataResp.Routes) > 0 && len(dataResp.Routes[0].Legs) > 0{
		for index , value := range dataResp.Routes[0].Legs[0].Steps{
			listPoints = append(listPoints , PointLocation{
				Lat: value.StartLocation.Lat,
				Lng: value.StartLocation.Lng,
			})
			if index + 1 == len(dataResp.Routes[0].Legs[0].Steps) {
				listPoints = append(listPoints , PointLocation{
					Lat: value.EndLocation.Lat,
					Lng: value.EndLocation.Lng,
				})
			}
		}
	}


	return
}
