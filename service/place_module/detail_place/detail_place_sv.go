package detail_place

import (
	"context"
	"parking-service/constants"
	"parking-service/core"
	"parking-service/model/api_model"
	"parking-service/proxy/goong_proxy"
)

type IDetailPlaceService interface {
	GetDetailPlaceByPlaceId(ctx context.Context , placeId string) (
		*api_model.DetailPlaceResponseDto,error)
}

type detailPlaceService struct {
	IGoongProxy goong_proxy.IGoongProxy
}

func NewDetailPlaceService(IGoongProxy goong_proxy.IGoongProxy) IDetailPlaceService {
	return &detailPlaceService{IGoongProxy: IGoongProxy}
}

func (sv *detailPlaceService) GetDetailPlaceByPlaceId(ctx context.Context , placeId string) (
	*api_model.DetailPlaceResponseDto,error) {
	var response api_model.DetailPlaceResponseDto

	if len(placeId) == 0 {
		return nil , core.NewBadRequestErrorMessage("placeId is required")
	}
	detailPlace , err := sv.IGoongProxy.GetDetailPlaceHolder(ctx , placeId)
	if err != nil {
		return nil , core.NewInternalError(err , constants.ERR_100001)
	}

	if detailPlace == nil {
		return nil , core.NewBadRequestErrorMessage("Not found valid place")
	}
	response.Lng = detailPlace.Result.Geometry.Location.Lng
	response.Lat = detailPlace.Result.Geometry.Location.Lat

	return &response , nil
}