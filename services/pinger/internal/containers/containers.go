package containersinfo

import (
	"context"
	"log/slog"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/qreaqtor/containers-monitoring/pinger/internal/config"
	"github.com/qreaqtor/containers-monitoring/pinger/internal/models"
	"github.com/qreaqtor/containers-monitoring/pinger/internal/ping"
)

type ContainersInfo struct {
	ctx context.Context

	dockerClient *client.Client

	optionss container.ListOptions

	pingTimeout time.Duration

	lengthConatinerID uint

	outboundIP string
}

func NewConmatinersInfo(ctx context.Context, dockerClient *client.Client, cfg config.Config) (*ContainersInfo, error) {
	opts := container.ListOptions{
		All: true,
	}

	ip, err := ping.GetOutboundIP()
	if err != nil {
		return nil, err
	}

	containersUC := &ContainersInfo{
		ctx: ctx,
		dockerClient: dockerClient,
		optionss: opts,
		pingTimeout: cfg.PingTimeout,
		lengthConatinerID: cfg.LengthConatinerID,
		outboundIP: ip.String(),
	}
	return containersUC, nil
}

func (c *ContainersInfo) GetInfo() ([]models.ContainerInfo, error) {
	containers, err := c.dockerClient.ContainerList(c.ctx, c.optionss)
	if err != nil {
		return nil, err
	}

	containersInfo := make([]models.ContainerInfo, 0, len(containers))
	for _, container := range containers {
		ports, err := ping.PingPorts(container.Ports, c.pingTimeout)
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		containerInfo := models.ContainerInfo{
			ID:     container.ID[:c.lengthConatinerID],
			Name:   container.Names[0],
			Image:  container.Image,
			State:  container.State,
			Status: container.Status,
			Ports:  ports,
			IPv4:   c.outboundIP,
		}
		containersInfo = append(containersInfo, containerInfo)
	}

	return containersInfo, nil
}
