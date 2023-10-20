package internal_controller

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
	"parking-service/model/api_model"
	"parking-service/service/parking_slot_module/add_parking_slot"
	"parking-service/service/parking_slot_module/update_parking_slot"
)

type IInternalParkingSlotController interface {
	CreateParkingSlotController(c echo.Context) error
	UpdateParkingSlotController(c echo.Context) error
}

type internalParkingSLotController struct {
	IAddParkingSlotService    add_parking_slot.IAddParkingSlotService
	IUpdateParkingSlotService update_parking_slot.IUpdateParkingSlotService
}

func NewInternalParkingSlotController(
	IAddParkingSlotService add_parking_slot.IAddParkingSlotService,
	IUpdateParkingSlotService update_parking_slot.IUpdateParkingSlotService) IInternalParkingSlotController {
	return &internalParkingSLotController{
		IAddParkingSlotService:    IAddParkingSlotService,
		IUpdateParkingSlotService: IUpdateParkingSlotService,
	}
}

// CreateParkingSlotController godoc
// @Summary Add Parking Slot
// @Description  Add Parking Slot
// @Tags internal
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Accept json
// @Produce json
// @param model body api_model.AddParkingSlotRequest true "model"
// @param parking_id path int true "parking_id"
// @Success 200 {object} api_model.BaseResponseAPIModel{data=boolean}
// @Router /api/internal/v1/:parking_id/parking-slot/add [post]
func (ctr *internalParkingSLotController) CreateParkingSlotController(c echo.Context) error {
	var request api_model.AddParkingSlotRequest
	if err := c.Bind(&request); err != nil {
		return api_model.FailResponse(c, err)
	}
	request.ParkingId = cast.ToInt(c.Param("parking_id"))
	svErr := ctr.IAddParkingSlotService.AddParkingSlot(context.Background(), request)
	if svErr != nil {
		return api_model.FailResponse(c, svErr)
	}
	return api_model.SuccessResponse(c, true)
}

// UpdateParkingSlotController godoc
// @Summary UpdateParkingSlotController
// @Description  UpdateParkingSlotController
// @Tags internal
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Accept json
// @Produce json
// @param model body api_model.UpdateParkingSlotRequest true "model"
// @param parking_id path int true "parking_id"
// @param parking_slot_id path int true "parking_slot_id"
// @Success 200 {object} api_model.BaseResponseAPIModel{data=boolean}
// @Router /api/internal/v1/:parking_id/parking-slot/:parking_slot_id/update [post]
func (ctr *internalParkingSLotController) UpdateParkingSlotController(c echo.Context) error {
	var request api_model.UpdateParkingSlotRequest
	if err := c.Bind(&request); err != nil {
		return api_model.FailResponse(c, err)
	}
	request.ParkingId = cast.ToInt(c.Param("parking_id"))
	request.ParkingSlotId = cast.ToInt(c.Param("parking_slot_id"))
	svErr := ctr.IUpdateParkingSlotService.UpdateParkingSlot(context.Background(), request)
	if svErr != nil {
		return api_model.FailResponse(c, svErr)
	}
	return api_model.SuccessResponse(c, true)
}
