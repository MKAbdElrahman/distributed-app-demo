package handlers

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type LogHandler struct {
	Logger HTTPLogger
}

func (bh *LogHandler) RegisterRoutes(r *chi.Mux) {
	r.Post("/log", bh.HandleLog)
}

func (bh *LogHandler) HandleLog(w http.ResponseWriter, r *http.Request) {
	msg, err := io.ReadAll(r.Body)

	if err != nil || len(msg) == 0 {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := bh.Logger.Log(string(msg)); err != nil {
		http.Error(w, "Error handling log", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Log received successfully"}`))
}
