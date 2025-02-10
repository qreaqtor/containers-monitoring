package api

import (
	"net/http"
	"strconv"

	"github.com/qreaqtor/containers-monitoring/monitoring-backend/internal/models"
)

const (
	defaultPageNumber = 0
	defaultPageSize   = 20

	pageNumber = "pageNumber"
	pageSize = "pageSize"
)

func getPaging(r *http.Request) models.Page {
	pageNumber, err := strconv.Atoi(r.URL.Query().Get(pageNumber))
	if err != nil {
		pageNumber = defaultPageNumber
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get(pageSize))
	if err != nil {
		pageSize = defaultPageSize
	}

	return models.Page{
		Number: pageNumber,
		Size:   pageSize,
	}
}
