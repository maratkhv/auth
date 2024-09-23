package grpcapp

import (
	server "auth/internal/grpc"
	"auth/internal/lib/logger"
	"auth/internal/services/auth"
	"log/slog"
	"net"

	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCserver *grpc.Server
	port       string
}

func New(log *slog.Logger, port string, api *auth.Auth) *App {

	iLog := logger.WrapInterceptorLogger(log)

	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(
		grpclog.UnaryServerInterceptor(iLog),
	))

	server.RegisterServer(srv, api)
	return &App{
		log:        log,
		gRPCserver: srv,
		port:       ":" + port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	lsn, err := net.Listen("tcp", a.port)
	if err != nil {
		return err
	}

	return a.gRPCserver.Serve(lsn)
}

func (a *App) GracefulStop() {
	a.gRPCserver.GracefulStop()
}
