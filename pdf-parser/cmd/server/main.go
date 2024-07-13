package main

import (
	"log"

	"pdf-parser/config"
	"pdf-parser/internal/repository"
	"pdf-parser/internal/server"
	"pdf-parser/internal/service"
)

func main() {
	config, err := config.Load("./config/config.yaml")
	if err != nil {
		log.Fatalf("error loading config: %s", err.Error())
	}

	repo, err := repository.NewSQLite()
	if err != nil {
		log.Fatalf("error creating repository: %s", err.Error())
	}

	httpClient := service.NewHttpClient()

	parser := service.NewParser()

	service := service.New(repo, httpClient, parser)

	controller := server.New(service)

	log.Printf("initing service: %s", config.Name)
	controller.Start()
}
