package person

import (
	"encoding/json"
	"fmt"
	"net/http"
	"test_api_go/configs"
	"test_api_go/internal/enrichment"
)

type PersonHandlerDeps struct {
	*configs.Config
	Enricher enrichment.EnrichmentService
}

type PersonHandler struct {
	*configs.Config
	Enricher enrichment.EnrichmentService
}

func NewPersonHandler(router *http.ServeMux, deps PersonHandlerDeps) {
	handler := &PersonHandler{
		Config:   deps.Config,
		Enricher: deps.Enricher,
	}
	router.HandleFunc("POST /person", handler.Create())
	router.HandleFunc("DELETE /person/{id}", handler.Delete())
	router.HandleFunc("GET /person/{id}", handler.GoTo())
	router.HandleFunc("GET /person", handler.GetAll())
	router.HandleFunc("PATCH /person/{id}", handler.Update())
}

func (handler *PersonHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var input PersonCreateRequest

		// парсим json из body
		err := json.NewDecoder(req.Body).Decode(&input)
		if err != nil {
			http.Error(w, "Invalid json", http.StatusBadRequest)
			return
		}

		// явная валидация полей
		switch {
		case input.Name == "":
			http.Error(w, "Field 'name' is required", http.StatusBadRequest)
			return
		case input.Surname == "":
			http.Error(w, "Field 'surname' is required", http.StatusBadRequest)
			return
		}

		// обогащение
		age, err := handler.Enricher.GetAge(input.Name)
		if err != nil {
			http.Error(w, "Faild to enrich age", http.StatusInternalServerError)
			return
		}

		response := PersonResponse{
			ID:         1,
			Name:       input.Name,
			Surname:    input.Surname,
			Patronymic: input.Patronymic,
			Age:        age,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)

	}
}

func (handler *PersonHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("Delete Person")
	}
}

func (handler *PersonHandler) GoTo() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("GoTo Person")
	}
}

func (handler *PersonHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("GetAll Person")
	}
}

func (handler *PersonHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("Update Person")
	}
}
