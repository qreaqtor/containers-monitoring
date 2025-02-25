package app

import (
	"github.com/qreaqtor/containers-monitoring/common/kafka/consumer"
	"github.com/qreaqtor/containers-monitoring/monitoring-backend/internal/config"
)

func getConsumerGroup(cfg config.KafkaConsumer) (*consumer.ConsumerGroup, error) {
	return consumer.NewConsumerGroup(cfg.Brokers, cfg.Topic, cfg.Group)
}
