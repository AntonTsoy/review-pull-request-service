package main

import (
	"fmt"
	"log"

	"github.com/AntonTsoy/reveiw-pull-request-service/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v", cfg)

	// TODO: squirrel + pgx database

	// TODO: repository

	// TODO: service

	// TODO: Fiber http handlers

	// TODO: start service

	// TODO: graceful shutdown
}
