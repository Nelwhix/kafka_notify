package pkg

import "net/http"

func NewUnprocessableEntityResponse(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	_, err := w.Write([]byte(message))
	if err != nil {
		return
	}
}
