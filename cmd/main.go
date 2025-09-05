package main

import (
	"log"

	"github.com/g0shi4ek/RIP_backend/internal/api"
)

func main() {
	log.Println("serving started")
	
	api.StartServer()
	log.Println("serving terminated")
}