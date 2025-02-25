package usecase

import (
	"context"
	"time"

	"github.com/qreaqtor/containers-monitoring/common/result"
	"github.com/qreaqtor/containers-monitoring/monitoring-backend/internal/models"
)

func (c *ContainerUC) GetInfo(ctx context.Context, page models.Page) (models.ContainersInfo, error) {
	return c.repo.GetInfo(ctx, page)
}

func (c *ContainerUC) GetInfoChan(ctx context.Context, page models.Page) <-chan result.Result[models.ContainersInfo] {
	output := make(chan result.Result[models.ContainersInfo])

	go func() {
		defer close(output)
		c.writeInfo(ctx, page, output)
	}()

	return output
}

func (c *ContainerUC) writeInfo(ctx context.Context, page models.Page, output chan result.Result[models.ContainersInfo]) {
	ticker := time.NewTicker(c.wsWritePeriod)
	defer ticker.Stop()

	for {
		containers, err := c.repo.GetInfo(ctx, page)
		if err != nil {
			output <- result.NewErrorResult[models.ContainersInfo](err)
			return
		}
		output <- result.NewValueResult(containers)

		select {
		case <-ticker.C:
		case <-ctx.Done():
			return
		}
	}
}
