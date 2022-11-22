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

func (s Service) SetUserGrade(ctx context.Context, grade models.UserGrade) {
	gradeStored, ok := s.storage.Get(grade.UserId)
	if !ok {
		s.storage.Set(grade.UserId, grade)
		return
	}
	if grade.Spp == 0 {
		grade.Spp = gradeStored.Spp
	}
	if grade.PostpaidLimit == 0 {
		grade.PostpaidLimit = gradeStored.PostpaidLimit
	}
	if grade.Spp == 0 {
		grade.Spp = gradeStored.Spp
	}
	if grade.ShippingFee == 0 {
		grade.ShippingFee = gradeStored.ShippingFee
	}
	if grade.ReturnFee == 0 {
		grade.ReturnFee = gradeStored.ReturnFee
	}
	s.storage.Set(grade.UserId, grade)
}
