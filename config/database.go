package config

import "errors"

type Database struct {
	Nickname               string `json:"nickname"`
	Name                   string `json:"name"`
	User                   string `json:"user"`
	Password               string `json:"password"`
	Host                   string `json:"host"`
	Port                   string `json:"port"`
	ConnectionsMaximum     int    `json:"connections_maximum"`
	ConnectionsMaximumIdle int    `json:"connections_maximum_idle"`
	MaximumTimeOpenTX      int    `json:"maximum_time_open_tx"`
	ReadOnly               bool   `json:"read_only"`
}

func (b *Database) Validate() error {
	if b.Nickname == "" {
		return errors.New("defina um apelido para o banco de dados")
	}

	if b.Name == "" {
		return errors.New("defina um nome para o banco de dados")
	}

	if b.User == "" {
		return errors.New("defina um usuario para o banco de dados")
	}

	if b.Password == "" {
		return errors.New("defina uma senha para o banco de dados")
	}

	if b.Host == "" {
		return errors.New("defina um host para o banco de dados")
	}

	if b.Port == "0" {
		return errors.New("defina uma porta para o banco de dados")
	}

	return nil
}
