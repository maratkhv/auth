package app

import (
	grpcapp "auth/internal/app/grpc"
	"auth/internal/config"
	"auth/internal/services/auth"
	strg "auth/internal/storage"
	"log/slog"
)

type App struct {
	GRPCserver *grpcapp.App
	storage    strg.Storage
	log        *slog.Logger
}

func New(log *slog.Logger, cfg *config.Config) *App {

	storage := strg.MustNew(cfg.ConnString, cfg.Redis.Port)

	api := auth.New(log, storage, storage)

	srv := grpcapp.New(log, cfg.Server.Port, api)

	return &App{
		log:        log,
		storage:    storage,
		GRPCserver: srv,
	}
}

func (a *App) Shutdown() {
	a.GRPCserver.GracefulStop()
	if err := a.storage.Close(); err != nil {
		a.log.Error("failed to close connections", slog.Any("error", err))
	}
}
