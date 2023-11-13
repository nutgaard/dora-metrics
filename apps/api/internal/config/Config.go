package config

import "github.com/caarlos0/env/v10"
import "github.com/rs/zerolog/log"

type Config struct {
	Port string `env:"PORT" envDefault:"8080"`
	DB   PostgresqlConfig
}
type PostgresqlConfig struct {
	ConnectionUrl string `env:"DB_URL" envDefault:"postgres://username:password@localhost:8082/dora-metrics"`
}

func ReadConfig() *Config {
	config := &Config{}
	err := env.Parse(config)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	return config
}
