package main

import (
	"net/http"

	"github.com/LucasBelusso1/MultithreadingChallange/internal/infra/webserver/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get(`/search_cep/{cep:^\d{8}$}`, handlers.GetAddress)

	http.ListenAndServe(":8000", r)
}
