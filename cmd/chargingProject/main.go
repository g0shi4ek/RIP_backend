package main

import (
	"fmt"

	"github.com/g0shi4ek/RIP_backend/internal/app/config"
	"github.com/g0shi4ek/RIP_backend/internal/app/dsn"
	"github.com/g0shi4ek/RIP_backend/internal/app/handlers"
	"github.com/g0shi4ek/RIP_backend/internal/app/repository"
	"github.com/g0shi4ek/RIP_backend/internal/pkg"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	router := gin.Default()
	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	postgresString := dsn.FromEnv()
	fmt.Println(postgresString)

	repo, err := repository.NewChargingRepository(postgresString)
	if err != nil {
		logrus.Fatalf("error initializing repository: %v", err)
	}

	handler, err := handlers.NewChargingHandler(repo)
	if err != nil {
		logrus.Fatalf("error initializing handlers: %v", err)
	}

	application := pkg.NewApp(conf, router, handler)
	application.RunApp()
}