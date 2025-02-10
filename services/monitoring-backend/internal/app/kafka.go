package app

import (
	"context"

	"github.com/qreaqtor/containers-monitoring/common/kafka/consumer"
	"github.com/qreaqtor/containers-monitoring/monitoring-backend/internal/config"
)

func startConsumer(ctx context.Context, cfg config.KafkaConsumer, handler consumer.HandleFunc) (*consumer.ConsumerGroup, error) {
	consumerGroup, err := consumer.NewConsumerGroup(cfg.Brokers, cfg.Topic, cfg.Group)
	if err != nil {
		return nil, err
	}

	err = consumerGroup.Start(ctx, handler)
	if err != nil {
		return nil, err
	}

	return consumerGroup, nil
}
