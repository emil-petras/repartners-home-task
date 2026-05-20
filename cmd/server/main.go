package main

import (
	"log"
	"net/http"

	"repartners-home-task/internal/config"
	"repartners-home-task/internal/database"
	"repartners-home-task/internal/handlers"
	"repartners-home-task/internal/services"
	web "repartners-home-task/web"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// load configuration from environment variables or defaults
	cfg := config.Load()

	// initialize sqlite database with the configured file path
	db, err := database.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	// create service layer instances for business logic
	packSizeService := services.NewPackSizeService(db)
	packagingService := services.NewPackagingService(packSizeService)
	packSizeHandler, err := handlers.NewPackSizeHandler(packSizeService, packagingService, web.FS)
	if err != nil {
		log.Fatalf("failed to create pack size handler: %v", err)
	}

	// create chi router for http request routing
	r := chi.NewRouter()

	// add middleware for logging and panic recovery
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// register all api routes from the handler
	packSizeHandler.RegisterRoutes(r)

	// health check endpoint for monitoring
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// start the http server and log the address
	log.Printf("starting server on %s", cfg.GetServerAddr())
	log.Printf("web ui available at http://localhost%s", cfg.GetServerAddr())
	if err := http.ListenAndServe(cfg.GetServerAddr(), r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
