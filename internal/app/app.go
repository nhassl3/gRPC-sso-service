package app

import (
	"log/slog"
	"time"

	"github.com/nhassl3/sso/internal/app/grpcapp"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	// TODO: initial storage

	// TODO: init auth service

	grpcApp := grpcapp.New(log, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
