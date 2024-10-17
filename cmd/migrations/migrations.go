package main

import (
	embeds "auth"
	"auth/internal/config"

	"database/sql"
	"log/slog"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {

	connStr := config.MustLoad(nil).ConnString

	logger := initLogger()

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		logger.Error("failed to connnect to db", slog.Any("error", err), slog.String("connection string", connStr))
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("failed to close db connection", slog.Any("error", err))
		}
	}()

	goose.SetBaseFS(embeds.MigrationsFS)

	if err := goose.SetDialect("postgres"); err != nil {
		logger.Error("failed to set goose dialect", slog.Any("error", err))
		os.Exit(1)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		logger.Error("failed to Up", slog.Any("error", err))
		os.Exit(1)
	}

	logger.Info("migrations completed successfully")

}

func initLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelInfo,
		},
	))
}
