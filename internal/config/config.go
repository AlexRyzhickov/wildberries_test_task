package config

import "github.com/caarlos0/env/v6"

type Config struct {
	GetPort              string `env:"GET_PORT" envDefault:"8080"`
	SetPort              string `env:"SET_PORT" envDefault:"8081"`
	Priority             uint   `env:"PRIORITY" envDefault:"1"`
	ReplicasAvailability bool   `env:"REPLICAS_AVAILABILITY" envDefault:"false"`
	ReplicasPort         string `env:"REPLICA_PORT" envDefault:"8082"`
}

func New() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
