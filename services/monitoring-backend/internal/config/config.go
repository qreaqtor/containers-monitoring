package config

import "time"

type Config struct {
	Env           string         `yaml:"env" env-required:"true"`
	API           uint           `yaml:"api_version" env-required:"true"`
	Port          uint           `yaml:"app_port" env-required:"true"`
	Postgres      PostgresConfig `yaml:"postgres" env-required:"true"`
	UpdatedPeriod time.Duration  `yaml:"updated_period" env-default:"1h"`
}

type PostgresConfig struct {
	User     string `env:"user" env-required:"true"`
	Password string `env:"password" env-required:"true"`
	DB       string `env:"db" env-required:"true"`

	Host string `env:"host" env-required:"true"`
	Port int    `env:"port" env-required:"true"`

	SSL bool `env:"POSTGRES_SSL" env-required:"true"`
}
