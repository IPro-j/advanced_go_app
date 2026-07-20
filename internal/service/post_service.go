package service

import (
	"blog-api/internal/model"
	"blog-api/internal/repository"
	"context"
	"errors"
	"fmt"
)

var (
	ErrPostNotFound = errors.New("post not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
)

type PostService struct {
	postRepo repository.PostRepository
	userRepo repository.UserRepository
}

func NewPostService(postRepo repository.PostRepository, userRepo repository.UserRepository) *PostService {
	return &PostService{
		postRepo: postRepo,
		userRepo: userRepo,
	}
}

func (s *PostService) Create(ctx context.Context, userID int, req *model.PostCreateRequest) (*model.Post, error) {
	// TODO: Создать новый пост
	// Шаги:
	// 1. Валидация данных (title не пустой и <= 200 символов, content не пустой)
	// 2. Создать модель поста с данными из запроса и userID
	// 3. Сохранить через репозиторий
	// 4. Вернуть созданный пост

	return nil, fmt.Errorf("not implemented")
}

func (s *PostService) GetByID(ctx context.Context, id int) (*model.Post, error) {
	// TODO: Получить пост по ID
	// Шаги:
	// 1. Получить пост через репозиторий
	// 2. Опционально: загрузить информацию об авторе
	// 3. Вернуть пост

	return nil, fmt.Errorf("not implemented")
}

func (s *PostService) GetAll(ctx context.Context, limit, offset int) ([]*model.Post, int, error) {
	// TODO: Получить список постов с пагинацией
	// Шаги:
	// 1. Валидировать и нормализовать параметры пагинации (limit по умолчанию 10, максимум 100)
	// 2. Получить посты через репозиторий
	// 3. Получить общее количество для пагинации
	// 4. Опционально: обогатить данные информацией об авторах
	// 5. Вернуть посты и общее количество

	return nil, 0, fmt.Errorf("not implemented")
}

func (s *PostService) Update(ctx context.Context, id int, userID int, req *model.PostUpdateRequest) (*model.Post, error) {
	// TODO: Обновить пост
	// Шаги:
	// 1. Получить существующий пост
	// 2. Проверить что userID является автором (иначе ErrForbidden)
	// 3. Валидировать новые данные (если предоставлены)
	// 4. Обновить только измененные поля
	// 5. Сохранить через репозиторий
	// 6. Вернуть обновленный пост

	return nil, fmt.Errorf("not implemented")
}

func (s *PostService) Delete(ctx context.Context, id int, userID int) error {
	// TODO: Удалить пост
	// Шаги:
	// 1. Найти пост и проверить существование
	// 2. Проверить что userID является автором
	// 3. Удалить через репозиторий
	// 4. Вернуть соответствующую ошибку при неудаче

	return fmt.Errorf("not implemented")
}

func (s *PostService) GetByAuthor(ctx context.Context, authorID int, limit, offset int) ([]*model.Post, int, error) {
	// TODO: Получить посты конкретного автора
	// Шаги:
	// 1. Валидировать параметры пагинации
	// 2. Получить посты автора через репозиторий
	// 3. Получить общее количество постов автора
	// 4. Опционально: добавить информацию об авторе к постам
	// 5. Вернуть результат с общим количеством

	return nil, 0, fmt.Errorf("not implemented")
}

// validatePostCreateRequest проверяет корректность данных для создания поста
func validatePostCreateRequest(req *model.PostCreateRequest) error {
	// TODO: Реализовать валидацию title и content

	return nil
}

// validatePostUpdateRequest проверяет корректность данных для обновления поста
func validatePostUpdateRequest(req *model.PostUpdateRequest) error {
	// TODO: Реализовать валидацию опциональных полей

	return nil
}
