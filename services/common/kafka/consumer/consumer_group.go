package consumer

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"sync"
	"sync/atomic"

	"github.com/IBM/sarama"
)

type HandleFunc func(*sarama.ConsumerMessage) error

type ConsumerGroup struct {
	ready   chan bool
	client  sarama.ConsumerGroup
	topic   string
	group   string
	handler HandleFunc
	started atomic.Bool
}

func NewConsumerGroup(brokers []string, topic string, group string) (*ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Version = sarama.DefaultVersion
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	client, err := sarama.NewConsumerGroup(brokers, group, config)
	if err != nil {
		return nil, err
	}

	return &ConsumerGroup{
		ready:  make(chan bool),
		client: client,
		group:  group,
		topic:  topic,
	}, nil
}

func (consumer *ConsumerGroup) Start(ctx context.Context, handler HandleFunc) error {
	if consumer.started.Swap(true) {
		return errors.New("already started")
	}

	consumer.handler = handler

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		
		for {
			if err := consumer.client.Consume(ctx, []string{consumer.topic}, consumer); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				log.Panicf("Error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()
	wg.Wait()

	return nil
}

func (consumer *ConsumerGroup) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.ready)
	return nil
}

func (consumer *ConsumerGroup) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *ConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				slog.Info("message channel was closed")
				return nil
			}
			slog.Debug("Message claimed",
				slog.String("value", string(message.Value)),
				slog.Time("timestamp", message.Timestamp),
				slog.String("topic", message.Topic),
			)
			err := consumer.handler(message)
			if err != nil {
				slog.Error(err.Error())
				continue
			}
			session.MarkMessage(message, "")
			session.Commit()
		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/IBM/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}

func (consumer *ConsumerGroup) Close() error {
	return consumer.client.Close()
}
