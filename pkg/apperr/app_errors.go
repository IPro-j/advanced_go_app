package apperr

import "errors"

var (
	ErrForbidden    = errors.New("forbidden")
	ErrUnauthorized = errors.New("unauthorized")
	ErrNotFound     = errors.New("not found")

	// Ошибки сущностей
	ErrCommentNotFound = errors.New("comment not found")
	ErrPostNotFound    = errors.New("post not found")
	ErrPostNotExists   = errors.New("post does not exist")
	ErrUserNotFound    = errors.New("user not found")

	// Ошибки валидации и бизнес‑логики
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrInvalidUsername  = errors.New("invalid username")
	ErrUsernameTooShort = errors.New("user name too short")
	ErrUsernameTooLong  = errors.New("user name too long")
	ErrInvalidEmail     = errors.New("invalid email")

	//ErrWeakPassword     = errors.New("password too weak")
	ErrPasswordTooShort = errors.New("password name too short")
	ErrPasswordTooLong  = errors.New("password name too long")

	ErrNoDigit       = errors.New("must contain at least one digit")
	ErrNoUpper       = errors.New("must contain at least one uppercase latter")
	ErrNoLower       = errors.New("must contain at least one lowrcase latter")
	ErrNoSpecialChar = errors.New("password must contain at least one special char")

	ErrInvalidTokenAlgorithm = errors.New("it's not hs256 token method ")
	ErrInvalidTokenGroup     = errors.New("it's not hmac token group")
)
