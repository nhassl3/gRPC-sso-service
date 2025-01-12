package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	authgRPC "github.com/nhassl3/sso/internal/grpc/auth"
	"google.golang.org/grpc"
)

const (
	opRun  = "grpcapp.Run"
	opStop = "grpcapp.Stop"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, port int) *App {
	gRPCServer := grpc.NewServer()

	authgRPC.Register(gRPCServer)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

// MustRun runs gRPC server and panics if any error occurs
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run gRPC server
func (a *App) Run() error {
	log := a.log.With(slog.String("op", opRun), slog.Int("port", a.port))

	log.Info("starting gRPC sevrer")

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", opRun, err)
	}

	log.Info("grpc server is running", slog.String("address", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", opRun, err)
	}

	return nil
}

// Stop stops gRPC server
func (a *App) Stop() {
	a.log.With(slog.String("op", opStop)).Info("stopping gRPC server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop() // Graceful shutdown
}
