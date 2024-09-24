package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/Nelwhix/kafka-notify/pkg"
	"github.com/Nelwhix/kafka-notify/pkg/models"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
)

type Handler struct {
	Producer  sarama.SyncProducer
	Model     models.BaseModel
	Logger    *slog.Logger
	Validator *validator.Validate
}

type SendMessageRequest struct {
	FromID  string `json:"fromID" validate:"required,ulid"`
	ToID    string `json:"toID" validate:"required,ulid"`
	Message string `json:"message" validate:"required"`
}

func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		pkg.NewUnprocessableEntityResponse(w, err.Error())
		return
	}
	defer r.Body.Close()

	var request SendMessageRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		pkg.NewUnprocessableEntityResponse(w, err.Error())
		return
	}

	err = h.Validator.Struct(request)
	if err != nil {
		pkg.NewUnprocessableEntityResponse(w, err.Error())
		return
	}

	err = h.sendKafkaMessage(r.Context(), request)
	if err != nil {
		if errors.Is(err, pkg.ErrUserNotFoundInProducer) {
			pkg.NewNotFoundResponse(w, err.Error())

			return
		}

		h.Logger.Error(err.Error())
		pkg.NewInternalServerErrorResponse(w, err.Error())

		return
	}

	pkg.NewOKResponse(w, "Notification sent successfully!")
}

func (h *Handler) sendKafkaMessage(ctx context.Context, request SendMessageRequest) error {
	fromUser, err := h.Model.FindUserByID(ctx, request.FromID)
	if err != nil {
		return pkg.ErrUserNotFoundInProducer
	}

	toUser, err := h.Model.FindUserByID(ctx, request.ToID)
	if err != nil {
		return pkg.ErrUserNotFoundInProducer
	}

	notification := models.Notification{
		From:    fromUser,
		To:      toUser,
		Message: request.Message,
	}
	notificationJSON, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: "notifications",
		Key:   sarama.StringEncoder(toUser.ID),
		Value: sarama.StringEncoder(notificationJSON),
	}

	_, _, err = h.Producer.SendMessage(msg)

	return err
}
