package app

import (
	"log/slog"
	"time"

	"github.com/nhassl3/sso/internal/app/grpcapp"
	"github.com/nhassl3/sso/internal/services/auth"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	// TODO: initial storage
	var usrSaver auth.UserSaver
	var usrProvider auth.UserProvider
	var appProvider auth.AppProvider

	// TODO: init auth service
	auth := auth.New(log, usrSaver, usrProvider, appProvider, tokenTTL)

	grpcApp := grpcapp.New(log, grpcPort, auth)

	return &App{
		GRPCServer: grpcApp,
	}
}
