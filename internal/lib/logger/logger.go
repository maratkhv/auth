package logger

import (
	"context"
	"log/slog"

	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

type Interceptor struct {
	log *slog.Logger
}

func WrapInterceptorLogger(log *slog.Logger) grpclog.Logger {
	return &Interceptor{log: log}
}

func (i *Interceptor) Log(ctx context.Context, level grpclog.Level, msg string, fields ...any) {
	if level == grpclog.LevelError {
		i.log.ErrorContext(ctx, msg, fields...)
	} else {
		i.log.Debug(msg)
	}
}
