package main

import (
	"demo/logr"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type LogHandler struct {
	Logger *logr.Logger
}

func (rh *LogHandler) RegisterRoutes(r *chi.Mux) {
	r.Post("/log", rh.HandleLog)
}

func DefaultLogToFileHandler(path string) (*LogHandler, error) {
	l, err := logr.DefaultFileLogger(path)
	if err != nil {
		return nil, err
	}
	return &LogHandler{
		Logger: l,
	}, nil
}

func (h *LogHandler) HandleLog(w http.ResponseWriter, r *http.Request) {
	msg, err := io.ReadAll(r.Body)

	if err != nil || len(msg) == 0 {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.Logger.Info(string(msg))

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Log received successfully"}`))
}
