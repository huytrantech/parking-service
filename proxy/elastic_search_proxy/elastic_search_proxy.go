package elastic_search_proxy

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"io/ioutil"
	"net/http"
	"parking-service/model/api_model"
	"parking-service/model/proxy_model/elastic_search"
	"parking-service/provider/viper_provider"
)

type IElasticSearchProxy interface {
	CreateParking(ctx context.Context, model elastic_search.ParkingModel) error
	GetCircleParkingLocation(ctx context.Context, request api_model.GetCircleParkingLocationRequestDto) (
		parking []elastic_search.ParkingModel, err error)
	RemoveParkingByQuery(ctx context.Context, query map[string]interface{}) error
	GetIndex() string
	GetDataElasticSearchCustomQuery(
		ctx context.Context, request map[string]interface{}) (parking []elastic_search.ParkingModel, err error)
}

type elasticSearchProxy struct {
	urlLink string
	key     string
	client  *http.Client
}

func NewElasticSearchProxy(IConfigProvider viper_provider.IConfigProvider) IElasticSearchProxy {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &elasticSearchProxy{
		urlLink: IConfigProvider.GetConfigEnv().ElasticSearchUrl,
		key:     IConfigProvider.GetConfigEnv().ElasticSearchToken,
		client:  &http.Client{Transport: tr},
	}
}

func (proxy *elasticSearchProxy) GetIndex() string {
	return "parking"
}
func (proxy *elasticSearchProxy) CreateParking(ctx context.Context, model elastic_search.ParkingModel) error {

	if model.ParkingId == 0 {
		fmt.Println("WARNING")
	}
	jsonData, err := json.Marshal(model)

	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", proxy.urlLink+"/"+proxy.GetIndex()+"/_doc/"+model.PublicId, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", proxy.key)
	_, err = proxy.client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (proxy *elasticSearchProxy) RemoveParkingByQuery(ctx context.Context, query map[string]interface{}) error {

	jsonData, err := json.Marshal(query)

	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", proxy.urlLink+"/"+proxy.GetIndex()+"/_delete_by_query", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", proxy.key)
	_, err = proxy.client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

type ElasticResponseModel struct {
	Hits struct {
		Hits []struct {
			Source elastic_search.ParkingModel `json:"_source"`
			Sort   []interface{}               `json:"sort"`
		} `json:"hits"`
	} `json:"hits"`
}

func (proxy *elasticSearchProxy) GetCircleParkingLocation(ctx context.Context, request api_model.GetCircleParkingLocationRequestDto) (
	parking []elastic_search.ParkingModel, err error) {

	filter := request.ConvertSearchCircleLocationRequest()
	jsonData, err := json.Marshal(filter)
	if err != nil {
		return
	}
	req, err := http.NewRequest("GET", proxy.urlLink+"/"+proxy.GetIndex()+"/_search", bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", proxy.key)
	resp, err := proxy.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return
	}

	var dataResp ElasticResponseModel
	if err = json.Unmarshal(body, &dataResp); err != nil {
		return
	}
	for _, value := range dataResp.Hits.Hits {
		item := value.Source
		if len(value.Sort) > 0 {
			item.Distance = cast.ToFloat64(value.Sort[0])
		}

		parking = append(parking, item)
	}
	return
}

func (proxy *elasticSearchProxy) GetDataElasticSearchCustomQuery(
	ctx context.Context, request map[string]interface{}) (parking []elastic_search.ParkingModel, err error) {

	jsonData, err := json.Marshal(request)
	if err != nil {
		return
	}
	req, err := http.NewRequest("GET", proxy.urlLink+"/"+proxy.GetIndex()+"/_search", bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", proxy.key)
	resp, err := proxy.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return
	}

	var dataResp ElasticResponseModel
	if err = json.Unmarshal(body, &dataResp); err != nil {
		return
	}
	for _, value := range dataResp.Hits.Hits {
		item := value.Source
		if len(value.Sort) > 0 {
			item.Distance = cast.ToFloat64(value.Sort[0])
		}

		parking = append(parking, item)
	}

	return
}
