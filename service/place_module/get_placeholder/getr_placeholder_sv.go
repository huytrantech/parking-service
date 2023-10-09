package get_placeholder

import (
	"context"
	"parking-service/constants"
	"parking-service/core"
	"parking-service/model/api_model"
	"parking-service/model/proxy_model/goong"
	"parking-service/proxy/goong_proxy"
)

type IGetPlaceHolderService interface {
	GetPlaceHolder(ctx context.Context ,
		request goong.GetPlaceHolderRequestDto)(response []api_model.GetPlaceHolderResponseDto , err error)
}

type getPlaceHolderService struct {
	IGoongProxy goong_proxy.IGoongProxy
}

func NewGetPlaceHolderService(IGoongProxy goong_proxy.IGoongProxy) IGetPlaceHolderService {
	return &getPlaceHolderService{IGoongProxy: IGoongProxy}
}

func (sv *getPlaceHolderService) GetPlaceHolder(ctx context.Context ,
	request goong.GetPlaceHolderRequestDto)(response []api_model.GetPlaceHolderResponseDto , err error) {

	if len(request.Input) == 0 {
		return
	}

	if request.Limit <= 0 {
		request.Limit = 10
	}
	if request.Radius <= 0 {
		request.Radius = 10
	}

	dataPlaceHolder , err := sv.IGoongProxy.GetPlaceHolder(ctx , request)
	if err != nil {
		err = core.NewInternalError(err , constants.ERR_100001)
		return
	}

	if dataPlaceHolder == nil {
		err = core.NewBadRequestErrorMessage("PlaceHolder Fail")
		return
	}

	response = make([]api_model.GetPlaceHolderResponseDto , len(dataPlaceHolder.Predictions))
	for index ,value := range dataPlaceHolder.Predictions {
		response[index] = api_model.GetPlaceHolderResponseDto{
			Place:   value.Description,
			PlaceId: value.PlaceId,
		}
	}


	return
}