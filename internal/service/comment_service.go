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

type CommentService struct {
	commentRepo repository.CommentRepository
	postRepo    repository.PostRepository
	userRepo    repository.UserRepository
}

func NewCommentService(
	commentRepo repository.CommentRepository,
	postRepo repository.PostRepository,
	userRepo repository.UserRepository,
) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		postRepo:    postRepo,
		userRepo:    userRepo,
	}
}

// Create создаёт новый комментарий к посту
func (s *CommentService) Create(ctx context.Context, postID, userID int, req *model.CommentCreateRequest) (*model.Comment, error) {
	if err := validateCommentCreateRequest(req); err != nil {
		return nil, err
	}

	// 1. Проверяем, что пост существует
	_, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		if errors.Is(err, apperr.ErrPostNotFound) {
			return nil, apperr.ErrPostNotExists // 404
		}
		return nil, fmt.Errorf("failed to check post existence: %w", err)
	}

	comment := &model.Comment{
		Content:  strings.TrimSpace(req.Content),
		PostID:   postID,
		AuthorID: userID,
	}

	if err := s.commentRepo.Create(ctx, comment); err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	return comment, nil
}

// GetByID получает комментарий по ID
func (s *CommentService) GetByID(ctx context.Context, id int) (*model.Comment, error) {
	comment, err := s.commentRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, apperr.ErrCommentNotFound) {
			return nil, apperr.ErrCommentNotFound
		}
		return nil, fmt.Errorf("failed to get comment by ID: %w", err)
	}
	return comment, nil
}

// GetByPost получает комментарии к посту с пагинацией
func (s *CommentService) GetByPost(ctx context.Context, postID int, limit, offset int) ([]*model.Comment, int, error) {
	const (
		defaultLimit = 20
		maxLimit     = 100
	)

	if limit <= 0 {
		limit = defaultLimit
	} else if limit > maxLimit {
		limit = maxLimit
	}
	if offset < 0 {
		offset = 0
	}

	// проверка существования поста
	_, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		if errors.Is(err, apperr.ErrPostNotFound) {
			return nil, 0, apperr.ErrPostNotExists
		}
		return nil, 0, fmt.Errorf("failed to check post existence for comments: %w", err)
	}

	comments, err := s.commentRepo.GetByPostID(ctx, postID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list comments by post: %w", err)
	}

	total, err := s.commentRepo.GetCountByPostID(ctx, postID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count comments by post: %w", err)
	}

	return comments, total, nil
}

// Update обновляет комментарий (только content)
func (s *CommentService) Update(ctx context.Context, id int, userID int, req *model.CommentUpdateRequest) (*model.Comment, error) {
	comment, err := s.commentRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, apperr.ErrCommentNotFound) {
			return nil, apperr.ErrCommentNotFound
		}
		return nil, fmt.Errorf("failed to get comment for update: %w", err)
	}

	// Используем метод модели для проверки прав
	if !comment.CanBeEditedBy(userID) {
		return nil, apperr.ErrForbidden
	}

	if err := validateCommentUpdateRequest(req); err != nil {
		return nil, err
	}

	comment.Content = strings.TrimSpace(req.Content)
	// updated_at обновится внутри репозитория

	if err := s.commentRepo.Update(ctx, comment); err != nil {
		if errors.Is(err, apperr.ErrCommentNotFound) {
			return nil, apperr.ErrCommentNotFound
		}
		return nil, fmt.Errorf("failed to update comment: %w", err)
	}

	return comment, nil
}

// Delete удаляет комментарий
func (s *CommentService) Delete(ctx context.Context, id int, userID int) error {
	comment, err := s.commentRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, apperr.ErrCommentNotFound) {
			return apperr.ErrCommentNotFound
		}
		return fmt.Errorf("failed to get comment for delete: %w", err)
	}

	if !comment.CanBeDeletedBy(userID) {
		return apperr.ErrForbidden
	}

	if err := s.commentRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, apperr.ErrCommentNotFound) {
			return apperr.ErrCommentNotFound
		}
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	return nil
}

// GetByAuthor получает комментарии конкретного автора с пагинацией

func (s *CommentService) GetByAuthor(ctx context.Context, authorID int, limit, offset int) ([]*model.Comment, int, error) {
	const (
		defaultLimit = 20
		maxLimit     = 100
	)

	if limit <= 0 {
		limit = defaultLimit
	} else if limit > maxLimit {
		limit = maxLimit
	}
	if offset < 0 {
		offset = 0
	}

	comments, err := s.commentRepo.ListByAuthor(ctx, authorID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list comments by author: %w", err)
	}

	total, err := s.commentRepo.CountByAuthor(ctx, authorID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count comments by author: %w", err)
	}

	return comments, total, nil
}

func validateCommentCreateRequest(req *model.CommentCreateRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	content := strings.TrimSpace(req.Content)
	if len(content) == 0 {
		return errors.New("content is required")
	}
	if len(content) > 1000 {
		return errors.New("content must be no more than 1000 characters")
	}

	return nil
}

func validateCommentUpdateRequest(req *model.CommentUpdateRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	content := strings.TrimSpace(req.Content)
	if len(content) == 0 {
		return errors.New("content cannot be empty")
	}
	if len(content) > 1000 {
		return errors.New("content must be no more than 1000 characters")
	}

	return nil
}
