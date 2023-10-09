//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"parking-service/consumer/consumer_app"
	"parking-service/provider/location_provider"
	"parking-service/provider/postgres_provider"
	"parking-service/provider/rabbitmq_provider"
	"parking-service/provider/viper_provider"
	"parking-service/proxy/elastic_search_proxy"
	"parking-service/repository"
	"parking-service/service/parking_module/sync_on_elastic"
)

func Initialize() consumer_app.ConsumerApp {
	wire.Build(
		viper_provider.NewConfigProvider,
		rabbitmq_provider.NewRabbitMqProvider,
		consumer_app.NewConsumerApp,
		sync_on_elastic.NewSyncOnElasticService,
		repository.NewParkingRepository,
		postgres_provider.NewPostgresProvider,
		elastic_search_proxy.NewElasticSearchProxy,
		location_provider.NewLocationProvider,
		repository.NewParkingSlotRepository,
	)
	return consumer_app.ConsumerApp{}
}
