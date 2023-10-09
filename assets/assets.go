package assets

import (
	"encoding/json"
	"io/ioutil"
)

type ConsumerDto struct {
	QueueName string `json:"QueueName"`
	Active    bool   `json:"Active"`
	Workers   int    `json:"Workers"`
}

func GetConsumerAssets() (consumer []ConsumerDto) {
	file, _ := ioutil.ReadFile("assets/consumer.json")

	_ = json.Unmarshal([]byte(file), &consumer)
	return
}
