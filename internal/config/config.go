package config

import (
	"errors"
	"flag"

	"github.com/caarlos0/env/v6"
)

type ServerConfig struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	DBUri                string `env:"DB_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func GetServerConfig() (*ServerConfig, error) {
	var endppointParam = flag.String("a", "localhost:6009", "Server endpoint")
	var dbURIParam = flag.String("d", "", "DB connection string")
	var accrualSystemAddressParam = flag.String("r", "localhost:9006", "Accrual system address")
	flag.Parse()
	var cfg ServerConfig
	err := env.Parse(&cfg)

	if err != nil {
		return nil, err
	}

	var endpoint *string
	var dbURI *string
	var accrualSystemAddress *string

	if len(cfg.DBUri) != 0 {
		dbURI = &cfg.DBUri
	} else if len(*dbURIParam) != 0 {
		dbURI = dbURIParam
	} else {
		return nil, errors.New("DB connection string is empty")
	}

	if len(cfg.RunAddress) != 0 {
		endpoint = &cfg.RunAddress
	} else {
		endpoint = endppointParam
	}

	if len(cfg.AccrualSystemAddress) != 0 {
		accrualSystemAddress = &cfg.AccrualSystemAddress
	} else {
		accrualSystemAddress = accrualSystemAddressParam
	}

	return &ServerConfig{
		RunAddress:           *endpoint,
		DBUri:                *dbURI,
		AccrualSystemAddress: *accrualSystemAddress,
	}, nil

}
