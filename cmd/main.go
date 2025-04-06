package main

import (
	"fmt"
	"net/http"
	"test_api_go/configs"
	"test_api_go/internal/person"
)

func main() {
	conf := configs.LoadConfig()
	router := http.NewServeMux()
	person.NewPersonHandler(router, person.PersonHandlerDeps{
		Config: conf,
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
