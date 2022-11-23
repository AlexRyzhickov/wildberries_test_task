package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/nats-io/nats.go"
	"log"
	"time"
	"wildberries_test_task/internal/models"
	"wildberries_test_task/internal/storage"
)

const topic = "topic"

type Service struct {
	storage  storage.Storage
	nc       *nats.Conn
	priority uint
}

func NewService(storage storage.Storage, nc *nats.Conn, priority uint) *Service {
	s := Service{
		storage:  storage,
		nc:       nc,
		priority: priority,
	}
	_, err := nc.Subscribe(topic, func(m *nats.Msg) {
		msg := models.Msg{}
		err := json.Unmarshal(m.Data, &msg)
		if err != nil {
			log.Println(err)
			return
		}
		s.storage.Set(msg)
	})
	if err != nil {
		log.Fatal(err)
	}
	return &s
}

func (s Service) GetUserGrade(ctx context.Context, userId string) (*models.UserGrade, error) {
	userGrade, ok := s.storage.Get(userId)
	if !ok {
		return &models.UserGrade{}, errors.New("storage error")
	}
	return userGrade, nil
}

func (s Service) SetUserGrade(ctx context.Context, grade models.UserGrade) error {
	gradeStored, ok := s.storage.Get(grade.UserId)
	if ok {
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
	}
	msg := models.Msg{
		Priority:  s.priority,
		Timestamp: time.Now().UnixNano(),
		UserGrade: grade,
	}
	err := publishUserGrade(s.nc, &msg)
	if err != nil {
		return err
	}
	s.storage.Set(msg)
	return nil
}

func publishUserGrade(nc *nats.Conn, msg *models.Msg) error {
	bytes, err := json.Marshal(&msg)
	if err != nil {
		return err
	}
	err = nc.Publish(topic, bytes)
	if err != nil {
		return err
	}
	return nil
}
