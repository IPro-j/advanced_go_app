```markdown
# Blog API — REST API для блог‑платформы на Go

## Описание проекта

Проект реализует REST API для блог‑платформы с аутентификацией на JWT, CRUD‑операциями для постов и комментариев, а также контролем прав доступа: только автор может редактировать или удалять свои ресурсы.

**Ключевые особенности:**
- Аутентификация и авторизация через JWT (с жёсткой фиксацией алгоритма, без доверия к полю `alg` из токена).
- Хеширование паролей с помощью `bcrypt`.
- Работа с PostgreSQL через `database/sql` и пул соединений.
- Middleware: логирование запросов, Recovery для паник, CORS, проверка JWT‑токена.
- Пагинация для списков постов и комментариев.
- Миграции БД и запуск инфраструктуры через `docker-compose`.

---

## Структура проекта

```
blog-api/
├── cmd/api/                          # Точка входа приложения
│   └── main.go
├── internal/                         # Внутренние пакеты приложения
│   ├── model/                        # Модели данных и DTO (User, Post, Comment)
│   ├── handler/                      # HTTP хендлеры (AuthHandler, PostHandler)
│   ├── service/                      # Бизнес‑логика (UserService, PostService)
│   ├── repository/                   # Репозитории (PostRepo, CommentRepo, UserRepo)
│   └── middleware/                   # Middleware (AuthMiddleware, Logger, Recovery)
├── pkg/                              # Переиспользуемые пакеты
│   ├── apperr/                       # Хранилище ошибок приложения
│   ├── auth/                         # JWT и пароли (jwt.go, password.go)
│   └── database/                     # Подключение к БД и миграции (postgres.go)
├── migrations/                       # SQL миграции
├── docker-compose.yml                # PostgreSQL и Adminer
├── .env.example                      # Пример конфигурации
├── go.mod                            # Зависимости проекта
└── README.md                         # Документация

```

---

## Начало работы

### 1. Подготовка окружения

```bash
# Установить зависимости
go mod download

# Создать файл конфигурации
cp .env.example .env

# Запустить PostgreSQL и Adminer
docker compose up -d

# Проверить, что БД работает
docker compose logs postgres

```

### 2. Применение миграций


миграции применяются автоматически при старте:

```bash
go run cmd/api/main.go
```

Проверьте наличие таблиц `users`, `posts`, `comments` (можно через Adminer или `psql`).

---

## API Эндпоинты

### Публичные (без аутентификации)

- `POST /api/register` — регистрация пользователя.
- `POST /api/login` — вход пользователя, возврат JWT‑токена.
- `GET /api/posts` — список постов (пагинация: `?limit=10&offset=0`).
- `GET /api/posts/{id}` — получить пост по ID.
- `GET /api/posts/{id}/comments` — комментарии к посту (пагинация).
- `GET /health` — health check (без префикса `/api`).

### Защищённые (требуется заголовок `Authorization: Bearer <token>`)

- `POST /api/protected/posts` — создать пост.
- `PUT /api/protected/posts/{id}` — обновить пост (только автор).
- `DELETE /api/protected/posts/{id}` — удалить пост (только автор).
- `POST /api/protected/posts/{id}/comments` — создать комментарий к посту.
- `PUT /api/protected/comments/{id}` — обновить комментарий (только автор комментария).
- `DELETE /api/protected/comments/{id}` — удалить комментарий (только автор комментария).

> **Важно:** Все защищённые эндпоинты требуют валидного JWT‑токена в заголовке `Authorization`.

---

## Примеры запросов (curl для Windows)

### Регистрация и вход

```bush
curl -X POST http://localhost:8080/api/register -H "Content-Type: application/json" -d "{\"username\":\"testuser\",\"email\":\"test@example.com\",\"password\":\"strongpassword\"}"

curl -X POST http://localhost:8080/api/login -H "Content-Type: application/json" -d "{\"email\":\"test@example.com\",\"password\":\"strongpassword\"}"
```

Сохраните `access_token` из ответа для последующих запросов.


### Health check

```bush
curl -v http://localhost:8080/health
```

Ожидаемый ответ:
```json
{"status":"ok","service":"blog-api"}
```

### Посты

```bash
# Создать пост
curl -X POST http://localhost:8080/api/protected/posts -H "Content-Type: application/json" -H "Authorization: Bearer YOUR_JWT_TOKEN"" -d "{\"title\":\"My 1st post\",\"content\":\"Hello world!\"}"

# Получить пост
curl -X GET http://localhost:8080/api/posts/1

# Список постов
curl -X GET http://localhost:8080/api/posts/


# Обновить пост (только автор)
curl -X PUT "http://localhost:8080/api/protected/posts/1" -H "Content-Type: application/json" -H "Authorization: Bearer OUR_JWT_TOKEN" -d "{\"title\":\"Мой 2 пост\",\"content\":\"Теперь тут другой текст\"}"

# Удалить пост (только автор)
curl -X DELETE "http://localhost:8080/api/protected/posts/2" -H "Authorization: Bearer OUR_JWT_TOKEN""
```

### Комментарии

```bash
# Создать комментарий к посту
curl -X POST http://localhost:8080/api/protected/posts/1/comments -H "Content-Type: application/json" -H "Authorization: Bearer YOUR_JWT_TOKEN" -d "{\"content\":\"Ерунда какая-та\"}"


# Получить комментарии к посту (с пагинацией)
curl -X GET "http://localhost:8080/api/posts/1/comments?limit=5&offset=0"

# Обновить комментарий (только автор)
curl -X PUT http://localhost:8080/api/protected/comments/1 -H "Content-Type: application/json" -H "Authorization: Bearer YOUR_JWT_TOKEN" -d "{\"content\":\"Обновлённый текст комментария\"}"


# Удалить комментарий (только автор)
curl -v -X DELETE http://localhost:8080/api/protected/comments/1 -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

Для отладки удобно использовать `-v`


---

## Полезные команды

```bash
# Запуск приложения
go run cmd/api/main.go


# База данных
docker compose up -d          # Запустить
docker compose down           # Остановить
docker compose logs -f        # Логи
docker exec -it blog-api_postgres_1 psql -U postgres -d blog_db
```



## Полезные ссылки

- [Go database/sql tutorial](http://go-database-sql.org/)
- [JWT in Go](https://github.com/golang-jwt/jwt)
- [Chi router](https://github.com/go-chi/chi)
- [bcrypt in Go](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
- [PostgreSQL documentation](https://www.postgresql.org/docs/)
- [Writing tests in Go](https://go.dev/doc/tutorial/testing)

---
