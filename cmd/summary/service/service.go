package service

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/pcpaulofariapc/voting/summary/config"
)

func CreateConsumer() (*kafka.Consumer, error) {
	kafkaConfig := config.Get().Services.Kafka
	groupID := "consumerVoting"

	config := &kafka.ConfigMap{
		"bootstrap.servers": kafkaConfig.Host,
		"group.id":          groupID,
		"auto.offset.reset": "latest",
	}
	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, err
	}
	return consumer, nil
}

func ConnectMemcache() (memcacheConnect *memcache.Client, err error) {
	memcacheConfig := config.Get().Services.Memcache
	memcacheConnect = memcache.New(memcacheConfig.Host)
	err = memcacheConnect.Ping()
	if err != nil {
		return nil, err
	}
	return memcacheConnect, nil
}
