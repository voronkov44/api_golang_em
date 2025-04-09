package person

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"test_api_go/configs"
	"test_api_go/internal/enrichment"
)

type PersonHandlerDeps struct {
	*configs.Config
	Enricher enrichment.EnrichmentService
	DB       *gorm.DB
}

type PersonHandler struct {
	*configs.Config
	Enricher enrichment.EnrichmentService
	DB       *gorm.DB
}

func NewPersonHandler(router *http.ServeMux, deps PersonHandlerDeps) {
	handler := &PersonHandler{
		Config:   deps.Config,
		Enricher: deps.Enricher,
		DB:       deps.DB,
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
			http.Error(w, "Failed to enrich age", http.StatusInternalServerError)
			return
		}
		gender, err := handler.Enricher.GetGender(input.Name)
		if err != nil {
			http.Error(w, "Failed to enrich gender", http.StatusInternalServerError)
			return
		}
		nationality, err := handler.Enricher.GetNationality(input.Name)
		if err != nil {
			http.Error(w, "Failed to enrich nationality", http.StatusInternalServerError)
			return
		}

		personModel := Person{
			Name:        input.Name,
			Surname:     input.Surname,
			Patronymic:  input.Patronymic,
			Age:         age,
			Gender:      gender,
			Nationality: nationality,
		}

		if err := handler.DB.Create(&personModel).Error; err != nil {
			http.Error(w, "Failed to save person", http.StatusInternalServerError)
			return
		}

		response := PersonResponse{
			ID:          personModel.ID,
			Name:        personModel.Name,
			Surname:     personModel.Surname,
			Patronymic:  personModel.Patronymic,
			Age:         personModel.Age,
			Gender:      personModel.Gender,
			Nationality: personModel.Nationality,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)

	}
}

func (handler *PersonHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		idStr := req.PathValue("id")
		if idStr == "" {
			http.Error(w, "Invalid id", http.StatusBadRequest)
			return
		}

		var personId int64
		if _, err := fmt.Sscanf(idStr, "%d", &personId); err != nil {
			http.Error(w, "Invalid id", http.StatusBadRequest)
			return
		}

		result := handler.DB.Delete(&Person{}, personId)
		if result.Error != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		if result.RowsAffected == 0 {
			http.Error(w, "Person not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (handler *PersonHandler) GoTo() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		idStr := req.PathValue("id")
		if idStr == "" {
			http.Error(w, "ID is required", http.StatusBadRequest)
			return
		}

		var personId int64
		if _, err := fmt.Sscanf(idStr, "%d", &personId); err != nil {
			http.Error(w, "Invalid ID format", http.StatusBadRequest)
			return
		}

		var person Person
		if err := handler.DB.First(&person, personId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, "Person not found", http.StatusNotFound)
			} else {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			return
		}

		response := PersonResponse{
			ID:          person.ID,
			Name:        person.Name,
			Surname:     person.Surname,
			Patronymic:  person.Patronymic,
			Age:         person.Age,
			Gender:      person.Gender,
			Nationality: person.Nationality,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
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
