package service

import (
	"context"
	"wildberries_test_task/internal/models"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s Service) GetUserGrade(ctx context.Context) (models.UserGrade, error) {

	return models.UserGrade{}, nil
}

func (s Service) SetUserGrade(ctx context.Context, u models.UserGrade) error {

	return nil
}
