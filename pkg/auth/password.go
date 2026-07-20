package auth

import (
	"errors"
)

var (
	ErrEmptyPassword    = errors.New("password cannot be empty")
	ErrPasswordTooShort = errors.New("password is too short")
)

// HashPassword хеширует пароль используя bcrypt
func HashPassword(password string) (string, error) {
	// TODO: Реализовать хеширование пароля
	// Шаги:
	// 1. Проверить что пароль не пустой
	// 2. Использовать bcrypt для хеширования
	// 3. Выбрать подходящий cost factor (например, 10-12)
	// 4. Вернуть хешированный пароль как строку
	//
	// Подсказка: используйте golang.org/x/crypto/bcrypt

	return "", errors.New("not implemented")
}

// CheckPassword проверяет соответствие пароля и его хеша
func CheckPassword(password, hash string) bool {
	// TODO: Реализовать проверку пароля
	// Шаги:
	// 1. Сравнить пароль с хешом используя bcrypt
	// 2. Вернуть true если пароль совпадает, false если нет
	// 3. При ошибке вернуть false
	//
	// Подсказка: bcrypt.CompareHashAndPassword

	return false
}

// ValidatePasswordStrength проверяет надежность пароля
func ValidatePasswordStrength(password string) error {
	// TODO: Реализовать проверку надежности пароля
	// Требования:
	// - Минимум 6 символов
	// - Опционально: содержит буквы и цифры
	// - Опционально: содержит заглавные и строчные буквы
	//
	// Вернуть соответствующую ошибку или nil

	return errors.New("not implemented")
}

// GenerateRandomPassword генерирует случайный пароль (опциональное задание)
func GenerateRandomPassword(length int) (string, error) {
	// TODO: Реализовать генерацию случайного пароля
	// Шаги:
	// 1. Создать набор допустимых символов
	// 2. Сгенерировать случайную последовательность заданной длины
	// 3. Вернуть пароль как строку
	//
	// Подсказка: используйте crypto/rand для криптографически стойкой генерации

	return "", errors.New("not implemented")
}
