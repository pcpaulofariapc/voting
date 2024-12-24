package service

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/pcpaulofariapc/voting/config"
)

func ConnectMemcache() (cache *memcache.Client, err error) {
	memcacheConfig := config.Get().Services.Memcache

	cache = memcache.New(memcacheConfig.Host)
	err = cache.Ping()
	if err != nil {
		return nil, err
	}

	return cache, nil
}

func CreateProducer() (*kafka.Producer, error) {
	kafkaConfig := config.Get().Services.Kafka

	config := &kafka.ConfigMap{"bootstrap.servers": kafkaConfig.Host}
	producer, err := kafka.NewProducer(config)
	if err != nil {
		return nil, err
	}

	return producer, nil
}
