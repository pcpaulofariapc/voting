package config

import (
	"errors"
	"fmt"
)

type Configuration struct {
	Application Application `json:"application"`
	Production  bool        `json:"production"`
	Database    Database    `json:"databases"`
	//Log         Log         `json:"log"`
	Services Services `json:"services"`
}

func (c *Configuration) Validate() (erro error) {

	if erro = c.Database.Validate(); erro != nil {
		return errors.Join(fmt.Errorf("erro ao validar o banco de dados"), erro)
	}

	if err := c.Application.Validate(); err != nil {
		return errors.Join(fmt.Errorf("erro ao validar informações da aplicacao"), erro)
	}

	if err := c.Services.Validate(); err != nil {
		return fmt.Errorf("erro ao validar serviços: %s", err.Error())
	}

	return
}
