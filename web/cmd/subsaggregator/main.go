package main

import (
	"log/slog"
	"net/http"
	"os"
	_ "subsaggregator/docs"
	"subsaggregator/internal/db"
	"subsaggregator/internal/router"

	"github.com/joho/godotenv"
	"github.com/natefinch/lumberjack"
)

// @title		Subsaggregator
// @version		1.0
// @description	Агрегатор подписок
// @host		localhost:8080
// @BasePath	/
// @schemes		http
func main() {
	godotenv.Load(".env")

	db.Init()

	r := router.NewRouter()

	hook := &lumberjack.Logger{
		Filename:   "./logs/app.log",
		MaxSize:    500,
		MaxBackups: 7,
		MaxAge:     7,
		Compress:   false,
	}

	var level slog.Level
	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	handler := slog.NewTextHandler(hook, &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)

	http.ListenAndServe(":8080", r)
}
