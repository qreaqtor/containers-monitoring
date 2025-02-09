package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	logmsg "github.com/qreaqtor/containers-monitoring/common/logging/message"
	"github.com/qreaqtor/containers-monitoring/common/web"
	"github.com/qreaqtor/containers-monitoring/monitoring-backend/internal/models"
)

type containersGetter interface {
	GetInfo(context.Context, models.Page) ([]models.ContainerInfo, error)
}

type ContainersAPI struct {
	containers containersGetter
}

func NewContainersAPI(containers containersGetter) *ContainersAPI {
	return &ContainersAPI{
		containers: containers,
	}
}

const (
	pageSize   = "pageSize"
	pageNumber = "pageNumber"
)

func (c *ContainersAPI) Register(r *mux.Router) {
	paging := []string{
		pageNumber, `{pageNumber:\d+}`,
		pageSize, `{pageSize:[1-9][\d+]?}`,
	}

	r.Path("/info").HandlerFunc(c.getInfo).Methods(http.MethodGet).Queries(paging...)
}

func (c *ContainersAPI) getInfo(w http.ResponseWriter, r *http.Request) {
	msg := logmsg.NewLogMsg(r.Context(), r.RequestURI, r.Method)

	number, err := strconv.Atoi(r.URL.Query().Get(pageNumber))
	if err != nil {
		web.WriteError(w, msg.With(err.Error(), http.StatusBadRequest))
		return
	}
	size, err := strconv.Atoi(r.URL.Query().Get(pageSize))
	if err != nil {
		web.WriteError(w, msg.With(err.Error(), http.StatusBadRequest))
		return
	}

	page := models.Page{
		Number: uint(number),
		Size:   uint(size),
	}

	containersInfo, err := c.containers.GetInfo(r.Context(), page)
	if err != nil {
		web.WriteError(w, msg.With(err.Error(), http.StatusBadRequest))
		return
	}

	web.WriteData(w,
		msg.With("OK", http.StatusOK),
		containersInfo,
	)
}
