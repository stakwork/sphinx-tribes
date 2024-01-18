package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stakwork/sphinx-tribes/auth"
)

func TestDeletePerson(t *testing.T) {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/person/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.WithValue(req.Context(), auth.ContextKey, "pubKey")
	req = req.WithContext(ctx)
	DeletePerson(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expected := `true`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

