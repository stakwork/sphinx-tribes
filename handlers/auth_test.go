package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestGetAdminPubkeys(t *testing.T) {

	os.Setenv("ADMIN_PUBKEYS", "test")

	req, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAdminPubkeys)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"pubkeys":["test"]}`
	if strings.TrimRight(rr.Body.String(), "\n") != expected {

		t.Errorf("handler returned unexpected body: expected %s pubkeys %s is there a space after?", expected, rr.Body.String())
	}
}
