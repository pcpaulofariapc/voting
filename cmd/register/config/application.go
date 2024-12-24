package config

import "errors"

type Application struct {
	Name         string `json:"name"`
	Package      string `json:"package"`
	TimeLocation string `json:"time_location"`
}

func (a *Application) Validate() error {
	if a.Name == "" {
		return errors.New("o nome da aplicação não foi definido")
	}
	if a.Package == "" {
		return errors.New("o pacote da aplicação não foi definido")
	}
	if a.TimeLocation == "" {
		return errors.New("o time location da aplicação não foi definida")
	}

	return nil
}
