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
	ErrCommentNotFound = errors.New("comment not found")
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
	// TODO: Реализовать создание комментария
	// 1. Подготовить SQL запрос INSERT INTO comments...
	// 2. Установить created_at и updated_at = time.Now()
	// 3. Выполнить запрос и получить ID созданной записи
	// 4. Установить ID в структуру comment
	//
	// HINT: Используйте QueryRowContext с RETURNING id
	// Пример запроса:
	// INSERT INTO comments (content, post_id, author_id, created_at, updated_at)
	// VALUES ($1, $2, $3, $4, $5) RETURNING id

	query := `
		INSERT INTO comments (content, post_id, author_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	now := time.Now()
	comment.CreatedAt = now
	comment.UpdatedAt = now

	// TODO: Выполнить запрос
	// err := r.db.QueryRowContext(ctx, query,
	//     comment.Content, comment.PostID, comment.AuthorID,
	//     comment.CreatedAt, comment.UpdatedAt,
	// ).Scan(&comment.ID)

	_ = query // Удалите эту строку после реализации
	return fmt.Errorf("not implemented")
}

// GetByID получает комментарий по ID
func (r *CommentRepo) GetByID(ctx context.Context, id int) (*model.Comment, error) {
	// TODO: Реализовать получение комментария по ID
	// 1. Подготовить SQL запрос SELECT ... FROM comments WHERE id = $1
	// 2. Выполнить запрос
	// 3. Просканировать результат в структуру Comment
	// 4. Обработать случай sql.ErrNoRows -> вернуть ErrCommentNotFound

	query := `
		SELECT id, content, post_id, author_id, created_at, updated_at
		FROM comments
		WHERE id = $1
	`

	var comment model.Comment
	// TODO: Выполнить запрос и просканировать результат
	// err := r.db.QueryRowContext(ctx, query, id).Scan(
	//     &comment.ID, &comment.Content, &comment.PostID,
	//     &comment.AuthorID, &comment.CreatedAt, &comment.UpdatedAt,
	// )
	// if err != nil {
	//     if err == sql.ErrNoRows {
	//         return nil, ErrCommentNotFound
	//     }
	//     return nil, fmt.Errorf("failed to get comment: %w", err)
	// }

	_ = query // Удалите эту строку после реализации
	return nil, fmt.Errorf("not implemented")
}

// GetByPostID получает комментарии к посту с пагинацией
func (r *CommentRepo) GetByPostID(ctx context.Context, postID int, limit, offset int) ([]*model.Comment, error) {
	// TODO: Реализовать получение комментариев к посту
	// 1. Подготовить SQL запрос с WHERE post_id = $1
	// 2. Добавить ORDER BY created_at ASC (комментарии по времени)
	// 3. Добавить LIMIT и OFFSET для пагинации
	// 4. Выполнить запрос и получить rows
	// 5. Итерировать по rows и собрать массив комментариев
	// 6. Не забудьте закрыть rows (defer rows.Close())
	//
	// HINT: Используйте QueryContext для получения множества записей
	// Пример запроса:
	// SELECT id, content, post_id, author_id, created_at, updated_at
	// FROM comments
	// WHERE post_id = $1
	// ORDER BY created_at ASC
	// LIMIT $2 OFFSET $3

	query := `
		SELECT id, content, post_id, author_id, created_at, updated_at
		FROM comments
		WHERE post_id = $1
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3
	`

	// TODO: Выполнить запрос
	// rows, err := r.db.QueryContext(ctx, query, postID, limit, offset)
	// if err != nil {
	//     return nil, fmt.Errorf("failed to get comments: %w", err)
	// }
	// defer rows.Close()

	// TODO: Итерировать по результатам
	// var comments []*model.Comment
	// for rows.Next() {
	//     var comment model.Comment
	//     err := rows.Scan(
	//         &comment.ID, &comment.Content, &comment.PostID,
	//         &comment.AuthorID, &comment.CreatedAt, &comment.UpdatedAt,
	//     )
	//     if err != nil {
	//         return nil, fmt.Errorf("failed to scan comment: %w", err)
	//     }
	//     comments = append(comments, &comment)
	// }
	//
	// if err = rows.Err(); err != nil {
	//     return nil, fmt.Errorf("failed to iterate comments: %w", err)
	// }
	//
	// return comments, nil

	_ = query // Удалите эту строку после реализации
	return nil, fmt.Errorf("not implemented")
}

// GetCountByPostID получает количество комментариев к посту
func (r *CommentRepo) GetCountByPostID(ctx context.Context, postID int) (int, error) {
	// TODO: Реализовать подсчет комментариев к посту
	// HINT: SELECT COUNT(*) FROM comments WHERE post_id = $1

	query := `SELECT COUNT(*) FROM comments WHERE post_id = $1`

	var count int
	// TODO: Выполнить запрос
	// err := r.db.QueryRowContext(ctx, query, postID).Scan(&count)

	_ = query // Удалите эту строку после реализации
	return 0, fmt.Errorf("not implemented")
}

// Update обновляет комментарий
func (r *CommentRepo) Update(ctx context.Context, comment *model.Comment) error {
	// TODO: (Опционально) Реализовать обновление комментария
	// 1. Обновить только content и updated_at
	// 2. Использовать UPDATE comments SET content = $1, updated_at = $2 WHERE id = $3
	// 3. Проверить RowsAffected

	query := `
		UPDATE comments
		SET content = $1, updated_at = $2
		WHERE id = $3
	`

	comment.UpdatedAt = time.Now()

	// TODO: Выполнить запрос

	_ = query // Удалите эту строку после реализации
	return fmt.Errorf("not implemented")
}

// Delete удаляет комментарий
func (r *CommentRepo) Delete(ctx context.Context, id int) error {
	// TODO: (Опционально) Реализовать удаление комментария
	// 1. DELETE FROM comments WHERE id = $1
	// 2. Проверить RowsAffected

	query := `DELETE FROM comments WHERE id = $1`

	// TODO: Выполнить запрос

	_ = query // Удалите эту строку после реализации
	return fmt.Errorf("not implemented")
}
