package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/felipeosantos/goexpert/rate-limiter/config"
	"github.com/felipeosantos/goexpert/rate-limiter/internal/limiter"
	custommiddleware "github.com/felipeosantos/goexpert/rate-limiter/internal/middleware"
	"github.com/felipeosantos/goexpert/rate-limiter/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.Load(".", "env")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize storage based on configuration
	var store storage.Storage
	store, err = storage.New(cfg.StorageType, cfg.Storage[cfg.StorageType])
	if err != nil {
		log.Fatalf("Failed to connect to storage: %v", err)
	}

	defer store.Close()

	// Initialize rate limiter
	rateLimiter := limiter.New(store, limiter.Config{
		IP:    cfg.IP,
		Token: cfg.Token,
	})

	// Initialize router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Apply rate limiter middleware
	r.Use(custommiddleware.RateLimiterMiddleware(rateLimiter))

	// Define routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	// Start server
	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server starting on %s", serverAddr)
	if err := http.ListenAndServe(serverAddr, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
