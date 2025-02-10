package api

import (
	"net/http"
	"strconv"

	"github.com/qreaqtor/containers-monitoring/monitoring-backend/internal/models"
)

const (
	defaultPageNumber = 0
	defaultPageSize   = 20
)

func getPaging(r *http.Request) models.Page {
	pageNumber, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		pageNumber = defaultPageNumber
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil {
		pageSize = defaultPageSize
	}

	return models.Page{
		Number: pageNumber,
		Size:   pageSize,
	}
}
