# Цитатник (Quotes Service)

Мини-сервис для хранения и управления цитатами, построенный на Go с использованием принципов чистой архитектуры и PostgreSQL в качестве базы данных.

## Быстрый старт

### Требования

- Docker
- Docker Compose

### Запуск:

1. Клонируйте репозиторий:
```bash
git clone https://github.com/shoksin/quotes-service
cd quotes-service
```

2. Скопируйте переменные окружения:
```
cp .env.example .env
```

### 1-ый способ: docker-compose:
3. Запустите приложение с помощью docker-compose:
```bash
docker-compose up --build
```

### 2-ой способ Docker:
3. Запустите базу данных c помощью Docker:
```bash
docker run --name quotes_postgres \
  -e POSTGRES_DB=quotes_db \
  -e POSTGRES_USER=quotes_user \
  -e POSTGRES_PASSWORD=quotes_pass \
  -p 5432:5432 \
  -d postgres:15-alpine
```
4. Установите зависимости:
```bash
go mod download
```

5. Запустите программу:
```bash
go run cmd/api/main.go
```

### Приложение будет доступно по адресу `http://localhost:8080`

## Архитектура

Проект следует принципам Clean Architecture и имеет следующую структуру:

```
├── cmd/api/                    # Точка входа в приложение
├── configs/                    # Конфигурация
├── internal/
│   ├── domain/                 # Бизнес-сущности
│   ├── usecase/                # Бизнес-логика
│   ├── repository/             # Слой доступа к данным (PostgreSQL)
│   ├── delivery/http/          # HTTP handlers и middleware
│   │   └── middleware/         # HTTP middleware
│   └── storage/                # Подключение к базе данных
├── migrations/                 # SQL миграции
├── docker-compose.yml
├── Dockerfile
└── README.md
```

## Функциональность

- Добавление новой цитаты (POST /quotes)
- Получение всех цитат (GET /quotes)
- Получение случайной цитаты (GET /quotes/random)
- Фильтрация по автору (GET /quotes?author=Confucius)
- Удаление цитаты по ID (DELETE /quotes/{id})
- Health Check (GET /health)

## Технологии

- **Go 1.24**
- **PostgreSQL**
- **Стандартная библиотека net/http**
- **Docker & Docker Compose**

### Проверка работы

Используйте следующие curl команды для тестирования API:

#### 1. Добавление новой цитаты
```bash
curl -X POST http://localhost:8080/quotes \
  -H "Content-Type: application/json" \
  -d '{"author":"Confucius", "quote":"Life is simple, but we insist on making it complicated."}'
```

#### 2. Получение всех цитат
```bash
curl http://localhost:8080/quotes
```

#### 3. Получение случайной цитаты
```bash
curl http://localhost:8080/quotes/random
```

#### 4. Фильтрация по автору
```bash
curl "http://localhost:8080/quotes?author=Confucius"
```

#### 5. Удаление цитаты
```bash
curl -X DELETE http://localhost:8080/quotes/1
```

#### 6. Health Check
```bash
curl http://localhost:8080/health
```

## Запросы:

### POST /quotes
Создание новой цитаты

**Request Body:**
```json
{
  "author": "string",
  "quote": "string"
}
```

**Response:**
```json
{
  "id": 1,
  "author": "Confucius",
  "quote": "Life is simple, but we insist on making it complicated.",
  "created_at": "2023-12-07T10:30:00Z"
}
```

### GET /quotes
Получение всех цитат с возможностью фильтрации

**Query Parameters:**
- `author` (optional) - фильтр по автору

**Response:**
```json
[
  {
    "id": 1,
    "author": "Confucius",
    "quote": "Life is simple, but we insist on making it complicated.",
    "created_at": "2023-12-07T10:30:00Z"
  }
]
```

### GET /quotes/random
Получение случайной цитаты

**Response:**
```json
{
  "id": 1,
  "author": "Confucius",
  "quote": "Life is simple, but we insist on making it complicated.",
  "created_at": "2023-12-07T10:30:00Z"
}
```

### DELETE /quotes/{id}
Удаление цитаты по ID

**Response:** 204 No Content

### GET /health
Health Check endpoint

**Response:**
```json
{
  "status": "healthy",
  "service": "quotes-service"
}
```

### Тестирование

Запуск unit-тестов:
```bash
go test ./...
```

Запуск тестов с покрытием:
```bash
go test -cover ./...
```

### Переменные окружения

Скопируйте `.env.example` в `.env` и измените значения по необходимости:

```bash
cp .env.example .env
```

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| SERVER_PORT | Порт HTTP сервера | 8080 |
| DB_HOST | Хост PostgreSQL | localhost |
| DB_PORT | Порт PostgreSQL | 5432 |
| DB_NAME | Имя базы данных | quotes_db |
| DB_USER | Пользователь БД | quotes_user |
| DB_PASSWORD | Пароль БД | quotes_pass |

## Структура базы данных

```sql
CREATE TABLE quotes (
    id SERIAL PRIMARY KEY,
    author VARCHAR(255) NOT NULL,
    quote TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

## Docker

### Сборка образа
```bash
docker build -t quotes-service .
```

### Запуск контейнера
```bash
docker run -p 8080:8080 \
  -e DB_HOST=localhost \
  -e DB_PORT=5432 \
  -e DB_NAME=quotes_db \
  -e DB_USER=quotes_user \
  -e DB_PASSWORD=quotes_pass \
  quotes-service
```

## Логирование

Приложение логирует:
- HTTP запросы с методом, путем, статус кодом и временем выполнения
- Подключение к базе данных
- Ошибки выполнения

## Обработка ошибок

API возвращает структурированные ошибки в формате:
```json
{
  "error": "error message"
}
```

Возможные HTTP статус коды:
- `200` - Успешный запрос
- `201` - Ресурс создан
- `204` - Ресурс удален
- `400` - Неверный запрос
- `404` - Ресурс не найден
- `500` - Внутренняя ошибка сервера

## Безопасность

- Использование подготовленных SQL запросов (защита от SQL инъекций)
- Валидация входных данных
- Обработка ошибок без раскрытия внутренней структуры

## Лицензия

MIT License