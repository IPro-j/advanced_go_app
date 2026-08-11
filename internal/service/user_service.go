package service

import (
	"blog-api/internal/model"
	"blog-api/internal/repository"
	"blog-api/pkg/apperr"
	"blog-api/pkg/auth"
	"context"
	"errors"
	"fmt"

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
	if err := ValidateUserCreateRequest(req); err != nil {
		return nil, fmt.Errorf("format failed: %w", err)
	}

	// 3. Проверка уникальности usernameм
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

	return &model.TokenResponse{
		Token:     accessToken,
		ExpiresAt: expiresAt,
		User:      user.ToResponse(),
	}, nil
}

// Login выполняет вход пользователя и возвращает TokenResponse аналогично Register
func (s *UserService) Login(ctx context.Context, req *model.UserLoginRequest) (*model.TokenResponse, error) {

	// Поиск пользователя по email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {

		return nil, err //fmt.Errorf("failed to get user by email: %w", err)
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

// validateUserCreateRequest проверяет корректность данных для регистрации
func ValidateUserCreateRequest(req *model.UserCreateRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	if err := validUserName(req.Username); err != nil {
		return fmt.Errorf("username format failed: %w", err)
	}

	if !isValidEmail(req.Email) {
		return apperr.ErrInvalidEmail
	}

	if err := validPassword(req.Password); err != nil {
		return fmt.Errorf("password format failed: %w", err)
	}

	return nil
}

// GetByID получает пользователя по ID
func (s *UserService) GetByID(ctx context.Context, id int) (*model.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

// GetByEmail получает пользователя по email
func (s *UserService) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}
