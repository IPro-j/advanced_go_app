package service

import (
	"blog-api/internal/model"
	"blog-api/internal/repository"
	"context"
	"errors"
	"fmt"
)

var (
	ErrCommentNotFound = errors.New("comment not found")
	ErrPostNotExists   = errors.New("post does not exist")
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

func (s *CommentService) Create(ctx context.Context, userID int, req *model.CommentCreateRequest) (*model.Comment, error) {
	// TODO: Создать новый комментарий
	// Шаги:
	// 1. Валидация данных (content не пустой и <= 1000 символов)
	// 2. Проверить что пост существует
	// 3. Создать модель комментария с userID как автором
	// 4. Сохранить через репозиторий
	// 5. Опционально: обогатить ответ информацией об авторе
	// 6. Вернуть созданный комментарий

	return nil, fmt.Errorf("not implemented")
}

func (s *CommentService) GetByID(ctx context.Context, id int) (*model.Comment, error) {
	// TODO: Получить комментарий по ID
	// Шаги:
	// 1. Получить комментарий через репозиторий
	// 2. Опционально: добавить информацию об авторе
	// 3. Вернуть результат или ErrCommentNotFound

	return nil, fmt.Errorf("not implemented")
}

func (s *CommentService) GetByPost(ctx context.Context, postID int, limit, offset int) ([]*model.Comment, int, error) {
	// TODO: Получить комментарии к посту с пагинацией
	// Шаги:
	// 1. Валидировать параметры пагинации (limit по умолчанию 20, максимум 100)
	// 2. Опционально: проверить существование поста
	// 3. Получить комментарии через репозиторий
	// 4. Получить общее количество для пагинации
	// 5. Опционально: обогатить данные информацией об авторах
	// 6. Вернуть комментарии и общее количество

	return nil, 0, fmt.Errorf("not implemented")
}

func (s *CommentService) Update(ctx context.Context, id int, userID int, req *model.CommentUpdateRequest) (*model.Comment, error) {
	// TODO: Обновить комментарий
	// Шаги:
	// 1. Найти существующий комментарий
	// 2. Проверить что userID является автором (иначе ErrForbidden)
	// 3. Валидировать новый content
	// 4. Обновить content и временную метку
	// 5. Сохранить через репозиторий
	// 6. Опционально: добавить информацию об авторе
	// 7. Вернуть обновленный комментарий

	return nil, fmt.Errorf("not implemented")
}

func (s *CommentService) Delete(ctx context.Context, id int, userID int) error {
	// TODO: Удалить комментарий
	// Шаги:
	// 1. Найти комментарий и проверить существование
	// 2. Проверить что userID является автором
	// 3. Удалить через репозиторий
	// 4. Вернуть соответствующую ошибку при неудаче

	return fmt.Errorf("not implemented")
}

func (s *CommentService) GetByAuthor(ctx context.Context, authorID int, limit, offset int) ([]*model.Comment, int, error) {
	// TODO: Получить комментарии конкретного автора
	// Шаги:
	// 1. Валидировать параметры пагинации
	// 2. Получить комментарии автора через репозиторий
	// 3. Получить общее количество комментариев автора
	// 4. Опционально: добавить информацию об авторе
	// 5. Вернуть результат с общим количеством

	return nil, 0, fmt.Errorf("not implemented")
}

// validateCommentCreateRequest проверяет корректность данных для создания комментария
func validateCommentCreateRequest(req *model.CommentCreateRequest) error {
	// TODO: Реализовать валидацию content и PostID

	return nil
}

// validateCommentUpdateRequest проверяет корректность данных для обновления комментария
func validateCommentUpdateRequest(req *model.CommentUpdateRequest) error {
	// TODO: Реализовать валидацию content

	return nil
}
