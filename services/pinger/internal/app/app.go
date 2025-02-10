package app

import (
	"context"
	"errors"
	"io"

	"github.com/qreaqtor/containers-monitoring/pinger/internal/config"

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

	app := &App{
		toClose: make([]io.Closer, 0),
	}

	err := app.setup(ctx, cfg)
	if err != nil {
		return nil, err
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
