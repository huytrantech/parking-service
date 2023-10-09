package goong_proxy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"parking-service/model/proxy_model/goong"
	"parking-service/provider/viper_provider"
	"strings"
)

type PlaceHolderPredictionDto struct {
	Description string `json:"description"`
	PlaceId string `json:"place_id"`
}

type PlaceHolderDto struct {
	Predictions []PlaceHolderPredictionDto `json:"predictions"`
	Error struct {
		Code string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type DetailPlaceDto struct {
	Result struct {
		PlaceId string `json:"place_id"`
		Geometry struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			}
		}
	} `json:"result"`
	Status string `json:"status"`
}
type IGoongProxy interface {
	GetPlaceHolder(ctx context.Context ,
		request goong.GetPlaceHolderRequestDto) (*PlaceHolderDto , error)
	GetDetailPlaceHolder(ctx context.Context ,
		placeId string) (*DetailPlaceDto , error)
}

type goongProxy struct {
	urlLink string
	key     string
	client  *http.Client
}

func NewGoongProxy(IConfigProvider viper_provider.IConfigProvider) IGoongProxy {
	//tr := &http.Transport{
	//	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	//}
	return &goongProxy{
		urlLink: strings.TrimRight(IConfigProvider.GetConfigEnv().GoongUrl,"/"),
		key:     IConfigProvider.GetConfigEnv().GoongKey,
		client:  &http.Client{},
	}
}

func (proxy *goongProxy) GetPlaceHolder(ctx context.Context ,
	request goong.GetPlaceHolderRequestDto) (*PlaceHolderDto , error)  {

	req, err := http.NewRequest("GET",fmt.Sprintf("%s/Place/AutoComplete?api_key=%s&input=%s&limit=%d&radius=%d",
		proxy.urlLink,proxy.key, url.QueryEscape(request.Input),request.Limit,request.Radius), nil)
	if err != nil {
		return nil , err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := proxy.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil , errors.New("Goong Error")
	}
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil , err
	}
	var response PlaceHolderDto
	if err = json.Unmarshal(body, &response); err != nil {
		return nil , err
	}

	if len(response.Error.Code) > 0 {
		return nil , errors.New(response.Error.Message)
	}

	return &response , nil
}

func (proxy *goongProxy) GetDetailPlaceHolder(ctx context.Context ,
	placeId string) (*DetailPlaceDto , error)  {

	req, err := http.NewRequest("GET",fmt.Sprintf("%s/Place/Detail?api_key=%s&place_id=%s",
		proxy.urlLink,proxy.key, url.QueryEscape(placeId)), nil)
	if err != nil {
		return nil , err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := proxy.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil , errors.New("Goong Error")
	}
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil , err
	}
	var response DetailPlaceDto
	if err = json.Unmarshal(body, &response); err != nil {
		return nil , err
	}

	if len(response.Result.PlaceId) == 0 {
		return nil , errors.New(response.Status)
	}

	return &response , nil
}