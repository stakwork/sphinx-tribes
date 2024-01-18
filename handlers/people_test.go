package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)



func TestGetPersonByIdWithInvalidAuthToken(t *testing.T) {
	req, err := http.NewRequest("GET", "/person/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "invalid_token")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetPersonById)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}

func TestDeletePersonSuccess(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/person/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "valid_token")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeletePerson)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `true`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestDeletePersonFailure(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/person/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "invalid_token")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeletePerson)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}
