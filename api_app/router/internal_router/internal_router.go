package internal_router

import (
	"github.com/labstack/echo/v4"
	"parking-service/api/controller/auth_controller"
	"parking-service/api/controller/internal_controller"
	"parking-service/middleware"
)

type IInternalRouter interface {
	RunRouterInternal(e *echo.Echo)
}

type Router struct {
	IInternalParkingController     internal_controller.IInternalParkingController
	IInternalParkingSlotController internal_controller.IInternalParkingSlotController
	IAuthController                auth_controller.IAuthController
	IInternalMiddleware            middleware.IInternalMiddleware
}

func NewInternalRouter(IInternalParkingController internal_controller.IInternalParkingController,
	IInternalParkingSlotController internal_controller.IInternalParkingSlotController,
	IAuthController auth_controller.IAuthController,
	IInternalMiddleware middleware.IInternalMiddleware) IInternalRouter {
	return &Router{IInternalParkingController: IInternalParkingController,
		IInternalParkingSlotController: IInternalParkingSlotController,
		IAuthController:                IAuthController,
		IInternalMiddleware:            IInternalMiddleware}
}

func (r *Router) RunRouterInternal(e *echo.Echo) {
	internalRouter := e.Group("/api/internal/v1", r.IInternalMiddleware.AuthTokenInside)
	{
		internalRouter.POST("/add", r.IInternalParkingController.CreateParkingInsideController)
		internalRouter.PUT("/:parking_id/update", r.IInternalParkingController.UpdateParkingController)
		internalRouter.PUT("/:parking_id/sync-es", r.IInternalParkingController.SyncParkingOnESController)
		internalRouter.PUT("/sync-multi-es/:parking_ids", r.IInternalParkingController.SyncMultiParkingOnESController)
		internalRouter.PUT("/approve-multi/:parking_ids", r.IInternalParkingController.ApprovedMultiParkingInsideController)
		internalRouter.PUT("/:parking_id/approve", r.IInternalParkingController.ApprovedParkingInsideController)
		internalRouter.PUT("/:parking_id/reopen", r.IInternalParkingController.ReopenParkingInternalController)
		internalRouter.GET("/:parking_id/detail", r.IInternalParkingController.DetailParkingController)
		internalRouter.PUT("/:parking_id/close", r.IInternalParkingController.CloseParkingInternalController)
		internalRouter.GET("/retrieve", r.IInternalParkingController.RetrieveListParkingController)
		internalRouter.POST("/generate-external-token", r.IInternalParkingController.GetExternalJWTToken)
		internalRouter.POST("/validate-external-token", r.IInternalParkingController.ValidateExternalToken)
		internalRouter.POST("/import", r.IInternalParkingController.ImportParkingController)

		internalParkingSlotRouter := internalRouter.Group("/:parking_id/parking-slot")
		{
			internalParkingSlotRouter.POST("/add", r.IInternalParkingSlotController.CreateParkingSlotController)
			internalParkingSlotRouter.PUT("/:parking_slot_id/update", r.IInternalParkingSlotController.UpdateParkingSlotController)
		}
	}

	authRouter := e.Group("api/internal/v1/auth")
	{
		authRouter.POST("/login", r.IAuthController.Login)
	}

}
