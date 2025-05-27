package main

import (
	"github.com/shoksin/quotes-service/configs"
	handler "github.com/shoksin/quotes-service/internal/delivery/http"
	"github.com/shoksin/quotes-service/internal/delivery/http/middleware"
	"github.com/shoksin/quotes-service/internal/repository"
	"github.com/shoksin/quotes-service/internal/storage"
	"github.com/shoksin/quotes-service/internal/usecase"
	"log"
	"net/http"
	"time"
)

func main() {
	cfg := configs.Load()

	db, err := storage.NewPostgresConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	quoteRepo := repository.NewQuoteRepository(db)

	quoteUseCase := usecase.NewQuoteUseCase(quoteRepo)

	quoteHandler := handler.NewQuoteHandler(quoteUseCase)

	router := http.NewServeMux()
	quoteHandler.RegisterRoutes(router)
	wrapped := middleware.LoggingMiddleware(router)

	addr := ":" + cfg.Server.Port
	log.Printf("Server starting on port %s", cfg.Server.Port)

	server := &http.Server{
		Addr:         addr,
		Handler:      wrapped,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err = server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}
