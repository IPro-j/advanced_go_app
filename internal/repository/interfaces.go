package repository

import (
	"blog-api/internal/model"
	"context"
)

// UserRepository определяет интерфейс для работы с пользователями
type UserRepository interface {
	// Create создает нового пользователя в базе данных
	Create(ctx context.Context, user *model.User) error

	// GetByID получает пользователя по ID
	GetByID(ctx context.Context, id int) (*model.User, error)

	// GetByEmail получает пользователя по email
	GetByEmail(ctx context.Context, email string) (*model.User, error)

	// GetByUsername получает пользователя по username
	GetByUsername(ctx context.Context, username string) (*model.User, error)

	// ExistsByEmail проверяет существование пользователя по email
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// ExistsByUsername проверяет существование пользователя по username
	ExistsByUsername(ctx context.Context, username string) (bool, error)

	// TODO: Добавить другие методы при необходимости
}

// PostRepository определяет интерфейс для работы с постами
type PostRepository interface {
	// Create создает новый пост
	Create(ctx context.Context, post *model.Post) error

	// GetByID получает пост по ID
	GetByID(ctx context.Context, id int) (*model.Post, error)

	// GetAll получает все посты с пагинацией
	// limit - количество записей на странице
	// offset - смещение от начала
	GetAll(ctx context.Context, limit, offset int) ([]*model.Post, error)

	// GetTotalCount получает общее количество постов
	GetTotalCount(ctx context.Context) (int, error)

	// Методы для фильтрации по автору
	ListByAuthor(ctx context.Context, authorID int, limit, offset int) ([]*model.Post, error)
	CountByAuthor(ctx context.Context, authorID int) (int, error)

	// Update обновляет пост
	Update(ctx context.Context, post *model.Post) error

	// Delete удаляет пост по ID
	Delete(ctx context.Context, id int) error

	// Exists проверяет существование поста по ID
	Exists(ctx context.Context, id int) (bool, error)
}

// CommentRepository определяет интерфейс для работы с комментариями
type CommentRepository interface {
	// Create создает новый комментарий
	Create(ctx context.Context, comment *model.Comment) error

	// GetByID получает комментарий по ID
	GetByID(ctx context.Context, id int) (*model.Comment, error)

	// GetByPostID получает комментарии к посту с пагинацией
	GetByPostID(ctx context.Context, postID int, limit, offset int) ([]*model.Comment, error)

	// GetCountByPostID получает количество комментариев к посту
	GetCountByPostID(ctx context.Context, postID int) (int, error)

	// TODO: Реализовать методы Update и Delete при необходимости

	Update(ctx context.Context, comment *model.Comment) error
	Delete(ctx context.Context, id int) error
	ListByAuthor(ctx context.Context, authorID int, limit, offset int) ([]*model.Comment, error)
	CountByAuthor(ctx context.Context, authorID int) (int, error)
}
