package containersinfo

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/qreaqtor/containers-monitoring/pinger/internal/config"
	"github.com/qreaqtor/containers-monitoring/pinger/internal/models"
)

type ContainersInfo struct {
	ctx context.Context

	dockerClient *client.Client

	optionss container.ListOptions

	pingTimeout time.Duration
	pingCount   int

	lengthConatinerID uint
}

func NewContainersInfo(ctx context.Context, dockerClient *client.Client, cfg config.Config) (*ContainersInfo, error) {
	opts := container.ListOptions{
		All: true,
	}

	containersUC := &ContainersInfo{
		ctx:               ctx,
		dockerClient:      dockerClient,
		optionss:          opts,
		pingTimeout:       cfg.PingTimeout,
		lengthConatinerID: cfg.LengthConatinerID,
		pingCount:         int(cfg.PingCount),
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
		ipAddress := ""
		for _, network := range container.NetworkSettings.Networks {
			if network.IPAddress == "" {
				continue
			} else if err := ping(network.IPAddress, c.pingCount, c.pingTimeout); err != nil {
				slog.Info(err.Error())
				continue
			}
			ipAddress = network.IPAddress
		}
		if ipAddress == "" {
			continue
		}

		containerInfo := models.ContainerInfo{
			ID:     container.ID[:c.lengthConatinerID],
			Name:   container.Names[0],
			Image:  container.Image,
			Status: container.Status,
			Ports:  convertPorts(container.Ports),
			IP:     ipAddress,
		}
		containersInfo = append(containersInfo, containerInfo)
	}

	return containersInfo, nil
}

func convertPorts(ports []types.Port) []string {
	portsConverted := make([]string, 0, len(ports))
	for _, port := range ports {
		portsConverted = append(portsConverted, fmt.Sprintf("%v:%v", port.PublicPort, port.PrivatePort))
	}
	return portsConverted
}
