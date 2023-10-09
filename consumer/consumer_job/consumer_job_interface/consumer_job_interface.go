package consumer_job_interface

import "github.com/rabbitmq/amqp091-go"

type IConsumerJob interface {
	Process(delivery amqp091.Delivery) (retry bool ,err error)
}
