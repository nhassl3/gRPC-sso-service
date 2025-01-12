package main

import (
	"log/slog"
	"os"

	"github.com/nhassl3/sso/internal/app"
	"github.com/nhassl3/sso/internal/config"
	"github.com/nhassl3/sso/internal/lib/logger/handlers/slogpretty"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// Load config
	cfg := config.MustLoad()
	// setup the logger
	log := setupLogger(cfg.Env)

	log.Info("starting application", slog.String("InformationLevel", cfg.Env))

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	application.GRPCServer.MustRun() // panic when erros occurs
}

// setupLogger this function provide logger for service
func setupLogger(env string) (log *slog.Logger) {
	switch env {
	case envLocal:
		log = setupPrettySlog(slog.LevelDebug) // for more perception
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = setupPrettySlog(slog.LevelDebug)
	}
	return
}

// setupPrettySlog for more perception information while
// service is running
func setupPrettySlog(level slog.Level) *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: level,
		},
	}
	return slog.New(
		opts.NewPrettyHandler(os.Stdout),
	)
}
