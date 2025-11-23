package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/AntonTsoy/reveiw-pull-request-service/internal/config"
	"github.com/AntonTsoy/reveiw-pull-request-service/internal/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	db, err := database.New(ctx, cfg.DB)
	if err != nil {
		log.Fatalf("failed to create DB instance: %v", err)
	}
	defer db.Close()

	if err = db.HealthCheck(ctx); err != nil {
		log.Fatalf("failed to open connection with database: %v", err)
	}

	log.Println("DB connection opened!")

	// TODO: repository

	// TODO: service

	// TODO: Fiber http handlers

	// TODO: start server

	<-ctx.Done()
	log.Println("Shutting down...")
}
