package handlers

import (
	"github.com/Nelwhix/kafka-notify/pkg"
	"github.com/Nelwhix/kafka-notify/pkg/models"
	"net/http"
)

type Handler struct {
	NotificationStore *models.NotificationStore
}

func (h *Handler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("userId")
	notifications := h.NotificationStore.Get(userId)
	if len(notifications) == 0 {
		pkg.NewOKResponseWithData(w, "No notifications found for user", make([]interface{}, 0))
		return
	}

	pkg.NewOKResponseWithData(w, "Get Notifications.", notifications)
}
