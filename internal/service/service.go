package service

import (
	"context"
	"errors"
	"wildberries_test_task/internal/models"
	"wildberries_test_task/internal/storage"
)

type Service struct {
	storage storage.Storage
}

func NewService(storage storage.Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s Service) GetUserGrade(ctx context.Context, userId string) (*models.UserGrade, error) {
	userGrade, ok := s.storage.Get(userId)
	if !ok {
		return &models.UserGrade{}, errors.New("storage error")
	}
	return userGrade, nil
}

func (s Service) SetUserGrade(ctx context.Context, u models.UserGrade) {
	s.storage.Set(u.UserId, u)
}
