package app

import (
	"context"
	"errors"
	"io"

	"github.com/docker/docker/client"
	"github.com/qreaqtor/containers-monitoring/pinger/internal/config"
	containersinfo "github.com/qreaqtor/containers-monitoring/pinger/internal/containers"
	"github.com/qreaqtor/containers-monitoring/pinger/internal/usecase"

	comlog "github.com/qreaqtor/containers-monitoring/common/logging"
)

type runner interface {
	Run() error
}

type App struct {
	toClose []io.Closer

	pinger runner
}

func NewApp(ctx context.Context, cfg config.Config) (*App, error) {
	comlog.SetDefaultLogger(cfg.Env)

	toClose := make([]io.Closer, 0)

	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	toClose = append(toClose, dockerClient)

	containers, err := containersinfo.NewConmatinersInfo(ctx, dockerClient, cfg)
	if err != nil {
		return nil, err
	}

	pinger := usecase.NewPingerUsecase(ctx, containers, cfg.UpdateTimeout)

	app := &App{
		toClose: toClose,
		pinger:  pinger,
	}
	return app, nil
}

func (a *App) Run() error {
	return a.pinger.Run()
}

func (a *App) Close() error {
	errs := make([]error, 0, len(a.toClose))

	for _, closer := range a.toClose {
		err := closer.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
