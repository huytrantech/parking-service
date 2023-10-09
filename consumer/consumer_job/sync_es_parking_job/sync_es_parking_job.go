package sync_es_parking_job

import (
	"context"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"parking-service/consumer/consumer_job/consumer_job_interface"
	"parking-service/core/utils/text_utils"
	"parking-service/model/consumer_model"
	"parking-service/service/parking_module/sync_on_elastic"
)

type syncESElasticJob struct {
	ISyncOnElasticService sync_on_elastic.ISyncOnElasticService
}

func NewSyncESElasticJob(ISyncOnElasticService sync_on_elastic.ISyncOnElasticService) consumer_job_interface.IConsumerJob {
	return &syncESElasticJob{
		ISyncOnElasticService: ISyncOnElasticService,
	}
}

func (job *syncESElasticJob) Process(delivery amqp091.Delivery) (retry bool ,err error) {
	fmt.Println("Message SyncES:" + string(delivery.Body))
	var request consumer_model.SyncEsJobRequest
	text_utils.ParseStringToStruct(string(delivery.Body) , &request)
	err = job.ISyncOnElasticService.SyncOnElastic(context.Background(),request.ParkingId)
	return
}