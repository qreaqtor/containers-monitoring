package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/qreaqtor/containers-monitoring/monitoring-backend/internal/models"
	"github.com/qreaqtor/containers-monitoring/monitoring-backend/internal/repo/schema"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

type ContainerRepo struct {
	db *bun.DB

	period time.Duration
}

func NewContainerRepo(sqldb *sql.DB, updatedPeriod time.Duration) *ContainerRepo {
	db := bun.NewDB(sqldb, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(false)))

	return &ContainerRepo{
		db: db,
		period: updatedPeriod,
	}
}

func (r *ContainerRepo) UpsertContainers(ctx context.Context, containers []models.ContainerInfo) error {
	containersSchema := make([]schema.ContainerInfo, 0, len(containers))
	for _, container := range containers {
		containersSchema = append(containersSchema, schema.NewContainerSchema(container))
	}

	_, err := r.db.NewInsert().
		Model(&containersSchema).
		On("CONFLICT (name) DO UPDATE").
		Set("id = EXCLUDED.id").
		Set("name = EXCLUDED.name").
		Set("image = EXCLUDED.image").
		Set("ipv4 = EXCLUDED.ipv4").
		Set("ports = EXCLUDED.ports").
		Set("state = EXCLUDED.state").
		Set("status = EXCLUDED.status").
		Set("updated_at = now()").
		Exec(ctx)
	return err
}

func (r *ContainerRepo) GetInfo(ctx context.Context, paging models.Page) ([]models.ContainerInfo, error) {
	offset := int(paging.Number * paging.Size)
	limit := int(paging.Size)

	containersSchema := make([]schema.ContainerInfo, 0, paging.Size)

	bottomUpdated := time.Now().Add(-1 * r.period)
	err := r.db.NewSelect().
		Model(&containersSchema).
		Where("updated_at > ?", bottomUpdated).
		Order("updated_at DESC").
		Offset(offset).
		Limit(limit).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	containers := make([]models.ContainerInfo, 0, len(containersSchema))
	for _, containerSchema := range containersSchema {
		containers = append(containers, containerSchema.ToDomainModel())
	}

	return containers, nil
}
