# laba_project

Учебный проект лабораторной работы №4 (командное проектирование и разработка).

Предметная область — микросервис «Витрина отзывов Wildberries»: модели пользователя и отзыва, репозитории CRUD, PostgreSQL, миграции goose.

## Стек

- Go
- PostgreSQL 16 (Docker Compose)
- `database/sql` + `pgx` (драйвер)
- `goose` (миграции)
- `godotenv` (конфигурация из `.env`)

## Быстрый старт

```bash
cp .env.example .env
docker compose up -d

# миграции
goose -dir migrations postgres \
  "host=localhost port=5432 user=postgres password=postgres dbname=laba sslmode=disable" up

# демо CRUD
go run ./cmd
```

## Структура

```
cmd/                 — точка входа
internal/config/     — загрузка .env
internal/model/      — User, Review
internal/repository/ — UserRepository, ReviewRepository
migrations/          — SQL-миграции goose
docker-compose.yaml  — PostgreSQL
```

## Репозиторий

https://github.com/Pabloidd/laba_project
