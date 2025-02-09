package schema

import (
	"time"

	"github.com/qreaqtor/containers-monitoring/monitoring-backend/internal/models"
	"github.com/uptrace/bun"
)

type ContainerInfo struct {
	bun.BaseModel `bun:"table:containers"`

	Name      string `bun:",pk"`
	ID        string `bun:",notnull"`
	Image     string `bun:",notnull"`
	IPv4      string
	Ports     []string  `bun:",array"`
	State     string    `bun:",notnull"`
	Status    string    `bun:",notnull"`
	UpdatedAt time.Time `bun:",default:current_timestamp"`
}

func NewContainerSchema(container models.ContainerInfo) ContainerInfo {
	return ContainerInfo{
		Name:   container.Name,
		ID:     container.ID,
		Image:  container.Image,
		IPv4:   container.IPv4,
		Ports:  container.Ports,
		State:  container.State,
		Status: container.Status,
	}
}

func (c *ContainerInfo) ToDomainModel() models.ContainerInfo {
	return models.ContainerInfo{
		Name:      c.Name,
		ID:        c.ID,
		Image:     c.Image,
		IPv4:      c.IPv4,
		Ports:     c.Ports,
		State:     c.State,
		Status:    c.Status,
		UpdatedAt: c.UpdatedAt,
	}
}
