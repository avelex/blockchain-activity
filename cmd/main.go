package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/avelex/blockchain-activity/config"
	"github.com/avelex/blockchain-activity/internal/app"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	context.AfterFunc(ctx, func() {
		slog.Info("Interrupting...")
	})

	cfg := config.InitConfig()

	slog.Info("App run")
	if err := app.Run(ctx, cfg); err != nil {
		slog.Error(err.Error())
	}
	slog.Info("App done")
}
