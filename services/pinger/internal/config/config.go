package config

import "time"

type Config struct {
	Env               string        `yaml:"env" env-required:"true"`
	PingTimeout       time.Duration `yaml:"ping_timeout" env-required:"true"`
	UpdateTimeout     time.Duration `yaml:"update_timeout" env-required:"true"`
	LengthConatinerID uint          `yaml:"container_id_length" env-default:"12"`
}
