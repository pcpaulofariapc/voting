package config

import "time"

func GetTimeLocation() *time.Location {
	fh, erro := time.LoadLocation(config.Application.TimeLocation)
	if erro != nil {
		fh, _ = time.LoadLocation("America/Sao_Paulo")
	}
	return fh
}
