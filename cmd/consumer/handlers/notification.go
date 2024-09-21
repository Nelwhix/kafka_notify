package handlers

import "net/http"

type Handler struct {
}

func (h *Handler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	idString := req.PathValue("id")
}
