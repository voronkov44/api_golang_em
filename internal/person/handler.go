package person

import (
	"encoding/json"
	"fmt"
	"net/http"
	"test_api_go/configs"
)

type PersonHandlerDeps struct {
	*configs.Config
}

type PersonHandler struct {
	*configs.Config
}

func NewPersonHandler(router *http.ServeMux, deps PersonHandlerDeps) {
	handler := &PersonHandler{
		Config: deps.Config,
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

		// валидация поля name и surname (patronymic не будем валидировать)
		if input.Name == "" || input.Surname == "" {
			http.Error(w, "Name and surname must be", http.StatusBadRequest)
			return
		}

		response := PersonResponse{
			ID:         1,
			Name:       input.Name,
			Surname:    input.Surname,
			Patronymic: input.Patronymic,
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
