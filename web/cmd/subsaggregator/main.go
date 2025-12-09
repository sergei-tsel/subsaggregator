package main

import (
	"log/slog"
	"net/http"
	"subsaggregator/internal/db"
	"subsaggregator/internal/router"

	"github.com/joho/godotenv"
	"github.com/natefinch/lumberjack"
)

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

	handler := slog.NewTextHandler(hook, &slog.HandlerOptions{
		AddSource: true,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)

	http.ListenAndServe(":8080", r)
}
