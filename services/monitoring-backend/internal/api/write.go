package api

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/qreaqtor/containers-monitoring/common/result"
	"github.com/qreaqtor/containers-monitoring/monitoring-backend/internal/models"
)

func writeWS(conn *websocket.Conn, containersResult result.Result[models.ContainersInfo]) error {
	if containersResult.Error != nil {
		return containersResult.Error
	}

	data, err := json.Marshal(containersResult.Value)
	if err != nil {
		return err
	}

	return conn.WriteJSON(data)
}
