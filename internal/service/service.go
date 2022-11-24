package service

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"wildberries_test_task/internal/models"
	"wildberries_test_task/internal/storage"
)

const topic = "topic"
const backupFileName = "backup.csv"
const csvColumns = "UserId,PostpaidLimit,Spp,ShippingFee,ReturnFee\n"
const storageErr = "storage error"

type Service struct {
	storage      storage.Storage
	nc           *nats.Conn
	priority     uint
	replicasPort string
}

func NewService(storage storage.Storage, nc *nats.Conn, priority uint, replicasAvailability bool, replicasPort string) *Service {
	s := Service{
		storage:      storage,
		nc:           nc,
		priority:     priority,
		replicasPort: replicasPort,
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
	if replicasAvailability {
		s.restoreStorage(replicasPort)
	}
	return &s
}

func (s Service) GetUserGrade(ctx context.Context, userId string) (*models.UserGrade, error) {
	userGrade, ok := s.storage.Get(userId)
	if !ok {
		return &models.UserGrade{}, errors.New(storageErr)
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
	b, err := json.Marshal(&msg)
	if err != nil {
		return err
	}
	err = nc.Publish(topic, b)
	if err != nil {
		return err
	}
	return nil
}

func (s Service) Backup(ctx context.Context) ([]byte, error) {
	var buf bytes.Buffer
	w, err := gzip.NewWriterLevel(io.Writer(&buf), gzip.BestCompression)
	if err != nil {
		return nil, err
	}
	grades, lastModTime := s.storage.GetAll()
	w.Header.Name = backupFileName
	w.Header.ModTime = time.Unix(0, lastModTime)
	_, err = w.Write([]byte(csvColumns))
	if err != nil {
		return nil, err
	}
	for _, grade := range grades {
		s := fmt.Sprintf("%v,%v,%v,%v,%v\n",
			grade.UserId,
			grade.PostpaidLimit,
			grade.Spp,
			grade.ShippingFee,
			grade.ReturnFee,
		)
		_, err = w.Write([]byte(s))
		if err != nil {
			return nil, err
		}
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s Service) restoreStorage(replicasPort string) {
	httpClient := &http.Client{
		Timeout: time.Minute * 5,
	}
	url := fmt.Sprintf("%v:%v/backup", "http://localhost", replicasPort)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	res, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	reader, err := gzip.NewReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	b, err := io.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(b), "\n")
	for i, line := range lines {
		if i > 0 && i < len(lines)-1 {
			msg := models.Msg{
				Priority:  s.priority,
				Timestamp: reader.ModTime.UnixNano(),
				UserGrade: models.UserGrade{},
			}
			_, err := fmt.Sscanf(strings.Replace(line, ",", " ", -1), "%s %d %d %d %d",
				&msg.UserGrade.UserId,
				&msg.UserGrade.PostpaidLimit,
				&msg.UserGrade.Spp,
				&msg.UserGrade.ShippingFee,
				&msg.UserGrade.ReturnFee,
			)
			if err != nil {
				log.Fatal(err)
			}
			s.storage.Set(msg)
		}
	}
}
