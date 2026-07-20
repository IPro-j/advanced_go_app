package service

import (
	"blog-api/internal/model"
	"blog-api/internal/repository"
	"blog-api/pkg/auth"
	"context"
	"errors"
	"fmt"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
)

type UserService struct {
	userRepo   repository.UserRepository
	jwtManager *auth.JWTManager
}

func NewUserService(userRepo repository.UserRepository, jwtManager *auth.JWTManager) *UserService {
	return &UserService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

func (s *UserService) Register(ctx context.Context, req *model.UserCreateRequest) (*model.TokenResponse, error) {
	// TODO: Реализовать регистрацию пользователя
	// Шаги:
	// 1. Валидация входных данных (username >= 3 символов, email валидный, пароль >= 6 символов)
	// 2. Проверить уникальность email через репозиторий
	// 3. Проверить уникальность username через репозиторий
	// 4. Захешировать пароль используя пакет auth
	// 5. Создать модель пользователя с хешированным паролем
	// 6. Сохранить пользователя через репозиторий
	// 7. Сгенерировать JWT токен для нового пользователя
	// 8. Вернуть TokenResponse с токеном и данными пользователя

	return nil, fmt.Errorf("not implemented")
}

func (s *UserService) Login(ctx context.Context, req *model.UserLoginRequest) (*model.TokenResponse, error) {
	// TODO: Реализовать вход пользователя
	// Шаги:
	// 1. Валидация входных данных
	// 2. Найти пользователя по email через репозиторий
	// 3. Проверить пароль используя функцию из пакета auth
	// 4. Сгенерировать JWT токен при успешной аутентификации
	// 5. Вернуть TokenResponse
	// ВАЖНО: При ошибке не раскрывать, что именно неправильно (email или пароль)

	return nil, fmt.Errorf("not implemented")
}

func (s *UserService) GetByID(ctx context.Context, id int) (*model.User, error) {
	// TODO: Получить пользователя по ID через репозиторий

	return nil, fmt.Errorf("not implemented")
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	// TODO: Получить пользователя по email через репозиторий

	return nil, fmt.Errorf("not implemented")
}

// validateUserCreateRequest проверяет корректность данных для регистрации
func validateUserCreateRequest(req *model.UserCreateRequest) error {
	// TODO: Реализовать проверку всех полей

	return nil
}

// validateUserLoginRequest проверяет корректность данных для входа
func validateUserLoginRequest(req *model.UserLoginRequest) error {
	// TODO: Реализовать проверку полей

	return nil
}
