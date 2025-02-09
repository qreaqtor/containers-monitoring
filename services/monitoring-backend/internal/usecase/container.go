package usecase

import (
	"context"

	"github.com/qreaqtor/containers-monitoring/monitoring-backend/internal/models"
)

type containersRepo interface {
	GetInfo(context.Context, models.Page) ([]models.ContainerInfo, error)
}

type ContainerUC struct {
	repo containersRepo
}

func NewContainerUC(repo containersRepo) *ContainerUC {
	return &ContainerUC{
		repo: repo,
	}
}

func (c *ContainerUC) GetInfo(ctx context.Context, page models.Page) ([]models.ContainerInfo, error) {
	return c.repo.GetInfo(ctx, page)
}
