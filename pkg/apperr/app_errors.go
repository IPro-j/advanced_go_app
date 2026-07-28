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
)
