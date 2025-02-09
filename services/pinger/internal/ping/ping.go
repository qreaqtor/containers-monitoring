package ping

import (
	"fmt"
	"net"
	"time"

	"github.com/docker/docker/api/types"
)

func PingPorts(containerPorts []types.Port, timeout time.Duration) ([]string, error) {
	ports := make([]string, 0, len(containerPorts))

	for _, port := range containerPorts {
		err := ping(port, timeout)
		if err != nil {
			return nil, err
		}
		ports = append(ports, convertPort(port))
	}

	return ports, nil
}

func ping(port types.Port, timeout time.Duration) error {
	portStr := fmt.Sprint(port.PublicPort)

	_, err := net.DialTimeout("tcp", net.JoinHostPort(port.IP, portStr), timeout)
	return err
}
