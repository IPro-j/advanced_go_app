package repository

import (
	"blog-api/internal/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var (
	ErrPostNotFound = errors.New("post not found")
)

type PostRepo struct {
	db *sql.DB
}

func NewPostRepo(db *sql.DB) *PostRepo {
	return &PostRepo{db: db}
}

// Create сохраняет новый пост.
func (r *PostRepo) Create(ctx context.Context, post *model.Post) error {
	const query = `
		INSERT INTO posts (title, content, author_id, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query, post.Title, post.Content, post.AuthorID).
		Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}

	return nil
}

func (r *PostRepo) Exists(ctx context.Context, id int) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM posts WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetByID получает пост по ID. Если не найден — возвращает ErrPostNotFound.
func (r *PostRepo) GetByID(ctx context.Context, id int) (*model.Post, error) {
	const query = `
		SELECT id, title, content, author_id, created_at, updated_at
		FROM posts
		WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var p model.Post
	err := row.Scan(&p.ID, &p.Title, &p.Content, &p.AuthorID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPostNotFound
		}
		return nil, fmt.Errorf("failed to get post by ID: %w", err)
	}

	return &p, nil
}

func (r *PostRepo) GetTotalCount(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM posts`
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *PostRepo) GetAll(ctx context.Context, limit, offset int) ([]*model.Post, error) {
	query := `
		SELECT id, title, content, author_id, created_at, updated_at
		FROM posts
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*model.Post
	for rows.Next() {
		var p model.Post
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.AuthorID,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &p)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return posts, nil
}

// Update обновляет существующие поля поста.
func (r *PostRepo) Update(ctx context.Context, post *model.Post) error {
	const query = `
		UPDATE posts
		SET title = $1, content = $2, updated_at = NOW()
		WHERE id = $3
	`

	res, err := r.db.ExecContext(ctx, query, post.Title, post.Content, post.ID)
	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if n == 0 {
		// Пост с таким ID не найден — можно вернуть NotFound, но обычно это проверяется до вызова Update
		return ErrPostNotFound
	}

	// Обновляем updated_at в памяти, чтобы не делать лишний SELECT
	post.UpdatedAt = time.Now()

	return nil
}

// Delete удаляет пост по ID.
func (r *PostRepo) Delete(ctx context.Context, id int) error {
	const query = `DELETE FROM posts WHERE id = $1`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if n == 0 {
		return ErrPostNotFound
	}

	return nil
}

// ListByAuthor возвращает посты конкретного автора с пагинацией.
func (r *PostRepo) ListByAuthor(ctx context.Context, authorID int, limit, offset int) ([]*model.Post, error) {
	const query = `
		SELECT id, title, content, author_id, created_at, updated_at
		FROM posts
		WHERE author_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, authorID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list posts by author: %w", err)
	}
	defer rows.Close()

	var posts []*model.Post
	for rows.Next() {
		var p model.Post
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.AuthorID, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post row for author: %w", err)
		}
		posts = append(posts, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating posts by author rows: %w", err)
	}

	return posts, nil
}

// CountByAuthor возвращает общее количество постов автора (для пагинации).
func (r *PostRepo) CountByAuthor(ctx context.Context, authorID int) (int, error) {
	const query = `SELECT COUNT(*) FROM posts WHERE author_id = $1`

	var total int
	err := r.db.QueryRowContext(ctx, query, authorID).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to count posts by author: %w", err)
	}

	return total, nil
}
