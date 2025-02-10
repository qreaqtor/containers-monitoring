package app

import (
	"context"
	"errors"
	"io"

	"github.com/gorilla/mux"
	appserver "github.com/qreaqtor/containers-monitoring/common/appServer"
	httpserver "github.com/qreaqtor/containers-monitoring/common/httpServer"
	"github.com/qreaqtor/containers-monitoring/common/kafka/consumer"
	comlog "github.com/qreaqtor/containers-monitoring/common/logging"
	"github.com/qreaqtor/containers-monitoring/monitoring-backend/internal/api"
	"github.com/qreaqtor/containers-monitoring/monitoring-backend/internal/config"
	"github.com/qreaqtor/containers-monitoring/monitoring-backend/internal/repo/postgres"
	"github.com/qreaqtor/containers-monitoring/monitoring-backend/internal/usecase"
)

type server interface {
	Start() error
	WaitAndClose() []error
}

type App struct {
	server server

	consumerGroup *consumer.ConsumerGroup
}

func StartNewApp(ctx context.Context, cfg config.Config) (*App, error) {
	comlog.SetLogger(cfg.Env)

	router := mux.NewRouter()

	conn, err := getPostgresConn(cfg.Postgres)
	if err != nil {
		return nil, err
	}

	repo := postgres.NewContainerRepo(conn, cfg.UpdatedPeriod)
	uc := usecase.NewContainerUC(ctx, repo, cfg.WsWritePeriod)
	api := api.NewContainersAPI(uc)
	api.Register(router)

	consumerGroup, err := getConsumerGroup(cfg.Kafka)
	if err != nil {
		return nil, err
	}

	httpServer := httpserver.NewHTTPServer(router)
	appServer := appserver.
		NewAppServer(ctx, httpServer, int(cfg.Port)).
		WithClosers([]io.Closer{conn, consumerGroup})

	app := &App{
		server:        appServer,
		consumerGroup: consumerGroup,
	}

	err = app.server.Start()
	if err != nil {
		return nil, err
	}

	err = app.consumerGroup.Start(ctx, uc.UpsertContainersHandler)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) Wait() error {
	return errors.Join(a.server.WaitAndClose()...)
}
