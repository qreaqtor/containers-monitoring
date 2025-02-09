package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	confcom "github.com/qreaqtor/containers-monitoring/common/config"
	"github.com/qreaqtor/containers-monitoring/pinger/internal/app"
	"github.com/qreaqtor/containers-monitoring/pinger/internal/config"
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
	defer app.Close()

	err = app.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
