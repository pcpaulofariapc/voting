package config

import "errors"

type Services struct {
	Memcache Memcache `json:"memcache"`
	Kafka    Kafka    `json:"kafka"`
}

func (s *Services) Validate() error {
	if err := s.Memcache.Validate(); err != nil {
		return err
	}
	if err := s.Kafka.Validate(); err != nil {
		return err
	}
	return nil
}

type Nats struct {
	Host string `json:"host"`
}

type Memcache struct {
	Host string `json:"host"`
}

type Kafka struct {
	Host string `json:"host"`
}

func (m Memcache) Validate() error {
	if m.Host == "" {
		return errors.New("defina uma URL para o Memcache")
	}
	return nil
}

func (k Kafka) Validate() error {
	if k.Host == "" {
		return errors.New("defina uma URL para o Kafka")
	}
	return nil
}
