package api

import (
	"log"

	"github.com/g0shi4ek/RIP_backend/internal/app/handlers"
	"github.com/g0shi4ek/RIP_backend/internal/app/repository"
	"github.com/sirupsen/logrus"
)

func StartServer() {
	log.Println("Starting server")
	
	repo, err := repository.NewChargingRepository()
	if err != nil{
		logrus.Errorf("couldn't init repository, %v", err)
		return
	}
	handler, err := handlers.NewChargingHandler(repo)
	if err != nil{
		logrus.Errorf("couldn't init handlers, %v", err)
		return
	}
	router := handler.InitRoutes()
	logrus.Println("Starting server")
	log.Fatal(router.Run(":8080"))

  	log.Println("Server down")
}