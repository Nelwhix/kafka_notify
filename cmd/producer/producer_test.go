package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/IBM/sarama/mocks"
	"github.com/Nelwhix/kafka-notify/cmd/producer/handlers"
	"github.com/Nelwhix/kafka-notify/pkg"
	"github.com/Nelwhix/kafka-notify/pkg/models"
	"github.com/go-faker/faker/v4"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var users []models.User
var store models.UserNotifications

type TestModel struct {
}

type TestProducer struct {
	store models.UserNotifications
}

func (t *TestModel) FindUserByID(ctx context.Context, userID string) (models.User, error) {
	for _, user := range users {
		if user.ID == userID {
			return user, nil
		}
	}

	return models.User{}, pkg.ErrUserNotFoundInProducer
}

func Setup(t *testing.T) *http.ServeMux {
	logger, err := pkg.CreateNewLogger(os.Stdout)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}

	model := &TestModel{}
	producer := mocks.NewSyncProducer(t, nil)

	validate = validator.New(validator.WithRequiredStructEnabled())

	handler := handlers.Handler{
		Producer:  producer,
		Model:     model,
		Logger:    logger,
		Validator: validate,
	}

	r := http.NewServeMux()
	r.HandleFunc("POST /send", handler.SendMessage)

	return r
}

func TestRequestMustHaveValidPayload(t *testing.T) {
	r := Setup(t)
	req, _ := http.NewRequest(http.MethodPost, "/send", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status code %d, got %d", http.StatusUnprocessableEntity, w.Code)
	}
}

func TestFromIDAndToIDMustBeUlid(t *testing.T) {
	r := Setup(t)
	request := handlers.SendMessageRequest{
		FromID:  "1",
		ToID:    "5",
		Message: faker.Word(),
	}

	jsonData, _ := json.Marshal(request)

	req, _ := http.NewRequest(http.MethodPost, "/send", bytes.NewBuffer(jsonData))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status code %d, got %d", http.StatusUnprocessableEntity, w.Code)
	}
}

func TestSendMessage(t *testing.T) {
	r := Setup(t)
	request := handlers.SendMessageRequest{
		FromID:  "",
		ToID:    "5",
		Message: faker.Word(),
	}

	jsonData, _ := json.Marshal(request)

	req, _ := http.NewRequest(http.MethodPost, "/send", bytes.NewBuffer(jsonData))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status code %d, got %d", http.StatusUnprocessableEntity, w.Code)
	}
}
