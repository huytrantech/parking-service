// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"parking-service/api/controller/auth_controller"
	"parking-service/api/controller/external_controller"
	"parking-service/api/controller/internal_controller"
	"parking-service/api/controller/public_controller"
	"parking-service/api/router"
	"parking-service/api/router/external_router"
	"parking-service/api/router/internal_router"
	"parking-service/api/router/public_router"
	"parking-service/middleware"
	"parking-service/provider/location_provider"
	"parking-service/provider/postgres_provider"
	"parking-service/provider/redis_provider"
	"parking-service/provider/viper_provider"
	"parking-service/proxy/account_proxy"
	"parking-service/proxy/elastic_search_proxy"
	"parking-service/proxy/google_map_proxy"
	"parking-service/proxy/goong_proxy"
	"parking-service/repository"
	"parking-service/service/parking_module/approve_parking"
	"parking-service/service/parking_module/close_parking"
	"parking-service/service/parking_module/create_parking"
	"parking-service/service/parking_module/detail_parking"
	"parking-service/service/parking_module/gen_public_id"
	"parking-service/service/parking_module/get_circle_parking_location"
	"parking-service/service/parking_module/get_direction"
	"parking-service/service/parking_module/import_parking_csv"
	"parking-service/service/parking_module/recommend_parking_location"
	"parking-service/service/parking_module/reopen_parking"
	"parking-service/service/parking_module/retrieve_list_parking"
	"parking-service/service/parking_module/sync_on_elastic"
	"parking-service/service/parking_module/sync_parking_es_by_id"
	"parking-service/service/parking_module/update_parking"
	"parking-service/service/parking_slot_module/add_parking_slot"
	"parking-service/service/parking_slot_module/update_parking_slot"
	"parking-service/service/place_module/detail_place"
	"parking-service/service/place_module/get_placeholder"
)

// Injectors from wire.go:

func Initialize() router.MainRouter {
	iConfigProvider := viper_provider.NewConfigProvider()
	iElasticSearchProxy := elastic_search_proxy.NewElasticSearchProxy(iConfigProvider)
	iGetCircleParkingLocationService := get_circle_parking_location.NewGetCircleParkingLocationService(iElasticSearchProxy)
	iGoogleMapProxy := google_map_proxy.NewGoogleMapProxy(iConfigProvider)
	iRedisProvider := redis_provider.NewRedisProvider(iConfigProvider)
	iPostgresProvider := postgres_provider.NewPostgresProvider(iConfigProvider)
	iParkingRepository := repository.NewParkingRepository(iPostgresProvider)
	iGetDirectionService := get_direction.NewGetDirectionService(iGoogleMapProxy, iRedisProvider, iParkingRepository)
	iParkingSlotRepository := repository.NewParkingSlotRepository(iPostgresProvider)
	iGenPublicIdService := gen_public_id.NewGenPublicIdService(iParkingRepository)
	iCreateParkingService := create_parking.NewCreateParkingService(iParkingRepository, iParkingSlotRepository, iGenPublicIdService)
	iRecommendParkingLocationService := recommend_parking_location.NewRecommendParkingLocationService(iParkingRepository, iElasticSearchProxy)
	iGoongProxy := goong_proxy.NewGoongProxy(iConfigProvider)
	iGetPlaceHolderService := get_placeholder.NewGetPlaceHolderService(iGoongProxy)
	iDetailPlaceService := detail_place.NewDetailPlaceService(iGoongProxy)
	iLocationProvider := location_provider.NewLocationProvider()
	iDetailParkingService := detail_parking.NewDetailParkingService(iParkingRepository, iParkingSlotRepository, iLocationProvider)
	iExternalParkingController := external_controller.NewExternalParkingController(iGetCircleParkingLocationService, iGetDirectionService, iCreateParkingService, iRecommendParkingLocationService, iGetPlaceHolderService, iDetailPlaceService, iDetailParkingService)
	iExternalMiddleware := middleware.NewExternalMiddleware()
	iExternalRouter := external_router.NewExternalRouter(iExternalParkingController, iExternalMiddleware)
	iSyncOnElasticService := sync_on_elastic.NewSyncOnElasticService(iElasticSearchProxy, iParkingRepository, iParkingSlotRepository, iLocationProvider)
	iApproveParkingService := approve_parking.NewApproveParkingService(iParkingRepository, iElasticSearchProxy, iSyncOnElasticService, iLocationProvider)
	iSyncMultiOnElasticService := sync_parking_es_by_id.NewSyncMultiOnElasticService()
	iSyncParkingESByIdService := sync_parking_es_by_id.NewSyncParkingESByIdService()
	iRetrieveListParkingService := retrieve_list_parking.NewRetrieveListParkingService(iParkingRepository)
	iCloseParkingService := close_parking.NewCloseParkingService(iParkingRepository, iElasticSearchProxy)
	iUpdateParkingService := update_parking.NewUpdateParkingService(iParkingRepository, iElasticSearchProxy)
	iReopenParkingService := reopen_parking.NewReopenParkingService(iParkingRepository)
	iImportParkingCSVService := import_parking_csv.NewImportParkingCSVService(iParkingRepository, iGenPublicIdService)
	iInternalParkingController := internal_controller.NewInternalParkingController(iCreateParkingService, iApproveParkingService, iSyncMultiOnElasticService, iSyncParkingESByIdService, iRetrieveListParkingService, iCloseParkingService, iSyncOnElasticService, iUpdateParkingService, iDetailParkingService, iReopenParkingService, iImportParkingCSVService)
	iAddParkingSlotService := add_parking_slot.NewAddParkingSlotService(iParkingSlotRepository)
	iUpdateParkingSlotService := update_parking_slot.NewUpdateParkingSlotService(iParkingSlotRepository)
	iInternalParkingSlotController := internal_controller.NewInternalParkingSlotController(iAddParkingSlotService, iUpdateParkingSlotService)
	iAccountProxy := account_proxy.NewAccountProxy()
	iAuthController := auth_controller.NewAuthController(iAccountProxy)
	iInternalMiddleware := middleware.NewInternalMiddleware()
	iInternalRouter := internal_router.NewInternalRouter(iInternalParkingController, iInternalParkingSlotController, iAuthController, iInternalMiddleware)
	iLocationController := public_controller.NewLocationController(iLocationProvider)
	iPublicRouter := public_router.NewPublicRouter(iLocationController)
	mainRouter := router.NewMainRouter(iExternalRouter, iInternalRouter, iPublicRouter)
	return mainRouter
}
