package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/yeralin-munar/tt-go-json-fernet/internal/service"
	"golang.org/x/sync/errgroup"
)

type App struct {
	Logger  log.Logger
	Service *service.TTGoJsonFernetService
}

func NewApp(
	logger log.Logger,
	service *service.TTGoJsonFernetService,
) *App {
	return &App{
		Logger:  logger,
		Service: service,
	}
}

func (a *App) Run() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	ctx, cancelFn := context.WithCancel(context.Background())

	go func() {
		<-quit
		log.Debug("Close app")
		cancelFn()
	}()

	g, gCtx := errgroup.WithContext(ctx)

	// Run
	g.Go(
		func() error {
			return a.Service.Run(gCtx)
		},
	)

	return g.Wait()
}
