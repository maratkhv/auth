package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type Config struct {
	Env        string `yaml:"env" env-default:"prod"`
	ConnString string `yaml:"postgresConnString" env-required:"true"`
	Server     server `yaml:"server"`
	Redis      redis  `yaml:"redis"`
}

type redis struct {
	Addr string `yaml:"addr" env-default:"6379"`
}

type server struct {
	Port string `yaml:"port" env-default:"8082"`
	// TODO: think about right timeout value for default
	Timeout time.Duration `yaml:"timeout" env-default:"15m"`
}

type Options struct {
	// WithPath must be a path to a config file.
	// It will overwrite env variable.
	// If WithPath is empty then env variable CONFIG_PATH
	// will be used
	WithPath string
}

func MustLoad(opts *Options) *Config {
	var configPath string

	if opts != nil && opts.WithPath != "" {
		configPath = opts.WithPath
	} else {
		configPath = envConfigPath()
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic(err)
	}

	return &cfg
}

func envConfigPath() string {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		panic("env varible CONFIG_PATH not found")
	}

	if _, err := os.Stat(configPath); err != nil {
		panic(err)
	}

	return configPath
}

func InitLogger(env string) *slog.Logger {
	switch env {
	case envLocal:
		return slog.New(slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level:     slog.LevelDebug,
				AddSource: true,
			},
		))
	case envDev:
		return slog.New(slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		))
	case envProd:
		return slog.New(slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelInfo,
			},
		))
	}

	panic("something went wrong initing logger")
}
