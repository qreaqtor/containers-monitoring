package config

import "time"

type Config struct {
	Env               string        `yaml:"env" env-required:"true"`
	PingTimeout       time.Duration `yaml:"ping_timeout" env-required:"true"`
	PingCount         uint          `yaml:"ping_count" env-default:"4"`
	UpdateTimeout     time.Duration `yaml:"update_timeout" env-required:"true"`
	LengthConatinerID uint          `yaml:"container_id_length" env-default:"12"`
	Kafka             KafkaProducer `yaml:"kafka_producer" env-required:"true"`
}

type KafkaProducer struct {
	Topic   string   `yaml:"topic" env-required:"true"`
	Brokers []string `yaml:"brokers" env-required:"true"`
}
