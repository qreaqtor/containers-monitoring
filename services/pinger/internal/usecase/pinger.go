package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/qreaqtor/containers-monitoring/pinger/internal/models"
)

type containersInfo interface {
	GetInfo() ([]models.ContainerInfo, error)
}

type Pinger struct {
	ctx context.Context

	tiker *time.Ticker

	containers containersInfo
}

func NewPingerUsecase(ctx context.Context, containers containersInfo, updateTimeout time.Duration) *Pinger {
	return &Pinger{
		ctx:        ctx,
		containers: containers,
		tiker:      time.NewTicker(updateTimeout),
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
			if err != nil {
				slog.Error(err.Error())
				continue
			}

			for _, container := range containers {
				fmt.Println(container)
			}
			fmt.Println()
		}
	}
}
