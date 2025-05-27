# Цитатник (Quotes Service)

Мини-сервис для хранения и управления цитатами, построенный на Go с использованием принципов чистой архитектуры и PostgreSQL в качестве базы данных.

## Архитектура

Проект следует принципам Clean Architecture и имеет следующую структуру:

```
├── cmd/api/                    # Точка входа в приложение
├── configs/                     # Конфигурация
├── internal/
│   ├── domain/                 # Бизнес-сущности и интерфейсы
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

- **Go 1.24** - основной язык
- **PostgreSQL** - надежная реляционная база данных
- **Стандартная библиотека net/http** - HTTP сервер без внешних зависимостей
- **Docker & Docker Compose** - контейнеризация
- **Clean Architecture** - архитектурный подход

## Быстрый старт

### Требования

- Docker
- Docker Compose

### Запуск

1. Клонируйте репозиторий:
```bash
git clone https://github.com/shoksin/quotes-service
cd quotes-service
```

2. Запустите приложение с помощью Docker Compose:
```bash
docker-compose up --build
```

Приложение будет доступно по адресу `http://localhost:8080`

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

## API Reference

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

## Разработка

### Локальная разработка

1. Запустите PostgreSQL:
```bash
docker run --name quotes_postgres \
  -e POSTGRES_DB=quotes_db \
  -e POSTGRES_USER=quotes_user \
  -e POSTGRES_PASSWORD=quotes_pass \
  -p 5432:5432 \
  -d postgres:15-alpine
```

2. Установите зависимости:
```bash
go mod download
```

3. Запустите приложение:
```bash
go run cmd/api/main.go
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

**Преимущества PostgreSQL:**
- **Надежность**
- **SQL стандарт**
- **Богатые возможности**
- **Масштабируемость**
- **Экосистема**

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

## Особенности реализации

- **Clean Architecture**: четкое разделение слоев
- **Dependency Injection**: зависимости внедряются через интерфейсы
- **Модульная структура**: middleware и инфраструктура в отдельных пакетах
- **Connection Pooling**: настроен пул соединений с БД
- **Graceful shutdown**: корректное завершение работы
- **Health checks**: проверка состояния сервиса
- **Unit tests**: тесты для бизнес-логики
- **Recovery middleware**: обработка паник
- **Структурированное логирование**: все HTTP запросы логируются


## Безопасность

- Использование подготовленных SQL запросов (защита от SQL инъекций)
- Валидация входных данных
- Обработка ошибок без раскрытия внутренней структуры

## Лицензия

MIT License