package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/ping", pong)

	return r
}

func pong(response http.ResponseWriter, request *http.Request) {
	response.Write([]byte("pong"))
}
