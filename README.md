# Go Backend — Система расчёта массы экзопланет

**Лабораторные 1–8** | Студент: Номоконов Владислав | Группа: ИУ5-53Б

## Запуск

```bash
docker-compose up --build
```

## Вход модератора

**Логин:** `moderator`  
**Пароль:** `password` (если не подходит, см. ниже)

Если пароль не подходит (в БД мог быть другой хеш), выполните SQL:

```bash
# В Adminer (http://localhost:8080) или psql:
# Выполните: scripts/fix_moderator_password.sql
```

Либо сгенерируйте свой хеш: `go run ./scripts/gen_password.go` → скопируйте хеш → `UPDATE users SET password_hash='...' WHERE login='moderator'`.

## Сервисы

| Сервис | URL |
|--------|-----|
| API | http://localhost:8081 |
| Swagger | http://localhost:8081/swagger/ |
| Adminer | http://localhost:8080 |

## Lab 8 — Межсервисное взаимодействие

- **POST /api/internal/submit-mass-result** — приём результата от async-сервиса (X-Service-Token: Lab8Token)
- **POST /api/admin/orders/{id}/trigger-calculation** — ручной запуск async-расчёта (модератор)
- **FormOrder** — не рассчитывает массу; расчёт в async-сервисе
- **GET заявок** — поля `calculated_count`, `mm_total`

## Что сделано

- **Услуги (инструменты)** — GET список с фильтрацией, GET по id.
- **Авторизация** — POST /auth/login, /auth/register, /auth/logout, GET/PUT /auth/me (JWT + сессии Redis).
- **Заявки** — CRUD, формирование, добавление/удаление/изменение услуг в заявке.
- **PostgreSQL** — пользователи, услуги, заявки, м-м.
- **CORS** — разрешено для GitHub Pages и localhost.

## Что показывать

1. **Swagger** — http://localhost:8081/swagger/ — документация, тест эндпоинтов.
2. **Бэкенд + фронт** — `docker-compose up` → `npm run dev` во фронте → данные в каталоге с API.
3. **Изменение в БД** — Adminer (http://localhost:8080) → правка `instruments` → обновить фронт/Tauri — данные изменились.

## Структура

```
internal/
├── handlers/    # HTTP обработчики
├── database/    # Запросы к PostgreSQL
├── models/      # Структуры данных
└── middleware/  # Auth, CORS
```

## Технологии

- Go 1.21
- PostgreSQL 15
- Redis (сессии)
- MinIO (изображения)
- Docker Compose

## Загрузка изображений

### Через MinIO Console

1. `docker-compose up --build`
2. MinIO: http://localhost:9001 (логин `minioadmin123`, пароль `minioadmin123456`)
3. Bucket `telescope-images` → Upload → в БД (Adminer) прописать `image_url` для инструмента

### Через API (модератор)

`POST /api/services/{id}/image` — multipart/form-data, поле `image`. Swagger: авторизоваться → приложить файл.

## API

- `GET /api/services` — список инструментов (фильтрация)
- `GET /api/services/{id}` — детали инструмента
- `POST /api/auth/login` — вход
- `POST /api/auth/register` — регистрация
- `POST /api/auth/logout` — выход
- `GET /api/auth/me` — текущий пользователь
- `PUT /api/auth/me` — обновить профиль
- `GET /api/orders` — заявки пользователя
- `POST /api/orders/services` — добавить в заявку
- `PUT /api/orders/{id}/form` — сформировать заявку
