package config

import (
	"encoding/json"
	"log"
	"os"
)

var config *Configuration

func Load() {
	path := "config.json"
	if val, set := os.LookupEnv("VOTING_CONFIG"); set && val != "" {
		path = val
	} else {
		log.Println("variável de ambiente `VOTING_CONFIG` não está definida, usado diretorio: ", path)
	}

	raw, erro := os.ReadFile(path)
	if erro != nil {
		log.Fatal(erro)
	}

	if erro = json.Unmarshal(raw, &config); erro != nil {
		log.Fatal(erro)
	}

	if erro = config.Validate(); erro != nil {
		log.Fatal(erro)
	}
}

func Get() *Configuration {
	if config == nil {
		Load()
	}
	return config
}
