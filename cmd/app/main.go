package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AntonTsoy/review-pull-request-service/internal/config"
	"github.com/AntonTsoy/review-pull-request-service/internal/database"
	"github.com/AntonTsoy/review-pull-request-service/internal/repository"
	"github.com/AntonTsoy/review-pull-request-service/internal/service"
	"github.com/AntonTsoy/review-pull-request-service/internal/transport/http/handlers"
	"github.com/AntonTsoy/review-pull-request-service/internal/transport/http/server"
)

func main() {
	log.Println("Initialize application...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	db, err := database.New(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to create DB instance: %v", err)
	}
	defer db.Close()

	if err = db.HealthCheck(ctx); err != nil {
		log.Fatalf("failed to open connection with database: %v", err)
	}
	log.Println("DB connection opened!")

	repository := repository.NewRepository()

	service := service.NewService(db, repository)

	handlers := handlers.NewHandlers(service)

	server := server.New(handlers)
	go func() {
		if err := server.Start(":8080"); err != nil {
			log.Printf("server error: %v", err)
		}
	}()

	sig := <-ctx.Done()
	log.Printf("Received signal: %v. Starting graceful shutdown...\n", sig)

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error during shutdown: %v", err)
		return
	}

	log.Println("Server stopped gracefully")
}
