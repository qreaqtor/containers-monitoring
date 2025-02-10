package usecase

import (
	"context"
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
	"github.com/qreaqtor/containers-monitoring/monitoring-backend/internal/models"
)

type containersRepo interface {
	GetInfo(context.Context, models.Page) (models.ContainersInfo, error)
	UpsertContainers(context.Context, models.ContainersInfo) error
}

type ContainerUC struct {
	repo containersRepo

	wsWritePeriod time.Duration

	ctx context.Context
}

func NewContainerUC(ctx context.Context, repo containersRepo, wsWritePeriod time.Duration) *ContainerUC {
	return &ContainerUC{
		ctx:           ctx,
		repo:          repo,
		wsWritePeriod: wsWritePeriod,
	}
}

func (c *ContainerUC) UpsertContainersHandler(consumerMsg *sarama.ConsumerMessage) error {
	msg := new(models.ContainersInfo)

	err := json.Unmarshal(consumerMsg.Value, msg)
	if err != nil {
		return err
	}
	return c.repo.UpsertContainers(c.ctx, *msg)
}
