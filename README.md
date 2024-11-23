# Music Library API

RESTful API сервис для управления музыкальной библиотекой с интеграцией внешнего API.

## Функциональность

- Получение данных библиотеки с фильтрацией и пагинацией
- Получение текста песни
- Удаление песни
- Изменение данных песни
- Добавление новой песни с обогащением данными из внешнего API

## Технологии

- Go 1.20
- Gin Web Framework
- GORM
- PostgreSQL
- Swagger (gin-swagger)
- godotenv

## Установка и запуск

1. Клонируйте репозиторий:
```bash
git clone https://github.com/your-username/music-library.git
cd music-library
```

2. Установите зависимости:
```bash
go mod download
```

3. Создайте базу данных PostgreSQL:
```bash
createdb music_library
```

4. Настройте переменные окружения в файле `.env`:
```
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=postgres
DB_NAME=music_library
MUSIC_API_URL=http://localhost:8081
SERVER_PORT=8080
```

5. Запустите приложение:
```bash
go run cmd/main.go
```

## API Endpoints

- `POST /api/v1/songs` - Создание новой песни
- `GET /api/v1/songs` - Получение списка песен с фильтрацией и пагинацией
- `PUT /api/v1/songs/:id` - Обновление информации о песне
- `DELETE /api/v1/songs/:id` - Удаление песни

## Swagger Documentation

Swagger документация доступна по адресу: `http://localhost:8080/swagger/index.html`

## Структура проекта

```
.
├── cmd/
│   └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── dto/
│   │   └── song.go
│   ├── handler/
│   │   └── song_handler.go
│   ├── model/
│   │   └── song.go
│   ├── repository/
│   │   └── song_repository.go
│   └── service/
│       └── music_api.go
├── .env
├── go.mod
└── README.md
```
