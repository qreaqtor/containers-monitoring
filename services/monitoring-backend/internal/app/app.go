package app

import (
	"context"
	"errors"
	"io"

	"github.com/gorilla/mux"
	appserver "github.com/qreaqtor/containers-monitoring/common/appServer"
	httpserver "github.com/qreaqtor/containers-monitoring/common/httpServer"
	comlog "github.com/qreaqtor/containers-monitoring/common/logging"
	"github.com/qreaqtor/containers-monitoring/monitoring-backend/internal/api"
	"github.com/qreaqtor/containers-monitoring/monitoring-backend/internal/config"
	"github.com/qreaqtor/containers-monitoring/monitoring-backend/internal/repo/postgres"
	"github.com/qreaqtor/containers-monitoring/monitoring-backend/internal/usecase"
)

type server interface {
	Start() error
	Wait() []error
}

type App struct {
	server server

	toClose []io.Closer
}

func NewApp(ctx context.Context, cfg config.Config) (*App, error) {
	comlog.SetDefaultLogger(cfg.Env)

	toClose := make([]io.Closer, 0)

	router := mux.NewRouter()

	appServer := appserver.NewAppServer(
		ctx,
		httpserver.NewHTTPServer(router),
		int(cfg.Port),
	)

	conn, err := getPostgresConn(cfg.Postgres)
	if err != nil {
		return nil, err
	}
	toClose = append(toClose, conn)

	repo := postgres.NewContainerRepo(conn, cfg.UpdatedPeriod)
	uc := usecase.NewContainerUC(ctx, repo, cfg.WsWritePeriod)
	api := api.NewContainersAPI(uc)
	api.Register(router)

	consumerGroup, err := startConsumer(ctx, cfg.Kafka, uc.UpsertContainersHandler)
	if err != nil {
		return nil, err
	}
	toClose = append(toClose, consumerGroup)

	app := &App{
		server:  appServer,
		toClose: toClose,
	}
	return app, nil
}

func (a *App) Start() error {
	return a.server.Start()
}

func (a *App) Wait() error {
	errs := a.server.Wait()

	for _, closer := range a.toClose {
		err := closer.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
