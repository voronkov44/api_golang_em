# People API Service

REST API сервис для управления людьми с автоматическим обогащением данных (возраст, пол, национальность)

## 🛠 Технологии
- Go 1.24 (чистая архитектура)
- PostgreSQL (хранение данных)
- GORM (ORM для работы с БД)
- Docker (для локального развертывания БД)
- Swagger (документация API)

## Установка и запуск проекта

### 1. Клонирование репозитория
```bash
git clone https://github.com/voronkov44/api_golang_em.git
```
### 2. Настройка окружения 
Требуется установка базы данных [PostgreSQL](https://www.postgresql.org/), если не установлен, смотрите [зависимости](https://github.com/voronkov44/api_golang_em?tab=readme-ov-file#%D0%B7%D0%B0%D0%B2%D0%B8%D1%81%D0%B8%D0%BC%D0%BE%D1%81%D1%82%D0%B8)

Создайте файл .env в корне проекта:

```ini
DSN="host=localhost user=ваш_пользователь password=ваш_пароль dbname=база_данных port=5432 sslmode=disable"
```
### 3. Запуск миграций
```bash
go run migrations/migrate.go
```

### 4. Запуск сервера
```bash
go run cmd/main.go
```
Сервер будет доступен на http://localhost:8081

## 📚 Документация API

После запуска сервера откройте Swagger UI:

```
http://localhost:8081/swagger/index.html

```

Доступные эндпоинты:

- POST /person - создание человека

- GET /person - список всех людей

- GET /person/{id} - получение конкретного человека

- PATCH /person/{id} - обновление данных

- DELETE /person/{id} - удаление

## **Зависимости**
Установка базы данных [PostgreSQL](https://www.postgresql.org/download/)
