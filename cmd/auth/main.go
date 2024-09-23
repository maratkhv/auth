package main

import (
	"auth/internal/app"
	"auth/internal/config"
	"auth/internal/lib/jwt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	log := config.InitLogger(cfg.Env)

	jwt.MustSetKeys(log)

	log.Info("starting application")

	application := app.New(log, cfg)

	go application.GRPCserver.MustRun()

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	sig := <-stop

	log.Info("stopping application", slog.String("signal", sig.String()))

	application.Shutdown()
}
