package schema

import (
	"time"

	"github.com/qreaqtor/containers-monitoring/monitoring-backend/internal/models"
	"github.com/uptrace/bun"
)

type ContainerInfo struct {
	bun.BaseModel `bun:"table:containers"`

	Name      string    `bun:",pk"`
	ID        string    `bun:",notnull"`
	Image     string    `bun:",notnull"`
	IP        string    `bun:",notnull"`
	Ports     []string  `bun:",array"`
	Status    string    `bun:",notnull"`
	UpdatedAt time.Time `bun:",default:current_timestamp"`
}

func NewContainerSchema(container models.ContainerInfo) ContainerInfo {
	return ContainerInfo{
		Name:   container.Name,
		ID:     container.ID,
		Image:  container.Image,
		IP:     container.IP,
		Ports:  container.Ports,
		Status: container.Status,
	}
}

func (c *ContainerInfo) ToDomainModel() models.ContainerInfo {
	return models.ContainerInfo{
		Name:      c.Name,
		ID:        c.ID,
		Image:     c.Image,
		IP:        c.IP,
		Ports:     c.Ports,
		Status:    c.Status,
		UpdatedAt: c.UpdatedAt,
	}
}
