package main

import (
	"log"

	"github.com/g0shi4ek/RIP_backend/internal/app/config"
	"github.com/g0shi4ek/RIP_backend/internal/app/handlers"
	"github.com/g0shi4ek/RIP_backend/internal/app/repository"
	"github.com/g0shi4ek/RIP_backend/internal/app/service"
	"github.com/g0shi4ek/RIP_backend/internal/pkg/server"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	chargingRepo, err := repository.NewChargingRepository()
	if err != nil {
		log.Fatalf("error initializing repository: %v", err)
	}

	chargingService, err := service.NewChargingService(chargingRepo)
	if err != nil {
		log.Fatalf("error initializing repository: %v", err)
	}

	chargingHandler, err := handlers.NewChargingHandler(chargingService)
	if err != nil {
		log.Fatalf("error initializing handlers: %v", err)
	}

	application := server.NewApp(conf, router, chargingHandler)
	application.RunApp()
}