# Go Backend — Система расчёта массы экзопланет

**Лабораторные 1–7** | Студент: Номоконов Владислав | Группа: ИУ5-53Б

## Запуск

```bash
docker-compose up --build
```

## Сервисы

| Сервис | URL |
|--------|-----|
| API | http://localhost:8081 |
| Swagger | http://localhost:8081/swagger/ |
| Adminer | http://localhost:8080 |

## Что сделано

- **Услуги (инструменты)** — GET список с фильтрацией, GET по id.
- **Авторизация** — POST /auth/login, /auth/register, /auth/logout, GET/PUT /auth/me (JWT + сессии Redis).
- **Заявки** — CRUD, формирование, добавление/удаление/изменение услуг в заявке, расчёт массы.
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
