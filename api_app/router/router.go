package router

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "parking-service/api/docs"
	"parking-service/api/router/external_router"
	"parking-service/api/router/internal_router"
	"parking-service/api/router/public_router"
	"parking-service/model/api_model"
)

type IMainRouter interface {
	RegisterRouter(e *echo.Echo)
}

type MainRouter struct {
	IExternalRouter external_router.IExternalRouter
	IInternalRouter internal_router.IInternalRouter
	IPublicRouter public_router.IPublicRouter
}

func NewMainRouter(
	IExternalRouter external_router.IExternalRouter,
	IInternalRouter internal_router.IInternalRouter,
	IPublicRouter public_router.IPublicRouter) MainRouter {
	return MainRouter{
		IExternalRouter: IExternalRouter,
		IInternalRouter: IInternalRouter,
		IPublicRouter: IPublicRouter,
	}
}

// @title Swagger Example API
// @version 1.0
// @description This is a sample server parking server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /v2
// @schemes http
func (r MainRouter) RegisterRouter(e *echo.Echo) {
	e.GET("", func(c echo.Context) error {
		return api_model.SuccessResponse(c,"Welcome api parking app")
	})

	e.GET("/swagger/*",echoSwagger.WrapHandler)
	r.IExternalRouter.RunRouter(e)
	r.IInternalRouter.RunRouterInternal(e)
	r.IPublicRouter.RunPublicRouter(e)
}
