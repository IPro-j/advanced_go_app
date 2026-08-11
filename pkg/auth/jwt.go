package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	jwt.RegisteredClaims
}

// JWTManager управляет созданием и валидацией JWT токенов
type JWTManager struct {
	secretKey []byte
	ttl       time.Duration
}

// NewJWTManager создает новый экземпляр JWT менеджера
func NewJWTManager(secretKey string, ttlHours int) (*JWTManager, error) {
	if secretKey == "" {
		return nil, errors.New("secret key cannot be empty")
	}
	if ttlHours <= 0 {
		return nil, errors.New("ttl hours must be greater than 0")
	}

	return &JWTManager{
		secretKey: []byte(secretKey),
		ttl:       time.Duration(ttlHours) * time.Hour,
	}, nil
}

// GenerateToken создает новый JWT токен для пользователя
func (m *JWTManager) GenerateToken(userID int, email, username string) (string, time.Time, error) {
	expirationTime := time.Now().UTC().Add(m.ttl)

	claims := &Claims{
		UserID:   userID,
		Email:    email,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(m.secretKey)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign JWT token: %w", err)
	}

	return tokenString, expirationTime, nil
}

// ValidateToken проверяет и парсит JWT токен
func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	if tokenString == "" {
		return nil, fmt.Errorf("empty token: %w", ErrInvalidToken)
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// ЖЁСТКАЯ ПРОВЕРКА АЛГОРИТМА: доверяем только HS256
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidToken
		}

		return m.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token expired: %w", ErrExpiredToken)
		}
		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, fmt.Errorf("token not valid yet: %w", ErrInvalidToken)
		}
		if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return nil, fmt.Errorf("invalid signature: %w", ErrInvalidToken)
		}
		// Другие ошибки (неверный формат и т.п.)
		return nil, fmt.Errorf("failed to validate token: %w", ErrInvalidToken)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid claims or token not valid: %w", ErrInvalidToken)
	}

	return claims, nil
}

// RefreshToken обновляет существующий токен
func (m *JWTManager) RefreshToken(tokenString string) (string, time.Time, error) {
	// 1. Валидируем старый токен
	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		return "", time.Time{}, err // возвращаем ту же ошибку валидации
	}

	// 2. Извлекаем данные пользователя из старого токена
	userID := claims.UserID
	username := claims.Username
	email := claims.Email

	// 3. Генерируем новый токен с теми же данными, но новым временем истечения
	newToken, newExp, err := m.GenerateToken(userID, email, username)
	if err != nil {
		return "", time.Time{}, err
	}

	// 4. Возвращаем новый токен и время истечения
	return newToken, newExp, nil
}

// GetUserIDFromToken быстро извлекает ID пользователя из токена без полной валидации
func (m *JWTManager) GetUserIDFromToken(tokenString string) (int, error) {
	claims, err := m.ValidateToken(tokenString) // <-- исправлено: было Validate
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}
