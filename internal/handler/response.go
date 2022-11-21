package handler

import (
	"encoding/json"
	"log"
	"mime"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

type badRequest struct {
	Error string `json:"error"`
}

func writeResponse(w http.ResponseWriter, r *http.Request, v interface{}) {
	code := http.StatusOK

	setISE := func(err error) {
		reqID := middleware.GetReqID(r.Context())
		log.Printf("[%v] Internal server error: %v\n", reqID, err)
		code = http.StatusInternalServerError
		v = badRequest{"Internal server error"}
	}

	if err, ok := v.(error); ok {
		setISE(err)
	}
	if _, ok := v.(badRequest); ok {
		code = http.StatusBadRequest
	}

	body, err := json.Marshal(v)
	if err != nil {
		setISE(err)
	}
	w.Header().Add("Content-Type", mime.TypeByExtension(".json"))
	w.WriteHeader(code)
	w.Write(body)
}
