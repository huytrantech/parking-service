package sync_parking_es_by_id

import (
	"context"
	"parking-service/constants"
	"parking-service/model/api_model"
	"parking-service/model/consumer_model"
	"parking-service/provider/rabbitmq_provider"
)

type ISyncMultiOnElasticService interface {
	SyncMultiOnElastic(ctx context.Context , parkingIds []int) (
		errRes api_model.ResponseErrorModel )
}

type syncMultiOnElasticService struct {
	IRabbitMQProvider rabbitmq_provider.IRabbitMQProvider
}

func NewSyncMultiOnElasticService() ISyncMultiOnElasticService {
	return &syncMultiOnElasticService{}
}

func(sv *syncMultiOnElasticService) SyncMultiOnElastic(ctx context.Context , parkingIds []int) (
	errRes api_model.ResponseErrorModel ) {

	for _ , parkingItem := range parkingIds {
		err := sv.IRabbitMQProvider.PublishData(ctx , constants.SyncESParking,consumer_model.SyncEsJobRequest{
			ParkingId: parkingItem,
			Source:    "parking-api",
		})

		if err != nil {
			continue
		}
	}


	return
}