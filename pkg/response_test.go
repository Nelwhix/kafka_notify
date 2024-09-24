package pkg

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewInternalServerErrorResponse(t *testing.T) {
	responseMsg := "Internal Server Error"
	recorder := httptest.NewRecorder()
	NewInternalServerErrorResponse(recorder, responseMsg)

	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, recorder.Code)
	}

	var actualResponse baseResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &actualResponse)
	if err != nil {
		t.Fatalf("Could not unmarshal actual response: %v", err)
	}

	expectedResponse := baseResponse{
		Message: responseMsg,
	}

	if actualResponse != expectedResponse {
		t.Errorf("Expected response %+v, got %+v", expectedResponse, actualResponse)
	}
}

func TestNewUnprocessableEntityResponse(t *testing.T) {
	responseMsg := "validation error"
	recorder := httptest.NewRecorder()
	NewUnprocessableEntityResponse(recorder, responseMsg)

	if recorder.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status code %d, got %d", http.StatusUnprocessableEntity, recorder.Code)
	}

	var actualResponse baseResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &actualResponse)
	if err != nil {
		t.Fatalf("Could not unmarshal actual response: %v", err)
	}

	expectedResponse := baseResponse{
		Message: responseMsg,
	}

	if actualResponse != expectedResponse {
		t.Errorf("Expected response %+v, got %+v", expectedResponse, actualResponse)
	}
}

func TestNewNotFoundResponse(t *testing.T) {
	responseMessage := "user not found"
	recorder := httptest.NewRecorder()
	NewNotFoundResponse(recorder, responseMessage)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, recorder.Code)
	}

	var actualResponse baseResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &actualResponse)
	if err != nil {
		t.Fatalf("Could not unmarshal actual response: %v", err)
	}

	expectedResponse := baseResponse{
		Message: responseMessage,
	}

	if actualResponse != expectedResponse {
		t.Errorf("Expected response %+v, got %+v", expectedResponse, actualResponse)
	}
}

func TestNewOkResponse(t *testing.T) {
	responseMessage := "Get tests."
	recorder := httptest.NewRecorder()
	NewOKResponse(recorder, responseMessage)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
	}

	var actualResponse baseResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &actualResponse)
	if err != nil {
		t.Fatalf("Could not unmarshal actual response: %v", err)
	}

	expectedResponse := baseResponse{
		Message: responseMessage,
	}

	if actualResponse != expectedResponse {
		t.Errorf("Expected response %+v, got %+v", expectedResponse, actualResponse)
	}
}
