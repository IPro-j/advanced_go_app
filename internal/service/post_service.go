package service

import (
	"blog-api/internal/model"
	"blog-api/internal/repository"
	"blog-api/pkg/apperr"
	"context"
	"errors"
	"fmt"
	"strings"
)

type PostService struct {
	postRepo repository.PostRepository
	//userRepo repository.UserRepository
}

func NewPostService(postRepo repository.PostRepository, userRepo repository.UserRepository) *PostService {
	return &PostService{
		postRepo: postRepo,
	}
}

// Create создаёт новый пост от имени пользователя
func (s *PostService) Create(ctx context.Context, userID int, req *model.PostCreateRequest) (*model.Post, error) {
	// 1. Валидация данных
	if err := validatePostCreateRequest(req); err != nil {
		return nil, err
	}

	// 2. Создать модель поста
	post := &model.Post{
		Title:    strings.TrimSpace(req.Title),
		Content:  strings.TrimSpace(req.Content),
		AuthorID: userID,
	}

	// 3. Сохранить через репозиторий
	if err := s.postRepo.Create(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	// 4. Вернуть созданный пост (с заполненным ID из БД)
	return post, nil
}

// GetByID получает пост по ID
func (s *PostService) GetByID(ctx context.Context, id int) (*model.Post, error) {
	post, err := s.postRepo.GetByID(ctx, id)
	if err != nil {

		if errors.Is(err, apperr.ErrPostNotFound) {
			return nil, apperr.ErrPostNotFound
		}
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	return post, nil
}

// GetAll получает список постов с пагинацией
func (s *PostService) GetAll(ctx context.Context, limit, offset int) ([]*model.Post, int, error) {
	// Никакой нормализации здесь! Сервис доверяет, что limit/offset уже валидны.

	posts, err := s.postRepo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list posts: %w", err)
	}

	total, err := s.postRepo.GetTotalCount(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count posts: %w", err)
	}

	return posts, total, nil
}

// Update обновляет пост. Проверяет, что пользователь — автор.
func (s *PostService) Update(ctx context.Context, id int, userID int, req *model.PostUpdateRequest) (*model.Post, error) {
	// 1. Получить существующий пост
	post, err := s.postRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, apperr.ErrPostNotFound) {
			return nil, apperr.ErrPostNotFound
		}
		return nil, fmt.Errorf("failed to get post for update: %w", err)
	}

	// 2. Проверить, что userID является автором
	if !post.CanBeEditedBy(userID) {
		return nil, apperr.ErrForbidden
	}

	// 3. Валидировать новые данные
	if err := validatePostUpdateRequest(req); err != nil {
		return nil, err
	}

	// 4. Обновить только изменённые поля
	if req.Title != "" {
		post.Title = strings.TrimSpace(req.Title)
	}
	if req.Content != "" {
		post.Content = strings.TrimSpace(req.Content)
	}

	// 5. Сохранить через репозиторий
	if err := s.postRepo.Update(ctx, post); err != nil {
		if errors.Is(err, apperr.ErrPostNotFound) {
			return nil, apperr.ErrPostNotFound
		}
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	// 6. Вернуть обновлённый пост
	return post, nil
}

// Delete удаляет пост. Проверяет, что пользователь — автор.
func (s *PostService) Delete(ctx context.Context, id int, userID int) error {
	// 1. Найти пост и проверить существование
	post, err := s.postRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, apperr.ErrPostNotFound) {
			return apperr.ErrPostNotFound
		}
		return fmt.Errorf("failed to get post for delete: %w", err)
	}

	// 2. Проверить, что userID является автором
	if !post.CanBeDeletedBy(userID) {
		return apperr.ErrForbidden
	}

	// 3. Удалить через репозиторий
	if err := s.postRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, apperr.ErrPostNotFound) {
			return apperr.ErrPostNotFound
		}
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}

// validatePostCreateRequest проверяет корректность данных для создания поста
func validatePostCreateRequest(req *model.PostCreateRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	title := strings.TrimSpace(req.Title)
	content := strings.TrimSpace(req.Content)

	if len(title) == 0 {
		return errors.New("title is required")
	}
	if len(title) > 200 {
		return errors.New("title must be no more than 200 characters")
	}
	if len(content) == 0 {
		return errors.New("content is required")
	}

	return nil
}

// validatePostUpdateRequest проверяет корректность данных для обновления поста
func validatePostUpdateRequest(req *model.PostUpdateRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	// Поля опциональны, но если переданы — должны быть валидны
	if req.Title != "" {
		title := strings.TrimSpace(req.Title)
		if len(title) == 0 {
			return errors.New("title cannot be empty if provided")
		}
		if len(title) > 200 {
			return errors.New("title must be no more than 200 characters")
		}
	}

	if req.Content != "" {
		content := strings.TrimSpace(req.Content)
		if len(content) == 0 {
			return errors.New("content cannot be empty if provided")
		}
	}

	return nil
}
