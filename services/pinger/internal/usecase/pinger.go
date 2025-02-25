package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/IBM/sarama"
	"github.com/qreaqtor/containers-monitoring/pinger/internal/config"
	"github.com/qreaqtor/containers-monitoring/pinger/internal/models"
)

type containersInfo interface {
	GetInfo() ([]models.ContainerInfo, error)
}

type Pinger struct {
	ctx context.Context

	tiker *time.Ticker

	containers containersInfo

	producer sarama.AsyncProducer
	topic    string
}

func NewPingerUsecase(ctx context.Context, containers containersInfo, producer sarama.AsyncProducer, cfg config.Config) *Pinger {
	return &Pinger{
		ctx:        ctx,
		containers: containers,
		tiker:      time.NewTicker(cfg.UpdateTimeout),
		producer:   producer,
		topic:      cfg.Kafka.Topic,
	}
}

func (p *Pinger) Run() error {
	for {
		select {
		case <-p.ctx.Done():
			p.tiker.Stop()
			return nil
		case <-p.tiker.C:
			containers, err := p.containers.GetInfo()
			if err == nil && len(containers) == 0 {
				err = errNoContainers
			}
			if err != nil {
				slog.Error(err.Error())
				continue
			}

			msg := models.NewContainersMsg(containers)
			p.sendMsg(msg)
		}
	}
}
