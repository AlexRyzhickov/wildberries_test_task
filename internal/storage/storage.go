package storage

import (
	"sync"
	"wildberries_test_task/internal/models"
)

type MemStorage struct {
	storage sync.Map
}

type Storage interface {
	Set(key string, userGrade models.UserGrade)
	Get(key string) (*models.UserGrade, bool)
}

func (m *MemStorage) Set(key string, userGrade models.UserGrade) {
	m.storage.Store(key, userGrade)
}

func (m *MemStorage) Get(key string) (*models.UserGrade, bool) {
	value, ok := m.storage.Load(key)
	if !ok {
		return nil, ok
	}
	userGrade, ok := value.(models.UserGrade)
	if ok {
		return &userGrade, ok
	}
	return nil, ok
}
