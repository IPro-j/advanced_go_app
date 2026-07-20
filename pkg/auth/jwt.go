package auth

import (
	"errors"
	"time"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token expired")
)

// Claims представляет данные, хранимые в JWT токене
type Claims struct {
	UserID   int    `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	// TODO: Добавить стандартные JWT claims
	// Подсказка: используйте jwt.RegisteredClaims или jwt.StandardClaims
}

// JWTManager управляет созданием и валидацией JWT токенов
type JWTManager struct {
	secretKey []byte
	ttl       time.Duration
}

// NewJWTManager создает новый экземпляр JWT менеджера
func NewJWTManager(secretKey string, ttlHours int) *JWTManager {
	// TODO: Инициализировать JWTManager
	// - Преобразовать secretKey в []byte
	// - Преобразовать ttlHours в time.Duration

	return &JWTManager{}
}

// GenerateToken создает новый JWT токен для пользователя
func (m *JWTManager) GenerateToken(userID int, email, username string) (string, time.Time, error) {
	// TODO: Реализовать генерацию JWT токена
	// Шаги:
	// 1. Создать Claims с данными пользователя
	// 2. Установить время истечения токена (текущее время + ttl)
	// 3. Создать токен используя алгоритм подписи (например, HS256)
	// 4. Подписать токен секретным ключом
	// 5. Вернуть подписанную строку токена и время истечения
	//
	// Подсказка: используйте библиотеку github.com/golang-jwt/jwt/v5

	return "", time.Time{}, errors.New("not implemented")
}

// ValidateToken проверяет и парсит JWT токен
func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	// TODO: Реализовать валидацию и парсинг JWT токена
	// Шаги:
	// 1. Распарсить токен с проверкой подписи
	// 2. Извлечь claims из токена
	// 3. Проверить время истечения токена
	// 4. Вернуть claims если токен валидный
	//
	// Обработать ошибки:
	// - Невалидная подпись -> ErrInvalidToken
	// - Истекший токен -> ErrExpiredToken
	// - Другие ошибки -> ErrInvalidToken

	return nil, errors.New("not implemented")
}

// RefreshToken обновляет существующий токен (опциональное задание)
func (m *JWTManager) RefreshToken(tokenString string) (string, time.Time, error) {
	// TODO: Реализовать обновление токена (продвинутое задание)
	// Шаги:
	// 1. Валидировать существующий токен
	// 2. Извлечь данные пользователя из старого токена
	// 3. Сгенерировать новый токен с теми же данными
	// 4. Вернуть новый токен

	return "", time.Time{}, errors.New("not implemented")
}

// GetUserIDFromToken быстро извлекает ID пользователя из токена без полной валидации
func (m *JWTManager) GetUserIDFromToken(tokenString string) (int, error) {
	// TODO: Извлечь UserID из токена (опциональное задание)
	// Может быть полезно для быстрой проверки

	return 0, errors.New("not implemented")
}
