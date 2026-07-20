package repository

import (
	"blog-api/internal/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")
)

// UserRepo представляет репозиторий для работы с пользователями
type UserRepo struct {
	db *sql.DB
}

// NewUserRepo создает новый репозиторий пользователей
func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

// Create создает нового пользователя
func (r *UserRepo) Create(ctx context.Context, user *model.User) error {
	// TODO: Реализовать создание пользователя
	// 1. Подготовить SQL запрос INSERT INTO users...
	// 2. Установить created_at и updated_at = time.Now()
	// 3. Выполнить запрос и получить ID созданной записи
	// 4. Установить ID в структуру user
	//
	// HINT: Используйте QueryRowContext с RETURNING id для получения ID
	// Пример запроса:
	// INSERT INTO users (username, email, password, created_at, updated_at)
	// VALUES ($1, $2, $3, $4, $5) RETURNING id

	query := `
		INSERT INTO users (username, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// TODO: Выполнить запрос и обработать результат

	return fmt.Errorf("not implemented")
}

// GetByID получает пользователя по ID
func (r *UserRepo) GetByID(ctx context.Context, id int) (*model.User, error) {
	// TODO: Реализовать получение пользователя по ID
	// 1. Подготовить SQL запрос SELECT ... FROM users WHERE id = $1
	// 2. Выполнить запрос
	// 3. Просканировать результат в структуру User
	// 4. Обработать случай, когда пользователь не найден (sql.ErrNoRows)
	//
	// HINT: Используйте QueryRowContext и Scan

	query := `
		SELECT id, username, email, password, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user model.User
	// TODO: Выполнить запрос и просканировать результат
	// Не забудьте обработать sql.ErrNoRows и вернуть ErrUserNotFound

	_ = query // Удалите эту строку после реализации
	return nil, fmt.Errorf("not implemented")
}

// GetByEmail получает пользователя по email
func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	// TODO: Реализовать получение пользователя по email
	// Аналогично GetByID, но поиск по полю email

	return nil, fmt.Errorf("not implemented")
}

// GetByUsername получает пользователя по username
func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	// TODO: Реализовать получение пользователя по username
	// Аналогично GetByID, но поиск по полю username

	return nil, fmt.Errorf("not implemented")
}

// ExistsByEmail проверяет существование пользователя по email
func (r *UserRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	// TODO: Реализовать проверку существования пользователя
	// HINT: Используйте SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	// TODO: Выполнить запрос и просканировать результат в переменную exists

	_ = query // Удалите эту строку после реализации
	return false, fmt.Errorf("not implemented")
}

// ExistsByUsername проверяет существование пользователя по username
func (r *UserRepo) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	// TODO: Реализовать проверку существования пользователя по username
	// Аналогично ExistsByEmail

	return false, fmt.Errorf("not implemented")
}

// Update обновляет данные пользователя
func (r *UserRepo) Update(ctx context.Context, user *model.User) error {
	// TODO: (Опционально) Реализовать обновление пользователя
	// 1. Подготовить SQL запрос UPDATE users SET ... WHERE id = $X
	// 2. Обновить updated_at = time.Now()
	// 3. Выполнить запрос
	// 4. Проверить, что запись была обновлена (RowsAffected)

	return fmt.Errorf("not implemented")
}

// Delete удаляет пользователя
func (r *UserRepo) Delete(ctx context.Context, id int) error {
	// TODO: (Опционально) Реализовать удаление пользователя
	// 1. Подготовить SQL запрос DELETE FROM users WHERE id = $1
	// 2. Выполнить запрос
	// 3. Проверить, что запись была удалена (RowsAffected)

	return fmt.Errorf("not implemented")
}
