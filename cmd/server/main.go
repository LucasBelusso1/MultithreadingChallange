package main

import (
	"net/http"

	"github.com/LucasBelusso1/MultithreadingChallange/internal/infra/webserver/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	_ "github.com/LucasBelusso1/MultithreadingChallange/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Multithreading CEP API
// @version         1.0
// @description     API to search CEP using multithreading
// @termsOfService  http://swagger.io/terms/

// @contact.name   Lucas belusso
// @contact.email  belussolucas@gmail.com

// @host      localhost:8000
// @BasePath  /search_cep
func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get(`/search_cep/{cep:^\d{8}$}`, handlers.GetAddress)

	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8000/docs/doc.json")))

	http.ListenAndServe(":8000", r)
}
