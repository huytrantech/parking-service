package public_router

import (
	"github.com/labstack/echo/v4"
	"parking-service/api/controller/public_controller"
)

type IPublicRouter interface {
	RunPublicRouter(e *echo.Echo)
}

type PublicRouter struct {
	ILocationController public_controller.ILocationController
}

func NewPublicRouter(ILocationController public_controller.ILocationController) IPublicRouter {
	return &PublicRouter{ILocationController: ILocationController}
}

func (r *PublicRouter) RunPublicRouter(e *echo.Echo) {
	publicRouter := e.Group("/api/public/v1")
	{
		locationRouter := publicRouter.Group("/location")
		{
			locationRouter.GET("/cities" , r.ILocationController.GetCitiesController)
			locationRouter.GET("/cities/:city_id/districts" , r.ILocationController.GetDistrictsController)
			locationRouter.GET("/cities/:city_id/districts/:district_id/wards" , r.ILocationController.GetWardsController)
		}
	}
}
