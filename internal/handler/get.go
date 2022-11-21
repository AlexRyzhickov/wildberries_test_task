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
	GetUserGrade(ctx context.Context, userId string) (*models.UserGrade, error)
}

func (h *GetUserGradeHandler) Method() string {
	return http.MethodGet
}

func (h *GetUserGradeHandler) Path() string {
	return "/get"
}

func (h *GetUserGradeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userId := r.FormValue("user_id")
	userGrade, err := h.Service.GetUserGrade(r.Context(), userId)
	if err != nil {
		writeResponse(w, r, err)
		return
	}
	writeResponse(w, r, &userGrade)
}
