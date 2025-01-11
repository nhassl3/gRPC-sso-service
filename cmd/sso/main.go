package main

import (
	"fmt"
	"log/slog"

	"github.com/nhassl3/sso/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)
}

func setupLogger(env string) *slog.Logger {
	
}