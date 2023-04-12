package routes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stakwork/sphinx-tribes/frontend"
)

func TestPingRoute(t *testing.T) {
	var result string = "pong"
	req, err := http.NewRequest("GET", "/ping", nil)

	if err != nil {
		t.Errorf("Error creating a new request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(frontend.PingRoute)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code. Expected: %d. Got: %d.", http.StatusOK, status)
	}

	var response string

	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	if response != result {
		t.Errorf("Expected: %v. Got: %v.", response, result)
	}
}

func TestIndexRoute(t *testing.T) {
	var result string = "Sphinx Community"
	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Errorf("Error creating a new request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(frontend.IndexRoute)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code. Expected: %d. Got: %d.", http.StatusOK, status)
	}

	if !strings.Contains(rr.Body.String(), result) {
		t.Error("Not the Index Page")
	}
}
