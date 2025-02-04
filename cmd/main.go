package main

import (
	"log/slog"
	"os"

	"github.com/Oxeeee/discont-bot/internal/bot"
	"github.com/Oxeeee/discont-bot/internal/config"
	"github.com/Oxeeee/discont-bot/internal/db"
	"github.com/Oxeeee/discont-bot/internal/repo"
	"github.com/Oxeeee/discont-bot/internal/services"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("starting application")

	db := db.ConntectDatabase(cfg).GetDB()
	
	userRepo := repo.NewUsersRepo(db)
	userService := services.NewUserService(userRepo, log)

	app := bot.New(log, cfg, userService)
	app.MustRun()
	
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
