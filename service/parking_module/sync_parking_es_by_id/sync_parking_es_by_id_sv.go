package sync_parking_es_by_id

import (
	"context"
	"parking-service/constants"
	"parking-service/model/consumer_model"
	"parking-service/provider/rabbitmq_provider"
)

type ISyncParkingESByIdService interface {
	SyncParkingESById(ctx context.Context , parkingId int) error
}

type syncParkingESByIdService struct {
	IRabbitMQProvider rabbitmq_provider.IRabbitMQProvider
}

func NewSyncParkingESByIdService() ISyncParkingESByIdService {
	return &syncParkingESByIdService{}
}

func(sv *syncParkingESByIdService) SyncParkingESById(ctx context.Context , parkingId int) error {

	err := sv.IRabbitMQProvider.PublishData(ctx , constants.SyncESParking,consumer_model.SyncEsJobRequest{
		ParkingId: parkingId,
		Source:    "parking-api",
	})

	if err != nil {
		return err
	}

	return nil
}