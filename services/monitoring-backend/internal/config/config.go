package config

import "time"

type Config struct {
	Env           string        `yaml:"env" env-required:"true"`
	API           uint          `yaml:"api_version" env-required:"true"`
	Port          uint          `yaml:"app_port" env-required:"true"`
	UpdatedPeriod time.Duration `yaml:"updated_period" env-default:"1h"`
	WsWritePeriod time.Duration `yaml:"ws_write_period" env-default:"10s"`
	Kafka         KafkaConsumer `yaml:"kafka_consumer" env-required:"true"`

	Postgres PostgresConfig `env-required:"true"`
}

type PostgresConfig struct {
	URL string `env:"PG_URL" env-required:"true"`
}

type KafkaConsumer struct {
	Topic   string   `yaml:"topic" env-required:"true"`
	Brokers []string `yaml:"brokers" env-required:"true"`
	Group   string   `yaml:"group" env-required:"true"`
}
