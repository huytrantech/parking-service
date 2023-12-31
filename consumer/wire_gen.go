// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"parking-service/consumer/consumer_app"
	"parking-service/provider/location_provider"
	"parking-service/provider/postgres_provider"
	"parking-service/provider/rabbitmq_provider"
	"parking-service/provider/viper_provider"
	"parking-service/proxy/elastic_search_proxy"
	"parking-service/repository"
	"parking-service/service/parking_module/sync_on_elastic"
)

// Injectors from wire.go:

func Initialize() consumer_app.ConsumerApp {
	iConfigProvider := viper_provider.NewConfigProvider()
	iRabbitMQProvider := rabbitmq_provider.NewRabbitMqProvider(iConfigProvider)
	iElasticSearchProxy := elastic_search_proxy.NewElasticSearchProxy(iConfigProvider)
	iPostgresProvider := postgres_provider.NewPostgresProvider(iConfigProvider)
	iParkingRepository := repository.NewParkingRepository(iPostgresProvider)
	iParkingSlotRepository := repository.NewParkingSlotRepository(iPostgresProvider)
	iLocationProvider := location_provider.NewLocationProvider()
	iSyncOnElasticService := sync_on_elastic.NewSyncOnElasticService(iElasticSearchProxy, iParkingRepository, iParkingSlotRepository, iLocationProvider)
	consumerApp := consumer_app.NewConsumerApp(iRabbitMQProvider, iSyncOnElasticService)
	return consumerApp
}
