package storage

import (
	"sync"
	"wildberries_test_task/internal/models"
)

type MemStorage struct {
	sync.RWMutex
	storage map[string]models.Msg
}

func InitializeMemoryStorage() *MemStorage {
	return &MemStorage{storage: make(map[string]models.Msg)}
}

type Storage interface {
	Set(msg models.Msg)
	Get(key string) (*models.UserGrade, bool)
}

func (m *MemStorage) Set(new models.Msg) {
	m.Lock()
	defer m.Unlock()
	key := new.UserGrade.UserId
	old, ok := m.storage[key]
	if ok {
		if new.Timestamp > old.Timestamp || new.Timestamp == old.Timestamp && new.Priority > old.Priority {
			m.storage[key] = new
		}
	} else {
		m.storage[key] = new
	}
}

func (m *MemStorage) Get(key string) (*models.UserGrade, bool) {
	m.RLock()
	defer m.RUnlock()
	msg, ok := m.storage[key]
	if ok {
		return &msg.UserGrade, true
	}
	return nil, false
}
