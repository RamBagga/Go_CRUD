package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetPatients(t *testing.T) {
	// Create a mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Mock the query
	rows := sqlmock.NewRows([]string{"id", "name", "email", "status"}).
		AddRow(1, "John Doe", "johndoe@example.com", "Healthy").
		AddRow(2, "Jane Doe", "janedoe@example.com", "Sick")
	mock.ExpectQuery("^SELECT (.+) FROM patients WHERE name != '' AND email != '' AND status != ''$").WillReturnRows(rows)

	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/patients", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := getPatients(db)

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := `[{"id":1,"name":"John Doe","email":"johndoe@example.com","status":"Healthy"},{"id":2,"name":"Jane Doe","email":"janedoe@example.com","status":"Sick"}]`
	if strings.TrimSpace(rr.Body.String()) != strings.TrimSpace(expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
func TestCreatePatient(t *testing.T) {
	// Create a mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Mock the query with expected error message
	mock.ExpectExec("^INSERT INTO patients (name, email, status) VALUES ($1, $2, $3) RETURNING id$").
		WithArgs("John Doe", "johndoe@example.com", "Healthy").
		WillReturnError(fmt.Errorf("some database error"))

	// Create a request to pass to our handler
	body := strings.NewReader(`{"name":"John Doe","email":"johndoe@example.com","status":"Healthy"}`)
	req, err := http.NewRequest("POST", "/patients", body)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := createPatient(db)

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code (should be Internal Server Error)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	// Check if the response body contains the expected error message
	expectedErrMsg := "some database error"
	if !strings.Contains(rr.Body.String(), expectedErrMsg) {
		t.Errorf("handler response body doesn't contain error message: got %v, want to contain: %v", rr.Body.String(), expectedErrMsg)
	}
}

func TestUpdatePatient(t *testing.T) {
	// Create a mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Mock the updated query with capture group for dynamic ID
	mock.ExpectExec("^UPDATE patients SET name = $1, email = $2, status = $3 WHERE id = (.+)$").
		WithArgs("John Doe", "johndoe@example.com", "Healthy").
		WillReturnResult(sqlmock.NewResult(1, 1)) // Update number of affected rows

	// Create a request to pass to our handler
	body := strings.NewReader(`{"name":"John Doe","email":"johndoe@example.com","status":"Healthy"}`)
	req, err := http.NewRequest("PUT", "/patients/1", body)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := updatePatient(db)

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code (should be OK)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body (should be empty as no data is returned)
	expected := ""
	if strings.TrimSpace(rr.Body.String()) != strings.TrimSpace(expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestDeletePatient(t *testing.T) {
	// Create a mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Mock the query
	mock.ExpectExec("^DELETE FROM patients WHERE name IS NULL AND email IS NULL AND status IS NULL OR status = 'Completed'$").
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Create a request to pass to our handler
	req, err := http.NewRequest("DELETE", "/patients", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := deletePatient(db)

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := `"Patients with NULL fields or status 'Completed' deleted"`
	if strings.TrimSpace(rr.Body.String()) != strings.TrimSpace(expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
