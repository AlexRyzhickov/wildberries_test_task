package handler

import (
	"context"
	"net/http"
)

type BackupHandler struct {
	Service BackupService
}

type BackupService interface {
	Backup(ctx context.Context) ([]byte, error)
}

func (h *BackupHandler) Method() string {
	return http.MethodGet
}

func (h *BackupHandler) Path() string {
	return "/backup"
}

func (h *BackupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bytes, err := h.Service.Backup(r.Context())
	if err != nil {
		writeResponse(w, r, err)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename=backup.gz")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
