package handlers

import (
	"github.com/IBM/sarama"
	"github.com/gorilla/schema"
	"github.com/jackc/pgx/v5"
	"net/http"
)

type Handler struct {
	Producer sarama.SyncProducer
	Conn     *pgx.Conn
}

var decoder = schema.NewDecoder()

type SendMessageRequest struct {
	fromID string
}

func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
}
