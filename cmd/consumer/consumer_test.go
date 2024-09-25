package main

import (
	"github.com/Nelwhix/kafka-notify/cmd/consumer/handlers"
	"github.com/Nelwhix/kafka-notify/pkg/models"
	"net/http"
	"testing"
)

func Setup(*testing.T) *http.ServeMux {
	store := &models.NotificationStore{
		Data: make(models.UserNotifications),
	}
	handler := handlers.Handler{
		NotificationStore: store,
	}

	r := http.NewServeMux()
	r.HandleFunc("GET /notifications/{userId}", handler.GetNotifications)

	return r
}
