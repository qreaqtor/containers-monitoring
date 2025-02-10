package app

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/docker/docker/client"
	"github.com/qreaqtor/containers-monitoring/pinger/internal/config"
	containersinfo "github.com/qreaqtor/containers-monitoring/pinger/internal/containers"
	"github.com/qreaqtor/containers-monitoring/pinger/internal/usecase"
)

func (app *App) setup(ctx context.Context, cfg config.Config) error {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	app.toClose = append(app.toClose, dockerClient)

	containers, err := containersinfo.NewConmatinersInfo(ctx, dockerClient, cfg)
	if err != nil {
		return err
	}

	kafkaConf := sarama.NewConfig()
	producer, err := sarama.NewAsyncProducer(cfg.Kafka.Brokers, kafkaConf)
	if err != nil {
		return err
	}
	app.toClose = append(app.toClose, producer)

	app.pinger = usecase.NewPingerUsecase(ctx, containers, producer, cfg)

	return nil
}
