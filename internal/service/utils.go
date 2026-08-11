package service

import (
	"blog-api/pkg/apperr"
	"regexp"
)

var (
	reDigit    = regexp.MustCompile(`\d`)
	reUpper    = regexp.MustCompile(`[A-Z]`)
	reLower    = regexp.MustCompile(`[a-z]`)
	reSpecial  = regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`)
	emailRegex = regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)
)

const (
	minLenPassword = 8
	maxLenPasword  = 128

	maxLenUsername = 16
	minLenUsername = 3
)

func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// ValidatePassword проверяет требования к паролю
func validPassword(password string) error {
	if len(password) < minLenPassword {
		return apperr.ErrPasswordTooShort
	}
	if len(password) > maxLenPasword {
		return apperr.ErrPasswordTooLong
	}
	if !reDigit.MatchString(password) {
		return apperr.ErrNoDigit
	}
	if !reUpper.MatchString(password) {
		return apperr.ErrNoUpper
	}
	if !reLower.MatchString(password) {
		return apperr.ErrNoLower
	}
	if !reSpecial.MatchString(password) {
		return apperr.ErrNoSpecialChar
	}

	return nil

}

func validUserName(username string) error {
	if len([]rune(username)) < minLenUsername {
		return apperr.ErrUsernameTooShort
	}

	if len([]rune(username)) > maxLenUsername {
		return apperr.ErrUsernameTooLong
	}

	return nil
}
