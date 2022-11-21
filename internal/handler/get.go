package handler

import (
	"context"
	"net/http"
	"wildberries_test_task/internal/models"
)

type GetUserGradeHandler struct {
	Service GetUserGradeService
}

type GetUserGradeService interface {
	GetUserGrade(ctx context.Context) (models.UserGrade, error)
}

func (h *GetUserGradeHandler) Method() string {
	return http.MethodGet
}

func (h *GetUserGradeHandler) Path() string {
	return "/get"
}

func (h *GetUserGradeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	writeResponse(w, r, "get")
}
