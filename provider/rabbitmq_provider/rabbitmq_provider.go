package rabbitmq_provider

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"parking-service/assets"
	"parking-service/provider/viper_provider"
)

type IRabbitMQProvider interface {
	Close()
	GetConnection() *amqp.Connection
	GetQueueChannel() *amqp.Channel
	QueueDeclare() (queueArr []amqp.Queue)
	PublishData(ctx context.Context , queueName string , payload interface{}) error
}

type rabbitMqProvider struct {
	qChannel *amqp.Channel
	connection *amqp.Connection
}

func NewRabbitMqProvider(IConfigProvider viper_provider.IConfigProvider) IRabbitMQProvider {
	consumerConfig := assets.GetConsumerAssets()
	fmt.Println(consumerConfig)
	conn, err := amqp.Dial(IConfigProvider.GetConfigEnv().RabbitMQUrl)
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	return &rabbitMqProvider{connection: conn,qChannel: ch}
}

func(provider *rabbitMqProvider) Close(){
	provider.connection.Close()
	provider.qChannel.Close()
}

func (provider *rabbitMqProvider) GetConnection() *amqp.Connection{
	return provider.connection
}

func (provider *rabbitMqProvider) GetQueueChannel() *amqp.Channel{
	return provider.qChannel
}

func (provider *rabbitMqProvider) QueueDeclare() (queueArr []amqp.Queue){

	if provider.qChannel == nil {
		return
	}
	consumerConfig := assets.GetConsumerAssets()
	for _ , v := range consumerConfig {
		if v.Active == false {
			continue
		}
		q , err := provider.qChannel.QueueDeclare(
			v.QueueName, // name
			false,   // durable
			false,   // delete when unused
			false,   // exclusive
			false,   // no-wait
			nil,     // arguments
		)
		if err != nil {
			continue
		}
		queueArr = append(queueArr , q)
	}

	return
}

func(provider *rabbitMqProvider) PublishData(ctx context.Context , queueName string , payload interface{}) error {
	body , _ := json.Marshal(&payload)
	err := provider.qChannel.PublishWithContext(ctx,
		"",     // exchange
		queueName, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing {
			ContentType: "text/plain",
			Body:        body,
		})
	if err != nil {
		return err
	}
	return nil
}