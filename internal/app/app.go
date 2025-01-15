package app

import (
	"log/slog"
	"time"

	"github.com/nhassl3/sso/internal/storage/sqlite"

	"github.com/nhassl3/sso/internal/app/grpcapp"
	"github.com/nhassl3/sso/internal/services/auth"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, grpcPort, authService)

	return &App{
		GRPCServer: grpcApp,
	}
}
