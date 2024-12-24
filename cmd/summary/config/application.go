package config

import "errors"

type Application struct {
	Name          string `json:"name"`
	Package       string `json:"package"`
	ExecutionTime int    `json:"execution_time"`
	TimeLocation  string `json:"time_location"`
}

func (a *Application) Validate() error {
	if a.Name == "" {
		return errors.New("o nome da aplicação não foi definido")
	}
	if a.Package == "" {
		return errors.New("o pacote da aplicação não foi definido")
	}

	if a.ExecutionTime < 0 {
		return errors.New("o intervalo de execução não pode ser menor ou igual a 0")
	}

	if a.TimeLocation == "" {
		return errors.New("o time location da aplicação não foi definida")
	}

	return nil
}
