package consumer_app

import (
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"parking-service/consumer/consumer_job"
	"parking-service/consumer/consumer_job/consumer_job_interface"
	"parking-service/provider/rabbitmq_provider"
	"parking-service/service/parking_module/sync_on_elastic"
)

type IConsumerApp interface {
	StartConsumer()
}

type ConsumerApp struct {
	IRabbitMQProvider rabbitmq_provider.IRabbitMQProvider
	ISyncOnElasticService sync_on_elastic.ISyncOnElasticService
}

func NewConsumerApp(IRabbitMQProvider rabbitmq_provider.IRabbitMQProvider,
	ISyncOnElasticService sync_on_elastic.ISyncOnElasticService) ConsumerApp {
	return ConsumerApp{IRabbitMQProvider: IRabbitMQProvider,
		ISyncOnElasticService: ISyncOnElasticService}
}

func (app *ConsumerApp) StartConsumer() {
	queueArr := app.IRabbitMQProvider.QueueDeclare()

	type queueHandle struct {
		Queue amqp091.Queue
		Job consumer_job_interface.IConsumerJob
	}
	qArr := make([]queueHandle , 0)
	for _ , v := range  queueArr {
		qArr = append(qArr , queueHandle{
			Queue: v,
			Job: consumer_job.SwitchJobByName(v.Name,app.ISyncOnElasticService),
		})
	}
	for _ , v := range qArr {
		fmt.Println("Init Job: "+ v.Queue.Name)
		go func(name string,job consumer_job_interface.IConsumerJob) {
			msgs, _ := app.IRabbitMQProvider.GetQueueChannel().Consume(
				name, // consumer_app
				"",     // consumer
				true,   // auto-ack
				false,  // exclusive
				false,  // no-local
				false,  // no-wait
				nil,    // args
			)
			go func() {
				for d := range msgs {
					retry  , err := job.Process(d)
					if retry {
						fmt.Println("Must retry")
					}
					if err != nil {
						fmt.Println("Error: " + err.Error())
					}
				}
			}()
		}(v.Queue.Name,v.Job)
	}

	var forever chan struct{}

	fmt.Println(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}