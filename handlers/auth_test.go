package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestGetAdminPubkeys(t *testing.T) {

	os.Setenv("ADMIN_PUBKEYS", "test")

	req, err := http.NewRequest("GET", "/admin_pubkeys", nil)
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

func TestCreateConnectionCode(t *testing.T) {
	// Want to test successful requesst here
	var jsonStr = []byte(`{"id":0,"connection_string":"string","is_used":false,"date_created":"2015-09-15T11:50:00Z"}`)
	req, err := http.NewRequest("POST", "/connectioncodes", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateConnectionCode)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := ``
	if strings.TrimRight(rr.Body.String(), "\n") != expected {

		t.Errorf("handler returned unexpected body: expected %s pubkeys %s is there a space after?", expected, rr.Body.String())
	}

	// Want to get 406 malformed error code here
	var jsonStr2 = []byte(`{"id":0,"connection_string":"string","is_used":false,"date_created":"5T11:50:00Z"}`)
	req2, err2 := http.NewRequest("POST", "/connectioncodes", bytes.NewBuffer(jsonStr2))
	if err2 != nil {
		t.Fatal(err2)
	}
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req2)
	if status := rr2.Code; status != http.StatusNotAcceptable {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Want to get xxx malformed error code here, using id == 1 as in mock this will simulate failed db response
	var jsonStr3 = []byte(`{"id":1,"connection_string":"string","is_used":false,"date_created":"2016-09-15T11:50:00Z"}`)
	req3, err3 := http.NewRequest("POST", "/connectioncodes", bytes.NewBuffer(jsonStr3))
	if err3 != nil {
		t.Fatal(err3)
	}
	rr3 := httptest.NewRecorder()
	handler.ServeHTTP(rr3, req3)
	if status := rr3.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestGetConnectionCode(t *testing.T) {
	req, err := http.NewRequest("GET", "/connectioncodes", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetConnectionCode)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"connection_string":"test","date_created":null}`
	if strings.TrimRight(rr.Body.String(), "\n") != expected {

		t.Errorf("handler returned unexpected body: expected %s pubkeys %s is there a space after?", expected, rr.Body.String())
	}

}
