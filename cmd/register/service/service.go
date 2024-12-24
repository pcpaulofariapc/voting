package service

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/pcpaulofariapc/voting/register/config"
)

func CreateConsumer() (*kafka.Consumer, error) {
	kafkaConfig := config.Get().Services.Kafka
	groupID := "consumerVoting"

	config := &kafka.ConfigMap{
		"bootstrap.servers": kafkaConfig.Host,
		"group.id":          groupID,
		"auto.offset.reset": "earliest", //latest
	}
	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, err
	}
	return consumer, nil
}
