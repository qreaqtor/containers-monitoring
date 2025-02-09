package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	confcom "github.com/qreaqtor/containers-monitoring/common/config"
	"github.com/qreaqtor/containers-monitoring/monitoring-backend/internal/app"
	"github.com/qreaqtor/containers-monitoring/monitoring-backend/internal/config"
)

func main() {
	cfg, err := confcom.Load[config.Config]()
	if err != nil {
		log.Fatalln(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	app, err := app.NewApp(ctx, *cfg)
	if err != nil {
		log.Fatalln(err)
	}

	err = app.Start()
	if err != nil {
		log.Fatalln(err)
	}

	err = app.Wait()
	if err != nil {
		log.Fatalln(err)
	}
}
