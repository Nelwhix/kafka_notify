package pkg

import (
	"log"
	"log/slog"
	"os"
	"path/filepath"
)

func CreateNewLogger() (*slog.Logger, error) {
	fileName := filepath.Join("logs", "app_logs.txt")
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	return slog.New(slog.NewJSONHandler(f, nil)), nil
}
