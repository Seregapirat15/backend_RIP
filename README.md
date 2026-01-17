# Go Backend — Система расчёта массы экзопланет

**Лабораторные работы 1-5** | Студент: Номоконов Владислав | Группа: ИУ5-53Б

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

## API

### Публичные
- `GET /api/services` — список инструментов (с фильтрацией)
- `GET /api/services/{id}` — детали инструмента

### С авторизацией
- `POST /api/auth/login` — вход
- `GET /api/orders` — заявки пользователя
- `POST /api/orders/services` — добавить в заявку
- `PUT /api/orders/{id}/form` — сформировать заявку
