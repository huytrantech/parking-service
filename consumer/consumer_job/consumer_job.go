package consumer_job

import (
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"parking-service/constants"
	"parking-service/consumer/consumer_job/consumer_job_interface"
	"parking-service/consumer/consumer_job/sync_es_parking_job"
	"parking-service/service/parking_module/sync_on_elastic"
)


type baseJob struct {

}

func NewBaseJob() consumer_job_interface.IConsumerJob {
	return &baseJob{}
}

func (job *baseJob) Process(delivery amqp091.Delivery) (retry bool ,err error) {
	fmt.Println("Message:" + string(delivery.Body))
	return
}

func SwitchJobByName(name string,
	ISyncOnElasticService sync_on_elastic.ISyncOnElasticService) consumer_job_interface.IConsumerJob {
	switch name {
	case constants.SyncESParking:
		return sync_es_parking_job.NewSyncESElasticJob(ISyncOnElasticService)
	default:
		return NewBaseJob()
	}
}
