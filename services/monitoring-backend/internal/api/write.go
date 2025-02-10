package api

import (
	"github.com/gorilla/websocket"
	"github.com/qreaqtor/containers-monitoring/common/result"
	"github.com/qreaqtor/containers-monitoring/monitoring-backend/internal/models"
)

func writeWS(conn *websocket.Conn, containersResult result.Result[models.ContainersInfo]) error {
	if containersResult.Error != nil {
		return containersResult.Error
	}

	return conn.WriteJSON(containersResult.Value)
}
