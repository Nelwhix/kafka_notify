package pkg

import (
	"encoding/json"
	"log"
	"net/http"
)

type baseResponse struct {
	Message string `json:"message"`
}

func NewBadRequestResponse(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	response := baseResponse{
		Message: message,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}

	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Fatal(err)
	}

	return
}

func NewInternalServerErrorResponse(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusInternalServerError)
	response := baseResponse{
		Message: message,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}

	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Fatal(err)
	}

	return
}

func NewUnprocessableEntityResponse(w http.ResponseWriter, message string) {
	response := baseResponse{
		Message: message,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}

	http.Error(w, string(jsonResponse), http.StatusUnprocessableEntity)
}

func NewNotFoundResponse(w http.ResponseWriter, message string) {
	response := baseResponse{
		Message: message,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write(jsonResponse)
}

func NewOKResponse(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusOK)
	response := baseResponse{
		Message: message,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}

	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Fatal(err)
	}

	return
}
