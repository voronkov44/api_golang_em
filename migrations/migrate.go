package main

import (
	"test_api_go/configs"
	"test_api_go/internal/person"
	"test_api_go/pkg/db"
)

func main() {
	conf := configs.LoadConfig()
	database := db.NewDb(conf)

	err := database.AutoMigrate(&person.Person{})
	if err != nil {
		panic(err)
	}
}
