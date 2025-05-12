package person

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"math"
	"net/http"
	"strconv"
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

// Create godoc
// @Summary Создать новую запись о человеке
// @Description Добавляет человека в БД с обогащением данных (возраст, пол, национальность)
// @Tags people
// @Accept json
// @Produce json
// @Param input body person.PersonCreateRequest true "Данные человека"
// @Success 201 {object} person.PersonResponse
// @Failure 400 {string} string "Неверный JSON"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /person [post]
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

// Delete godoc
// @Summary Удалить человека
// @Description Удаляет запись о человеке из БД
// @Tags people
// @Param id path int true "ID человека"
// @Success 200
// @Failure 400 {string} string "Неверный ID"
// @Failure 404 {string} string "Человек не найден"
// @Router /person/{id} [delete]
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

// GoTo godoc
// @Summary Получить человека по ID
// @Description Возвращает данные конкретного человека
// @Tags people
// @Produce json
// @Param id path int true "ID человека"
// @Success 200 {object} person.PersonResponse
// @Failure 400 {string} string "Неверный ID"
// @Failure 404 {string} string "Человек не найден"
// @Router /person/{id} [get]
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

// GetAll godoc
// @Summary Получить список людей
// @Description Возвращает список людей с пагинацией и фильтрацией
// @Tags people
// @Produce json
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Лимит записей" default(10)
// @Param name query string false "Фильтр по имени"
// @Param surname query string false "Фильтр по фамилии"
// @Param age query int false "Фильтр по возрасту"
// @Param gender query string false "Фильтр по полу"
// @Param nationality query string false "Фильтр по национальности"
// @Success 200 {object} map[string]interface{} "data: список людей, meta: метаданные пагинации"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /person [get]
func (handler *PersonHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		page, _ := strconv.Atoi(req.URL.Query().Get("page"))
		if page <= 0 {
			page = 1
		}

		limit, _ := strconv.Atoi(req.URL.Query().Get("limit"))
		switch {
		case limit > 100:
			limit = 100
		case limit <= 0:
			limit = 10
		}

		offset := (page - 1) * limit

		name := req.URL.Query().Get("name")
		surname := req.URL.Query().Get("surname")
		age, _ := strconv.Atoi(req.URL.Query().Get("age"))
		gender := req.URL.Query().Get("gender")
		nationality := req.URL.Query().Get("nationality")

		query := handler.DB.Model(&Person{})

		if name != "" {
			query = query.Where("name LIKE ?", "%"+name+"%")
		}
		if surname != "" {
			query = query.Where("surname LIKE ?", "%"+surname+"%")
		}
		if age > 0 {
			query = query.Where("age = ?", age)
		}
		if gender != "" {
			query = query.Where("gender = ?", gender)
		}
		if nationality != "" {
			query = query.Where("nationality = ?", nationality)
		}

		var total int64
		if err := query.Count(&total).Error; err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		var persons []Person
		if err := query.Offset(offset).Limit(limit).Find(&persons).Error; err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		response := make([]PersonResponse, len(persons))
		for i, person := range persons {
			response[i] = PersonResponse{
				ID:          person.ID,
				Name:        person.Name,
				Surname:     person.Surname,
				Patronymic:  person.Patronymic,
				Age:         person.Age,
				Gender:      person.Gender,
				Nationality: person.Nationality,
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": response,
			"meta": map[string]interface{}{
				"total":     total,
				"page":      page,
				"limit":     limit,
				"last_page": int(math.Ceil(float64(total) / float64(limit))),
			},
		})
	}
}

// Update godoc
// @Summary Обновить данные человека
// @Description Изменяет данные человека (частичное обновление)
// @Tags people
// @Accept json
// @Produce json
// @Param id path int true "ID человека"
// @Param input body person.PersonUpdateRequest true "Новые данные"
// @Success 200 {object} person.PersonResponse
// @Failure 400 {string} string "Неверный запрос"
// @Failure 404 {string} string "Человек не найден"
// @Router /person/{id} [patch]
func (handler *PersonHandler) Update() http.HandlerFunc {
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

		var input PersonUpdateRequest
		if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
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

		if input.Name != "" {
			person.Name = input.Name
		}
		if input.Surname != "" {
			person.Surname = input.Surname
		}
		if input.Patronymic != "" {
			person.Patronymic = input.Patronymic
		}
		if input.Age > 0 {
			person.Age = input.Age
		}
		if input.Gender != "" {
			person.Gender = input.Gender
		}
		if input.Nationality != "" {
			person.Nationality = input.Nationality
		}

		if err := handler.DB.Save(&person).Error; err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
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
