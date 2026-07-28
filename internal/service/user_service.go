package service

import (
	"blog-api/internal/model"
	"blog-api/internal/repository"
	"blog-api/pkg/apperr"
	"blog-api/pkg/auth"
	"context"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
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

// Register регистрирует нового пользователя и возвращает TokenResponse
func (s *UserService) Register(ctx context.Context, req *model.UserCreateRequest) (*model.TokenResponse, error) {
	// 1. Валидация входных данных
	if err := validateUserCreateRequest(req); err != nil {
		return nil, err
	}

	// 3. Проверка уникальности username
	exists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check username existence: %w", err)
	}
	if exists {
		return nil, apperr.ErrUserAlreadyExists
	}

	// 2. Проверка уникальности email
	exists, err = s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return nil, apperr.ErrEmailAlreadyExists
	}

	// 4. Хеширование пароля
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 5. Создание модели пользователя
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(passwordHash), // в model.go это поле Password, тег json:"-"
	}

	// 6. Сохранение пользователя через репозиторий
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err // репозиторий может вернуть ErrUserAlreadyExists при нарушении уникальности
	}

	// 7. Генерация JWT токенов
	accessToken, expiresAt, err := s.jwtManager.GenerateToken(user.ID, user.Email, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// 8. Подготовка ответа под неизменяемую модель model.TokenResponse:
	// - Token = accessToken
	// - ExpiresAt = время истечения accessToken (если JWTManager не возвращает, можно взять константу из конфига)
	// - User = user.ToResponse()

	return &model.TokenResponse{
		Token:     accessToken,
		ExpiresAt: expiresAt,
		User:      user.ToResponse(),
	}, nil
}

// Login выполняет вход пользователя и возвращает TokenResponse аналогично Register
func (s *UserService) Login(ctx context.Context, req *model.UserLoginRequest) (*model.TokenResponse, error) {
	// 1. Валидация входных данных
	if err := validateUserLoginRequest(req); err != nil {
		return nil, err
	}

	// 2. Поиск пользователя по email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, apperr.ErrUserNotFound) {
			// Не раскрываем, что именно не найдено: email или пользователь
			return nil, apperr.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	// 3. Проверка пароля (bcrypt)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, apperr.ErrInvalidCredentials
	}

	// 4. Генерация JWT токенов
	accessToken, expiresAt, err := s.jwtManager.GenerateToken(user.ID, user.Email, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// 5. Подготовка ответа

	return &model.TokenResponse{
		Token:     accessToken,
		ExpiresAt: expiresAt,
		User:      user.ToResponse(),
	}, nil
}

// GetByID получает пользователя по ID
func (s *UserService) GetByID(ctx context.Context, id int) (*model.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

// GetByEmail получает пользователя по email
func (s *UserService) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

// validateUserCreateRequest проверяет корректность данных для регистрации
func validateUserCreateRequest(req *model.UserCreateRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	username := strings.TrimSpace(req.Username)
	email := strings.TrimSpace(req.Email)
	password := req.Password

	if len(username) < 3 {
		return errors.New("username must be at least 3 characters")
	}

	if len(email) == 0 {
		return errors.New("email is required")
	}

	if !isValidEmail(req.Email) {
		return errors.New("invalid email format")
	}

	if len(password) < 6 {
		return errors.New("password must be at least 6 characters")
	}

	return nil
}

// validateUserLoginRequest проверяет корректность данных для входа
func validateUserLoginRequest(req *model.UserLoginRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	email := strings.TrimSpace(req.Email)
	password := strings.TrimSpace(req.Password)

	if len(email) == 0 {
		return errors.New("email is required")
	}

	if !isValidEmail(req.Email) {
		return errors.New("invalid email format")
	}

	if len(password) == 0 {
		return errors.New("password is required")
	}

	return nil
}
