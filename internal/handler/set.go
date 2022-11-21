package handler

import (
	"context"
	"net/http"
	"wildberries_test_task/internal/models"
)

type SetUserGradeHandler struct {
	Service SetUserGradeService
}

type SetUserGradeService interface {
	SetUserGrade(ctx context.Context, u models.UserGrade) error
}

func (h *SetUserGradeHandler) Method() string {
	return http.MethodPost
}

func (h *SetUserGradeHandler) Path() string {
	return "/set"
}

func (h *SetUserGradeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	writeResponse(w, r, "post")
}
