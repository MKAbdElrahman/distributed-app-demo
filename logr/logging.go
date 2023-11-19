package logr

import (
	"log/slog"
	"os"
	"path/filepath"
)

type Logger struct {
	*slog.Logger
}

type option func(*Logger) error

func NewLogWriter(options ...option) (*Logger, error) {
	logWriter := &Logger{
		Logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}

	for _, opt := range options {
		err := opt(logWriter)
		if err != nil {
			return nil, err
		}
	}

	return logWriter, nil
}

func DefaultFileLogger(path string) (*Logger, error) {
	w, err := NewLogWriter(WithFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600))

	if err != nil {
		return nil, err
	}
	return w, nil
}

// WithFile sets the log writer to a file.
func WithFile(path string, flag int, perm os.FileMode) option {
	return func(lw *Logger) error {
		// Ensure the directory exists before attempting to open the file
		dir := filepath.Dir(path)
		if err := ensureDir(dir); err != nil {
			return err
		}

		f, err := os.OpenFile(path, flag, perm)
		if err != nil {
			return err
		}

		lw.Logger = slog.New(slog.NewJSONHandler(f, nil))
		return nil
	}
}

// ensureDir function ensures that the directory exists; creates it if not.
func ensureDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}
