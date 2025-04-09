package main

import (
	"fmt"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"test_api_go/configs"
	_ "test_api_go/docs"
	"test_api_go/internal/enrichment"
	"test_api_go/internal/person"
	"test_api_go/pkg/db"
)

// @title People API
// @version 1.0
// @description API для управления людьми с обогащением данных
// @host localhost:8081
// @BasePath /
func main() {
	conf := configs.LoadConfig()
	database := db.NewDb(conf)
	router := http.NewServeMux()

	router.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	enricher := enrichment.NewEnrichmentService()

	person.NewPersonHandler(router, person.PersonHandlerDeps{
		Config:   conf,
		Enricher: enricher,
		DB:       database.DB,
	})

	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	fmt.Println("Server is listening on port :8081")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Failed to start server", err)
	}
}
