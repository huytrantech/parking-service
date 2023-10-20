package internal_controller

import (
	"context"
	"encoding/csv"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
	"parking-service/constants"
	"parking-service/core/utils/jwt_utils"
	"parking-service/model/api_model"
	"parking-service/service/parking_module/approve_parking"
	"parking-service/service/parking_module/close_parking"
	"parking-service/service/parking_module/create_parking"
	"parking-service/service/parking_module/detail_parking"
	"parking-service/service/parking_module/import_parking_csv"
	"parking-service/service/parking_module/reopen_parking"
	"parking-service/service/parking_module/retrieve_list_parking"
	"parking-service/service/parking_module/sync_on_elastic"
	"parking-service/service/parking_module/sync_parking_es_by_id"
	"parking-service/service/parking_module/update_parking"
	"strconv"
	"strings"
)

type IInternalParkingController interface {
	CreateParkingInsideController(c echo.Context) error
	ApprovedParkingInsideController(c echo.Context) error
	SyncParkingOnESController(c echo.Context) error
	SyncMultiParkingOnESController(c echo.Context) error
	RetrieveListParkingController(c echo.Context) error
	GetExternalJWTToken(c echo.Context) error
	ValidateExternalToken(c echo.Context) error
	CloseParkingInternalController(c echo.Context) error
	UpdateParkingController(c echo.Context) error
	ReopenParkingInternalController(c echo.Context) error
	ApprovedMultiParkingInsideController(c echo.Context) error
	DetailParkingController(c echo.Context) error
	ImportParkingController(c echo.Context) error
}

type internalParkingController struct {
	ICreateParkingService       create_parking.ICreateParkingService
	IApproveParkingService      approve_parking.IApproveParkingService
	ISyncMultiOnElasticService  sync_parking_es_by_id.ISyncMultiOnElasticService
	IRetrieveListParkingService retrieve_list_parking.IRetrieveListParkingService
	ICloseParkingService        close_parking.ICloseParkingService
	ISyncParkingESByIdService   sync_parking_es_by_id.ISyncParkingESByIdService
	IUpdateParkingService       update_parking.IUpdateParkingService
	ISyncOnElasticService       sync_on_elastic.ISyncOnElasticService
	IReopenParkingService       reopen_parking.IReopenParkingService
	IDetailParkingService       detail_parking.IDetailParkingService
	IImportParkingCSVService    import_parking_csv.IImportParkingCSVService
}

func NewInternalParkingController(
	ICreateParkingService create_parking.ICreateParkingService,
	IApproveParkingService approve_parking.IApproveParkingService,
	ISyncMultiOnElasticService sync_parking_es_by_id.ISyncMultiOnElasticService,
	ISyncParkingESByIdService sync_parking_es_by_id.ISyncParkingESByIdService,
	IRetrieveListParkingService retrieve_list_parking.IRetrieveListParkingService,
	ICloseParkingService close_parking.ICloseParkingService,
	ISyncOnElasticService sync_on_elastic.ISyncOnElasticService,
	IUpdateParkingService update_parking.IUpdateParkingService,
	IDetailParkingService detail_parking.IDetailParkingService,
	IReopenParkingService reopen_parking.IReopenParkingService,
	IImportParkingCSVService import_parking_csv.IImportParkingCSVService) IInternalParkingController {
	return &internalParkingController{ICreateParkingService: ICreateParkingService,
		IApproveParkingService:      IApproveParkingService,
		ISyncMultiOnElasticService:  ISyncMultiOnElasticService,
		IRetrieveListParkingService: IRetrieveListParkingService,
		ICloseParkingService:        ICloseParkingService,
		IUpdateParkingService:       IUpdateParkingService,
		ISyncOnElasticService:       ISyncOnElasticService,
		ISyncParkingESByIdService:   ISyncParkingESByIdService,
		IReopenParkingService:       IReopenParkingService,
		IDetailParkingService:       IDetailParkingService,
		IImportParkingCSVService:    IImportParkingCSVService,
	}
}

// CreateParkingInsideController godoc
// @Summary Add Parking
// @Description  Add Parking
// @Tags internal
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Accept json
// @Produce json
// @param model body api_model.CreateParkingRequestDto true "model"
// @Success 200 {object} api_model.BaseResponseAPIModel{data=boolean}
// @Router /api/internal/v1/add [post]
func (ctr *internalParkingController) CreateParkingInsideController(c echo.Context) error {
	var request api_model.CreateParkingRequestDto
	if err := c.Bind(&request); err != nil {
		return api_model.FailResponse(c, err)
	}
	request.Username = cast.ToString(c.Get(constants.HEADER_USERNAME))
	response, svErr := ctr.ICreateParkingService.CreateParking(context.Background(), request)
	if svErr != nil {
		return api_model.FailResponse(c, svErr)
	}
	return api_model.SuccessResponse(c, response)
}

// ApprovedParkingInsideController godoc
// @Summary Approved Parking
// @Description  Approved Parking
// @Tags internal
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Accept json
// @Produce json
// @param parking_id path string true "parking_id"
// @Success 200 {object} api_model.BaseResponseAPIModel{data=boolean}
// @Router /api/internal/v1/{parking_id}/approve [put]
func (ctr *internalParkingController) ApprovedParkingInsideController(c echo.Context) error {

	parkingId := c.Param("parking_id")

	svErr := ctr.IApproveParkingService.ApproveParking(context.Background(), parkingId,
		cast.ToString(c.Get(constants.HEADER_USERNAME)))
	if svErr != nil {
		return api_model.FailResponse(c, svErr)
	}
	return api_model.SuccessResponse(c, true)
}

// ApprovedMultiParkingInsideController godoc
// @Summary Approved Multi Parking
// @Description  Approved Multi Parking
// @Tags internal
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Accept json
// @Produce json
// @param parking_ids path string true "parking_ids"
// @Success 200 {object} api_model.BaseResponseAPIModel{data=boolean}
// @Router /api/internal/v1/approve-multi/{parking_ids} [put]
func (ctr *internalParkingController) ApprovedMultiParkingInsideController(c echo.Context) error {

	parkingIds := c.Param("parking_ids")

	arrParkingIds := strings.Split(parkingIds, ",")
	for _, id := range arrParkingIds {
		svErr := ctr.IApproveParkingService.ApproveParking(context.Background(), id,
			cast.ToString(c.Get(constants.HEADER_USERNAME)))
		if svErr != nil {
			continue
		}
	}
	return api_model.SuccessResponse(c, true)
}

// CloseParkingInternalController godoc
// @Summary Close Parking
// @Description  Close Parking
// @Tags internal
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Accept json
// @Produce json
// @param parking_id path string true "parking_id"
// @Success 200 {object} api_model.BaseResponseAPIModel{data=boolean}
// @Router /api/internal/v1/{parking_id}/close [put]
func (ctr *internalParkingController) CloseParkingInternalController(c echo.Context) error {

	parkingId := c.Param("parking_id")

	svErr := ctr.ICloseParkingService.CloseParking(context.Background(),
		parkingId, cast.ToString(c.Get(constants.HEADER_USERNAME)))
	if svErr != nil {
		return api_model.FailResponse(c, svErr)
	}
	return api_model.SuccessResponse(c, true)
}

// ReopenParkingInternalController godoc
// @Summary Reopen Parking
// @Description  Reopen Parking
// @Tags internal
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Accept json
// @Produce json
// @param parking_id path string true "parking_id"
// @Success 200 {object} api_model.BaseResponseAPIModel{data=boolean}
// @Router /api/internal/v1/{parking_id}/reopen [put]
func (ctr *internalParkingController) ReopenParkingInternalController(c echo.Context) error {

	parkingId := c.Param("parking_id")

	svErr := ctr.IReopenParkingService.ReopenParking(context.Background(),
		parkingId, cast.ToString(c.Get(constants.HEADER_USERNAME)))
	if svErr != nil {
		return api_model.FailResponse(c, svErr)
	}
	return api_model.SuccessResponse(c, true)
}

// SyncParkingOnESController godoc
// @Summary Sync Parking
// @Description  Sync Parking
// @Tags internal
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Accept json
// @Produce json
// @param parking_id path int true "parking_id"
// @Success 200 {object} api_model.BaseResponseAPIModel{data=boolean}
// @Router /api/internal/v1/:parking_id/sync-es [put]
func (ctr *internalParkingController) SyncParkingOnESController(c echo.Context) error {

	parkingId, _ := strconv.Atoi(c.Param("parking_id"))
	svErr := ctr.ISyncOnElasticService.SyncOnElastic(context.Background(), parkingId)
	if svErr != nil {
		return api_model.FailResponse(c, svErr)
	}
	return api_model.SuccessResponse(c, true)
}

// SyncMultiParkingOnESController godoc
// @Summary Sync multi Parking
// @Description  Sync multi Parking
// @Tags internal
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Accept json
// @Produce json
// @param parking_ids path string true "parking_ids"
// @Success 200 {object} api_model.BaseResponseAPIModel{data=boolean}
// @Router /api/internal/v1/sync-multi-es/{parking_ids} [put]
func (ctr *internalParkingController) SyncMultiParkingOnESController(c echo.Context) error {

	parkingIds := c.Param("parking_ids")
	if parkingIds == "all" {
		svErr := ctr.ISyncOnElasticService.SyncAllOnElastic(context.Background())
		if svErr != nil {
			return api_model.FailResponse(c, svErr)
		}
	}
	for _, value := range strings.Split(parkingIds, ",") {
		svErr := ctr.ISyncOnElasticService.SyncOnElastic(context.Background(), cast.ToInt(value))
		if svErr != nil {
			continue
		}
	}
	return api_model.SuccessResponse(c, true)
}

// RetrieveListParkingController godoc
// @Summary Retrieve paging list parking
// @Description  Retrieve paging list parking
// @Tags internal
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Accept json
// @Produce json
// @param page_index query int false "page_index"
// @param page_limit query int false "page_limit"
// @Success 200 {object} api_model.BaseResponseAPIModel{data=[]api_model.RetrieveListParkingResponseDto}
// @Router /api/internal/v1/retrieve [get]
func (ctr *internalParkingController) RetrieveListParkingController(c echo.Context) error {
	request := api_model.RetrieveListParkingRequestDto{
		PageIndex: cast.ToInt16(c.QueryParam("page_index")),
		PageLimit: cast.ToInt16(c.QueryParam("page_limit")),
	}
	response, svErr := ctr.IRetrieveListParkingService.RetrieveListParkingDatabase(context.Background(), request)
	if svErr != nil {
		return api_model.FailResponse(c, svErr)
	}
	return api_model.SuccessResponse(c, response)
}

// GetExternalJWTToken godoc
// @Summary Get External JWT Token
// @Description  Get External JWT Token
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Tags auth
// @Accept json
// @Produce json
// @Param model body api_model.GenJWTTokenExternalRequestDto true "model"
// @Success 200 {object} api_model.BaseResponseAPIModel{data=string}
// @Router /api/internal/v1/generate-external-token [post]
func (ctr *internalParkingController) GetExternalJWTToken(c echo.Context) error {

	var request api_model.GenJWTTokenExternalRequestDto
	if err := c.Bind(&request); err != nil {
		return api_model.FailResponse(c, err)
	}
	jwtToken, err := jwt_utils.GenJWTToken(request.Token)
	if err != nil {
		return api_model.FailResponse(c, err)
	}
	return api_model.SuccessResponse(c, jwtToken)
}

// UpdateParkingController godoc
// @Summary Update parking information
// @Description  Update parking information
// @Tags internal
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Accept json
// @Produce json
// @param parking_id path string true "parking_id"
// @param model body api_model.UpdateParkingRequestDto false "model"
// @Success 200 {object} api_model.BaseResponseAPIModel{data=boolean}
// @Router /api/internal/v1/{parking_id}/update [put]
func (ctr *internalParkingController) UpdateParkingController(c echo.Context) error {
	var request api_model.UpdateParkingRequestDto
	if err := c.Bind(&request); err != nil {
		return api_model.FailResponse(c, err)
	}
	request.Username = cast.ToString(c.Get(constants.HEADER_USERNAME))
	parkingId := c.Param("parking_id")
	err := ctr.IUpdateParkingService.UpdateParking(context.Background(), request, parkingId)
	if err != nil {
		return api_model.FailResponse(c, err)
	}
	return api_model.SuccessResponse(c, true)
}

func (ctr *internalParkingController) ValidateExternalToken(c echo.Context) error {
	type JwtTokenRequest struct {
		Token string `json:"token"`
	}
	var request JwtTokenRequest
	if err := c.Bind(&request); err != nil {
		return api_model.FailResponse(c, err)
	}
	jwtToken, err := jwt_utils.VerifyJWTTokenWithExpiredRequest(request.Token)
	if err != nil {
		return api_model.FailResponse(c, err)
	}
	return api_model.SuccessResponse(c, jwtToken)
}

// DetailParkingController godoc
// @Summary Detail Parking
// @Description  Detail Parking
// @Tags internal
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Accept json
// @Produce json
// @param parking_id path string true "parking_id"
// @Success 200 {object} api_model.BaseResponseAPIModel{data=boolean}
// @Router /api/internal/v1/{parking_id}/detail [get]
func (ctr *internalParkingController) DetailParkingController(c echo.Context) error {

	parkingId := c.Param("parking_id")
	resp, svErr := ctr.IDetailParkingService.GetDetailParking(context.Background(), parkingId)
	if svErr != nil {
		return api_model.FailResponse(c, svErr)
	}
	return api_model.SuccessResponse(c, resp)
}

// ImportParkingController godoc
// @Summary Import Parking
// @Description  Import Parking
// @Tags internal
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Accept json
// @Produce json
// @param file formData file true "file"
// @Success 200 {object} api_model.BaseResponseAPIModel{data=api_model.ImportParkingResponseDto}
// @Router /api/internal/v1/import [post]
func (ctr *internalParkingController) ImportParkingController(c echo.Context) error {

	file, err := c.FormFile("file")
	if err != nil {
		return api_model.FailResponse(c, err)
	}
	f, err := file.Open()
	if err != nil {
		return api_model.FailResponse(c, err)
	}

	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		return api_model.FailResponse(c, err)
	}
	resp := ctr.IImportParkingCSVService.ImportParkingCSV(context.Background(), data, cast.ToString(c.Get(constants.HEADER_USERNAME)))
	return api_model.SuccessResponse(c, resp)
}
