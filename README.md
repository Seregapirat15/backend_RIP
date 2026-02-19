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

## Загрузка фотографий (MinIO)

### Через веб-интерфейс MinIO

1. Запустите бэкенд: `docker-compose up --build`.
2. Откройте **MinIO Console**: http://localhost:9001  
   Логин: `minioadmin123`  
   Пароль: `minioadmin123456`
3. Слева выберите **Buckets** → откройте бакет **telescope-images** (если его нет — создайте: **Create Bucket** → имя `telescope-images`).
4. В бакете нажмите **Upload** → **Upload file** и выберите файл (JPG/PNG). Имя файла должно быть **на латинице** (например `harps.jpg`, `service_1_123.jpg`).
5. Чтобы картинка отображалась у инструмента, в БД нужно прописать **точное** имя файла (как в MinIO, с учётом регистра) в поле `image_url`:
   - откройте **Adminer**: http://localhost:8080 (сервер: `postgres`, пользователь: `postgres`, пароль: `postgres123`, БД: `exoplanet_calculations`);
   - выполните SQL, подставив свой файл и id инструмента:
   ```sql
   UPDATE instruments SET image_url = 'harps.jpg' WHERE id = 1;
   ```
   Либо запустите готовый скрипт для стандартных имён (harps.jpg, jameswebb.jpg, espresso.jpg, SPIRou.JPG): в Adminer → SQL Command → вставьте содержимое файла `scripts/fix_instrument_images.sql` и выполните.

### Через API (для модератора)

1. Запустите бэкенд с Docker (MinIO поднимется на 9000/9001).
2. Загрузка только для **модератора**: `POST /api/services/{id}/image`, тело — `multipart/form-data`, поле `image` (файл JPG/PNG).
3. В **Swagger** (http://localhost:8081/swagger/): авторизуйтесь под модератором, откройте `POST /api/services/{id}/image`, укажите `id` инструмента и приложите файл в поле `image`.
4. Изображения сохраняются в MinIO (бакет `telescope-images`). В ответах API подставляются presigned URL — фронт подгружает фото по ним автоматически.

## API

### Публичные
- `GET /api/services` — список инструментов (с фильтрацией)
- `GET /api/services/{id}` — детали инструмента

### С авторизацией
- `POST /api/auth/login` — вход
- `GET /api/orders` — заявки пользователя
- `POST /api/orders/services` — добавить в заявку
- `PUT /api/orders/{id}/form` — сформировать заявку
