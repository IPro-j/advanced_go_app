package model

import (
	"time"
)

// User представляет модель пользователя в системе
type User struct {
	ID        int       `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"` // Хешированный пароль, не отдаем в JSON
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Post представляет модель поста в блоге
type Post struct {
	ID        int       `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
	AuthorID  int       `json:"author_id" db:"author_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Comment представляет модель комментария к посту
type Comment struct {
	ID        int       `json:"id" db:"id"`
	Content   string    `json:"content" db:"content"`
	PostID    int       `json:"post_id" db:"post_id"`
	AuthorID  int       `json:"author_id" db:"author_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// UserCreateRequest представляет запрос на создание пользователя
type UserCreateRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// UserLoginRequest представляет запрос на вход пользователя
type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// PostCreateRequest представляет запрос на создание поста
type PostCreateRequest struct {
	Title   string `json:"title" validate:"required,min=1,max=200"`
	Content string `json:"content" validate:"required,min=1"`
}

// PostUpdateRequest представляет запрос на обновление поста
type PostUpdateRequest struct {
	Title   string `json:"title" validate:"required,min=1,max=200"`
	Content string `json:"content" validate:"required,min=1"`
}

// CommentCreateRequest представляет запрос на создание комментария
type CommentCreateRequest struct {
	Content string `json:"content" validate:"required,min=1,max=1000"`
}

// TODO: Добавить следующие структуры и методы:

// UserResponse - структура для ответа с данными пользователя (без пароля)
// Поля: ID, Username, Email, CreatedAt

// TokenResponse - структура для ответа с JWT токеном
// Поля: Token (string), ExpiresAt (time.Time), User (UserResponse)

// PostResponse - структура для ответа с данными поста
// Поля: ID, Title, Content, Author (UserResponse), CreatedAt, UpdatedAt

// CommentResponse - структура для ответа с данными комментария
// Поля: ID, Content, PostID, Author (UserResponse), CreatedAt, UpdatedAt

// TODO: Реализовать методы для моделей:

// User.ToResponse() UserResponse - преобразует User в UserResponse

// Post.CanBeEditedBy(userID int) bool - проверяет, может ли пользователь редактировать пост

// Post.CanBeDeletedBy(userID int) bool - проверяет, может ли пользователь удалить пост

// Comment.CanBeEditedBy(userID int) bool - проверяет, может ли пользователь редактировать комментарий

// Comment.CanBeDeletedBy(userID int) bool - проверяет, может ли пользователь удалить комментарий

// HINT: Пользователь может редактировать/удалять только свои посты и комментарии
// (сравните AuthorID с переданным userID)
