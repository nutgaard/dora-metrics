package config

import "github.com/caarlos0/env/v10"
import "github.com/rs/zerolog/log"

type Config struct {
	Port string `env:"PORT" envDefault:"8080"`
	DB   PostgresqlConfig
}
type PostgresqlConfig struct {
	Host     string `env:"DB_HOST" envDefault:"localhost:8082/dora-metrics"`
	Username string `env:"DB_USERNAME" envDefault:"username"`
	Password string `env:"DB_PASSWORD" envDefault:"password"`
}

func ReadConfig() *Config {
	config := &Config{}
	err := env.Parse(config)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	return config
}
