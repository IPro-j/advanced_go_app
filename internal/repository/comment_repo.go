package repository

import (
	"blog-api/internal/model"
	"blog-api/pkg/apperr"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// CommentRepo представляет репозиторий для работы с комментариями
type CommentRepo struct {
	db *sql.DB
}

// NewCommentRepo создает новый репозиторий комментариев
func NewCommentRepo(db *sql.DB) *CommentRepo {
	return &CommentRepo{db: db}
}

// Create создает новый комментарий
func (r *CommentRepo) Create(ctx context.Context, comment *model.Comment) error {

	query := `
		INSERT INTO comments (content, post_id, author_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	var id int
	var createdAt, updatedAt time.Time
	err := r.db.QueryRowContext(ctx, query, comment.Content, comment.PostID, comment.AuthorID).
		Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}

	comment.ID = id
	comment.CreatedAt = createdAt
	comment.UpdatedAt = updatedAt
	return nil
}

// GetByID получает комментарий по ID
func (r *CommentRepo) GetByID(ctx context.Context, id int) (*model.Comment, error) {
	var c model.Comment
	query := `SELECT id, content, post_id, author_id, created_at, updated_at FROM comments WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID, &c.Content, &c.PostID, &c.AuthorID, &c.CreatedAt, &c.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperr.ErrCommentNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get comment by id: %w", err)
	}
	return &c, nil
}

// GetByPostID получает комментарии к посту с пагинацией
func (r *CommentRepo) GetByPostID(ctx context.Context, postID int, limit, offset int) ([]*model.Comment, error) {

	// Получаем комментарии
	query := `
		SELECT id, content, post_id, author_id, created_at, updated_at
		FROM comments
		WHERE post_id = $1
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, postID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query comments: %w", err)
	}
	defer rows.Close()

	var comments []*model.Comment
	for rows.Next() {
		var c model.Comment
		err := rows.Scan(&c.ID, &c.Content, &c.PostID, &c.AuthorID, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, &c)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate comments: %w", err)
	}

	return comments, nil
}

// GetCountByPostID получает количество комментариев к посту
func (r *CommentRepo) GetCountByPostID(ctx context.Context, postID int) (int, error) {
	query := `SELECT COUNT(*) FROM comments WHERE post_id = $1`

	var count int
	err := r.db.QueryRowContext(ctx, query, postID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count comments: %w", err)
	}

	return count, nil
}

// Update обновляет комментарий
func (r *CommentRepo) Update(ctx context.Context, comment *model.Comment) error {
	comment.UpdatedAt = time.Now()

	query := `
		UPDATE comments
		SET content = $1, updated_at = $2
		WHERE id = $3
	`
	res, err := r.db.ExecContext(ctx, query, comment.Content, comment.UpdatedAt, comment.ID)
	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if n == 0 {
		return apperr.ErrCommentNotFound
	}
	return nil
}

// Delete удаляет комментарий
func (r *CommentRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM comments WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if n == 0 {
		return apperr.ErrCommentNotFound
	}
	return nil
}

// ListByAuthor возвращает комментарии автора с пагинацией
func (r *CommentRepo) ListByAuthor(ctx context.Context, authorID int, limit, offset int) ([]*model.Comment, error) {
	query := `
		SELECT id, content, post_id, author_id, created_at, updated_at
		FROM comments
		WHERE author_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, authorID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query comments by author: %w", err)
	}
	defer rows.Close()

	var comments []*model.Comment
	for rows.Next() {
		var c model.Comment
		err := rows.Scan(&c.ID, &c.Content, &c.PostID, &c.AuthorID, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, &c)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate comments by author: %w", err)
	}

	return comments, nil
}

// CountByAuthor возвращает количество комментариев автора
func (r *CommentRepo) CountByAuthor(ctx context.Context, authorID int) (int, error) {
	query := `SELECT COUNT(*) FROM comments WHERE author_id = $1`
	var count int
	err := r.db.QueryRowContext(ctx, query, authorID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count comments by author: %w", err)
	}
	return count, nil
}
