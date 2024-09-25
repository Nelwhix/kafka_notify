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
	"github.com/oklog/ulid/v2"
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

func (t *TestModel) FindUserByID(ctx context.Context, userID string) (models.User, error) {
	for _, user := range users {
		if user.ID == userID {
			return user, nil
		}
	}

	return models.User{}, pkg.ErrUserNotFoundInProducer
}

func Setup(t *testing.T) (*http.ServeMux, *mocks.SyncProducer) {
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

	return r, producer
}

func TestRequestMustHaveValidPayload(t *testing.T) {
	r, _ := Setup(t)
	req, _ := http.NewRequest(http.MethodPost, "/send", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status code %d, got %d", http.StatusUnprocessableEntity, w.Code)
	}
}

func TestFromIDAndToIDMustBeUlid(t *testing.T) {
	r, _ := Setup(t)
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

func TestUserNotFound(t *testing.T) {
	r, _ := Setup(t)

	request := handlers.SendMessageRequest{
		FromID:  ulid.Make().String(),
		ToID:    ulid.Make().String(),
		Message: faker.Word(),
	}

	jsonData, _ := json.Marshal(request)

	req, _ := http.NewRequest(http.MethodPost, "/send", bytes.NewBuffer(jsonData))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusUnprocessableEntity, w.Code)
	}
}

func TestSendMessage(t *testing.T) {
	r, producer := Setup(t)
	fromID := ulid.Make().String()
	toID := ulid.Make().String()

	users = append(users, models.User{
		ID:   fromID,
		Name: faker.Word(),
	})
	users = append(users, models.User{
		ID:   toID,
		Name: faker.Word(),
	})

	request := handlers.SendMessageRequest{
		FromID:  fromID,
		ToID:    toID,
		Message: faker.Word(),
	}

	// mocking that a certain message will be received
	producer.ExpectSendMessageWithCheckerFunctionAndSucceed(func(val []byte) error {
		return nil // Simulate a successful send
	})

	jsonData, _ := json.Marshal(request)

	req, _ := http.NewRequest(http.MethodPost, "/send", bytes.NewBuffer(jsonData))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}
