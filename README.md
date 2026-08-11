```markdown
# Blog API — REST API для блог‑платформы на Go

## Описание проекта

Проект реализует REST API для блог‑платформы с аутентификацией на JWT, CRUD‑операциями для постов и комментариев, а также контролем прав доступа: только автор может редактировать или удалять свои ресурсы.

**Ключевые особенности:**
- Аутентификация и авторизация через JWT (с жёсткой фиксацией алгоритма, доверяем только HS256 ).
- Хеширование паролей с помощью `bcrypt`.
- Работа с PostgreSQL через `database/`.
- Middleware: структурированное логирование (JSON/text), Recovery для паник, CORS, проверка JWT‑токена.
- Пагинация для списков постов и комментариев.
- Миграции БД и запуск инфраструктуры через `docker-compose`.
- Graceful shutdown сервера с корректным завершением фоновых задач.
- Конфигурирование через переменные окружения, включая таймауты HTTP‑сервера.

---

## Структура проекта

```
blog-api/
├── cmd/api/                          # Точка входа приложения
│   └── main.go
├── internal/                         # Внутренние пакеты приложения
│   ├── model/                        # Модели данных и DTO (User, Post, Comment)
│   ├── handler/                      # HTTP хендлеры (AuthHandler, PostHandler, CommentHandler)
│   ├── service/                      # Бизнес‑логика (UserService, PostService, CommentService)
│   ├── repository/                   # Репозитории (PostRepo, CommentRepo, UserRepo)
│   └── middleware/                   # Middleware (AuthMiddleware, Logger, Recovery)
├── pkg/                              # Переиспользуемые пакеты
│   ├── apperr/                       # Хранилище ошибок приложения
│   ├── auth/                         # JWT и пароли (jwt.go)
│   └── database/                     # Подключение к БД и миграции (postgres.go)
├── docker-compose.yml                # PostgreSQL и Adminer
├── .env.example                      # Пример конфигурации
├── go.mod                            # Зависимости проекта
└── README.md                         # Документация




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

Миграции применяются автоматически при старте приложения:

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

- `POST /api/posts` — создать пост.
- `PATCH /api/posts/{id}` — обновить пост (только автор).
- `DELETE /api/posts/{id}` — удалить пост (только автор).
- `POST /api/posts/{id}/comments` — создать комментарий к посту.
- `PUT /api/comments/{id}` — обновить комментарий (только автор комментария).
- `DELETE /api/comments/{id}` — удалить комментарий (только автор комментария).

> **Важно:** Все защищённые эндпоинты требуют валидного JWT‑токена в заголовке `Authorization`.  

---

## Примеры запросов (curl)

### Регистрация и вход

```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"strongpassword"}'

curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"strongpassword"}'
```

Сохраните `access_token` из ответа для последующих запросов.

### Health check

```bash
curl http://localhost:8080/api/health
```

Ожидаемый ответ:
```json
{"status":"ok","service":"blog-api"}
```

### Посты

```bash
# Создать пост
curl -X POST http://localhost:8080/api/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"title":"My 1st post","content":"Hello world!"}'

# Получить пост
curl -X GET http://localhost:8080/api/posts/1

# Список постов (с пагинацией)
curl -X GET "http://localhost:8080/api/posts?limit=5&offset=0"

# Обновить пост (только автор)
curl -X PATCH "http://localhost:8080/api/posts/1" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"title":"Мой обновлённый пост"}'

# Удалить пост (только автор)
curl -X DELETE "http://localhost:8080/api/posts/2" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Комментарии

```bash
# Создать комментарий к посту
curl -X POST "http://localhost:8080/api/posts/1/comments" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"content":"Отличный пост!"}'

# Получить комментарии к посту (с пагинацией)
curl -X GET "http://localhost:8080/api/posts/1/comments?limit=5&offset=0"

# Обновить комментарий (только автор)
curl -X PUT "http://localhost:8080/api/comments/1" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"content":"Обновлённый текст комментария"}'

# Удалить комментарий (только автор)
curl -v -X DELETE "http://localhost:8080/api/comments/1" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

Для отладки удобно использовать флаг `-v`.

---

## Конфигурация и настройки

Приложение читает настройки из переменных окружения. Ключевые параметры:

- `HTTP_READ_TIMEOUT`, `HTTP_WRITE_TIMEOUT`, `HTTP_IDLE_TIMEOUT` — таймауты сервера.
- `HTTP_SHUTDOWN_GRACE_PERIOD` — время на завершение текущих запросов при остановке.

Пример `.env`:
```env
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=blouser
DB_PASSWORD=blogpassword
DB_NAME=blogdb
DB_SSLMODE=disable
JWT_SECRET=supersecretkey
JWT_EXPIRY_HOURS=24
HTTP_READ_TIMEOUT=15s
HTTP_WRITE_TIMEOUT=15s
HTTP_IDLE_TIMEOUT=60s
HTTP_SHUTDOWN_GRACE_PERIOD=30s
```

---

## Полезные команды

```bash
# Запуск приложения
go run cmd/api/main.go

# База данных
docker compose up -d          # Запустить
docker compose down           # Остановить
docker compose logs -f        # Логи
docker exec -it blog_postgres psql -U postgres -d blogdb
```

---

## Особенности реализации и эксплуатации

- **Пагинация**: параметры `limit` и `offset` нормализуются в сервисе (например, `limit` ограничивается сверху значением 100), а в ответе возвращается фактически применённое значение.
- **Обработка ошибок**: ошибки приложения оборачиваются в доменные типы (`apperr`), что позволяет возвращать корректные HTTP‑статусы и сообщения без раскрытия деталей реализации.
- **Graceful shutdown**: сервер корректно завершает активные запросы и останавливает фоновые процессы , что важно при работе в Docker.

---

## Полезные ссылки

- [Go database/sql tutorial](http://go-database-sql.org/)
- [JWT in Go](https://github.com/golang-jwt/jwt)
- [Chi router](https://github.com/go-chi/chi)
- [bcrypt in Go](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
- [PostgreSQL documentation](https://www.postgresql.org/docs/)
- [Writing tests in Go](https://go.dev/doc/tutorial/testing)
```
