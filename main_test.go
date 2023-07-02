package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestHandleGet(t *testing.T) {
	av := NewAppVersions()
	av.SetVersion("dev", "app1", 1)

	req, err := http.NewRequest("GET", "/api/version?env=dev&app=app1", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(av.APIHandler)
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `1` + "\n"
	if rec.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rec.Body.String(), expected)
	}
}

func TestHandlePost(t *testing.T) {
	av := NewAppVersions()

	form := new(bytes.Buffer)
	fmt.Fprintf(form, "env=dev&app=app1&version=%s", strconv.Itoa(2))

	req, err := http.NewRequest("POST", "/api/version", form)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(av.APIHandler)
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	if version := av.GetVersions("dev")["app1"]; version != 2 {
		t.Errorf("Handler did not update version correctly: got %v want %v", version, 2)
	}
}
