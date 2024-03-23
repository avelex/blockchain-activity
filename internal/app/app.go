package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/avelex/blockchain-activity/config"
	"github.com/avelex/blockchain-activity/internal/adapters/rpc/getblock"
	http_v1 "github.com/avelex/blockchain-activity/internal/controllers/http/v1"
	"github.com/avelex/blockchain-activity/internal/service"
	inmemory "github.com/avelex/blockchain-activity/internal/service/repository/in-memory"
	"github.com/avelex/blockchain-activity/pkg/jsonrpc"
)

func Run(ctx context.Context, cfg config.Config) error {
	jsonRPC := jsonrpc.New()
	getblockClient := getblock.New(jsonRPC, cfg.RPC.AccessToken)

	ctxAvailable, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := getblockClient.CheckAvailability(ctxAvailable); err != nil {
		return fmt.Errorf("getblock rpc not available: %w", err)
	}

	// TODO: clean cache with LRU
	// TODO: prepare and fill cache in background
	blockRepo := inmemory.New()
	blockchainService := service.New(blockRepo, getblockClient)

	v1Handler := http_v1.New(blockchainService)

	mux := http.NewServeMux()
	v1Handler.RegisterRoutes(mux)

	srv := http.Server{
		Addr:         cfg.HTTP.Host + ":" + cfg.HTTP.Port,
		Handler:      mux,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	slog.Info("Start http server", "host", cfg.HTTP.Host, "port", cfg.HTTP.Port)
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			slog.Warn(err.Error())
		}
	}()

	<-ctx.Done()

	ctxShutdown, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		return err
	}

	return nil
}
