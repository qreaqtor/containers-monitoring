package api

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	logmsg "github.com/qreaqtor/containers-monitoring/common/logging/message"
	"github.com/qreaqtor/containers-monitoring/common/result"
	"github.com/qreaqtor/containers-monitoring/common/web"
	"github.com/qreaqtor/containers-monitoring/monitoring-backend/internal/models"
)

type containersGetter interface {
	GetInfoChan(context.Context, models.Page) <-chan result.Result[models.ContainersInfo]
	GetInfo(context.Context, models.Page) (models.ContainersInfo, error)
}

type ContainersAPI struct {
	containers containersGetter

	upgrader *websocket.Upgrader
}

func NewContainersAPI(containers containersGetter) *ContainersAPI {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return &ContainersAPI{
		containers: containers,
		upgrader:   &upgrader,
	}
}

const (
	pageSize   = "pageSize"
	pageNumber = "pageNumber"
)

func (c *ContainersAPI) Register(r *mux.Router) {
	paging := []string{
		pageNumber, `{pageNumber:[0-9]+}`,
		pageSize, `{pageSize:[0-9]+}`,
	}

	r.Path("/info/ws").HandlerFunc(c.getInfoWS).Methods(http.MethodGet).Queries(paging...)
	r.Path("/info").HandlerFunc(c.getInfo).Methods(http.MethodGet).Queries(paging...)
}

func (c *ContainersAPI) getInfo(w http.ResponseWriter, r *http.Request) {
	msg := logmsg.NewLogMsg(r.Context(), r.RequestURI, r.Method)

	page := getPaging(r)

	containersInfo, err := c.containers.GetInfo(r.Context(), page)
	if err != nil {
		web.WriteError(w, msg.WithText(err.Error()).WithStatus(http.StatusBadRequest))
		return
	}

	web.WriteData(w,
		msg.WithText("OK").WithStatus(http.StatusOK),
		containersInfo,
	)
}

func (c *ContainersAPI) getInfoWS(w http.ResponseWriter, r *http.Request) {
	logMsg := logmsg.NewLogMsg(r.Context(), r.RequestURI, r.Method)

	conn, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logMsg.WithText(err.Error()).Error()
		return
	}
	defer conn.Close()

	page := getPaging(r)

	containersChan := c.containers.GetInfoChan(r.Context(), page)

	for containersResult := range containersChan {
		err := writeWS(conn, containersResult)
		if err != nil {
			logMsg.WithText(err.Error()).Error()
			break
		}
	}
}
