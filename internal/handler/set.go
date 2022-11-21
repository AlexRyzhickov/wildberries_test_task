package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"wildberries_test_task/internal/models"
)

type SetUserGradeHandler struct {
	Service SetUserGradeService
}

type SetUserGradeService interface {
	SetUserGrade(ctx context.Context, u models.UserGrade)
}

func (h *SetUserGradeHandler) Method() string {
	return http.MethodPost
}

func (h *SetUserGradeHandler) Path() string {
	return "/set"
}

func (h *SetUserGradeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userGrade := models.UserGrade{}
	err := json.NewDecoder(r.Body).Decode(&userGrade)
	if err != nil {
		writeResponse(w, r, err)
		return
	}
	h.Service.SetUserGrade(r.Context(), userGrade)
	writeResponse(w, r, "User Grade stored successfully")
}
