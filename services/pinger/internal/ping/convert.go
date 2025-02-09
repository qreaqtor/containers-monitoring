package ping

import (
	"fmt"

	"github.com/docker/docker/api/types"
)

func convertPort(port types.Port) string {
	return fmt.Sprintf("%v:%v", port.PublicPort, port.PrivatePort)
}
