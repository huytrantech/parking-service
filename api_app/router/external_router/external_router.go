package external_router

import (
	"github.com/labstack/echo/v4"
	"parking-service/api/controller/external_controller"
	"parking-service/middleware"
)

type IExternalRouter interface {
	RunRouter(e *echo.Echo)
}

type Router struct {
	IExternalParkingController external_controller.IExternalParkingController
	IExternalMiddleware        middleware.IExternalMiddleware
}

func NewExternalRouter(
	IExternalParkingController external_controller.IExternalParkingController,
	IExternalMiddleware middleware.IExternalMiddleware) IExternalRouter {
	return &Router{
		IExternalParkingController: IExternalParkingController,
		IExternalMiddleware:        IExternalMiddleware}
}

func (r *Router) RunRouter(e *echo.Echo) {
	externalRouter := e.Group("/api/external/v1")
	{
		externalRouter.POST("/add", r.IExternalParkingController.AddParking)
		externalRouter.GET("/circle-location", r.IExternalParkingController.GetCircleParkingLocationController)
		externalRouter.GET("/directions", r.IExternalParkingController.GetDirectionController, r.IExternalMiddleware.AuthTokenPublic)
		externalRouter.GET("/recommend", r.IExternalParkingController.RecommendParkingLocationController)
		externalRouter.GET("/:parking_id", r.IExternalParkingController.GetDetailParkingController)
		externalRouter.GET("/placeholder", r.IExternalParkingController.GetPlaceHolder)
		externalRouter.GET("/placeholder/:place_id/detail", r.IExternalParkingController.GetDetailPlaceHolder)
	}
}
