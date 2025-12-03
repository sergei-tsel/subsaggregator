package main

import (
	"net/http"
	"subsaggregator/internal/db"
	"subsaggregator/internal/router"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	db.Init()

	r := router.NewRouter()

	http.ListenAndServe(":8080", r)
}
