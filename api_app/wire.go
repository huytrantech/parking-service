//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
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

func Initialize() router.MainRouter {
	wire.Build(
		//Router
		router.NewMainRouter,
		external_router.NewExternalRouter,
		internal_router.NewInternalRouter,
		public_router.NewPublicRouter,
		//Controller
		internal_controller.NewInternalParkingController,
		public_controller.NewLocationController,
		auth_controller.NewAuthController,
		external_controller.NewExternalParkingController,
		internal_controller.NewInternalParkingSlotController,
		//Service
		create_parking.NewCreateParkingService,
		approve_parking.NewApproveParkingService,
		detail_parking.NewDetailParkingService,
		get_circle_parking_location.NewGetCircleParkingLocationService,
		recommend_parking_location.NewRecommendParkingLocationService,
		sync_parking_es_by_id.NewSyncMultiOnElasticService,
		update_parking.NewUpdateParkingService,
		sync_parking_es_by_id.NewSyncParkingESByIdService,
		reopen_parking.NewReopenParkingService,
		sync_on_elastic.NewSyncOnElasticService,
		get_direction.NewGetDirectionService,
		close_parking.NewCloseParkingService,
		retrieve_list_parking.NewRetrieveListParkingService,
		get_placeholder.NewGetPlaceHolderService,
		add_parking_slot.NewAddParkingSlotService,
		update_parking_slot.NewUpdateParkingSlotService,
		detail_place.NewDetailPlaceService,
		gen_public_id.NewGenPublicIdService,
		import_parking_csv.NewImportParkingCSVService,
		//Proxy
		elastic_search_proxy.NewElasticSearchProxy,
		google_map_proxy.NewGoogleMapProxy,
		account_proxy.NewAccountProxy,
		goong_proxy.NewGoongProxy,
		//Provider
		location_provider.NewLocationProvider,
		viper_provider.NewConfigProvider,
		redis_provider.NewRedisProvider,
		postgres_provider.NewPostgresProvider,
		//Middleware
		middleware.NewInternalMiddleware,
		middleware.NewExternalMiddleware,
		//Repository
		repository.NewParkingRepository,
		repository.NewParkingSlotRepository,
	)
	return router.MainRouter{}
}
