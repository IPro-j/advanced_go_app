package auth

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmptyPassword    = errors.New("password cannot be empty")
	ErrPasswordTooShort = errors.New("password is too short")
)

const (
	defaultCost = 10
	minLength   = 6
)

// HashPassword хеширует пароль используя bcrypt
func HashPassword(password string) (string, error) {
	// 1. Проверка на пустоту
	if strings.TrimSpace(password) == "" {
		return "", ErrEmptyPassword
	}

	// 2. Хеширование с заданным cost factor
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), defaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedBytes), nil
}

// CheckPassword проверяет соответствие пароля и его хеша
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == nil {
		return true
	}
	// Если ошибка — считаем, что пароли не совпадают
	return false
}

// ValidatePasswordStrength проверяет надежность пароля
func ValidatePasswordStrength(password string) error {
	if strings.TrimSpace(password) == "" {
		return ErrEmptyPassword
	}

	if len(password) < minLength {
		return ErrPasswordTooShort
	}

	hasLower := false
	hasUpper := false
	hasDigit := false

	for _, r := range password {
		switch {
		case r >= 'a' && r <= 'z':
			hasLower = true
		case r >= 'A' && r <= 'Z':
			hasUpper = true
		case r >= '0' && r <= '9':
			hasDigit = true
		}
	}

	// Опционально: требуем хотя бы 2 из 3 категорий (строчные, заглавные, цифры)
	count := 0
	if hasLower {
		count++
	}
	if hasUpper {
		count++
	}
	if hasDigit {
		count++
	}

	if count < 2 {
		return errors.New("password must contain at least two of the following: lowercase letters, uppercase letters, digits")
	}

	return nil
}

// GenerateRandomPassword генерирует случайный пароль
func GenerateRandomPassword(length int) (string, error) {
	if length <= 0 {
		return "", errors.New("length must be positive")
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	setLen := len(charset)

	buf := make([]byte, length)
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(setLen)))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		buf[i] = charset[n.Int64()]
	}

	return string(buf), nil
}
