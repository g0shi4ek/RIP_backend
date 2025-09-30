package main

import (
	"log"

	"github.com/g0shi4ek/RIP_backend/internal/app/config"
	"github.com/g0shi4ek/RIP_backend/internal/app/handlers"
	"github.com/g0shi4ek/RIP_backend/internal/app/repository"
	"github.com/g0shi4ek/RIP_backend/internal/app/service"
	"github.com/g0shi4ek/RIP_backend/internal/pkg/server"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/g0shi4ek/RIP_backend/docs"
)

// @title Charging Station API
// @version 1.0
// @description API for electric vehicle charging station management
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	application := server.NewApp(conf, router, chargingHandler)
	application.RunApp()
}