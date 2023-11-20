package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type LogHandler interface {
	Log(message string) error
}

type HTTPLogHandler struct {
	LoggerURL string
}

func (lh *HTTPLogHandler) Log(message string) error {
	resp, err := http.Post(lh.LoggerURL+"/log", "application/text", strings.NewReader(message))
	if err != nil {
		return fmt.Errorf("error sending log to logging service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("logging service responded with status code: %d", resp.StatusCode)
	}

	return nil
}

type BusinessHandler struct {
	LogHandler LogHandler
}

func (bh *BusinessHandler) RegisterRoutes(r *chi.Mux) {
	r.Post("/log", bh.HandleLog)
}

func (bh *BusinessHandler) HandleLog(w http.ResponseWriter, r *http.Request) {
	msg, err := io.ReadAll(r.Body)

	if err != nil || len(msg) == 0 {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := bh.LogHandler.Log(string(msg)); err != nil {
		http.Error(w, "Error handling log", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Log received successfully"}`))
}
