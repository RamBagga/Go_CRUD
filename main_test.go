package patients_test

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"your_app/patients"
)

// MockDB defines the interface for our mock database
type MockDB interface {
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
	Exec(string, ...interface{}) (sql.Result, error)
}

func TestGetPatients(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockDB := MockDB{}
	patientsHandler := patients.GetPatients(mockDB)

	// Positive case: Mock successful query with valid patients
	expectedPatients := []patients.Patient{{ID: 1, Name: "John Doe", Email: "john.doe@example.com", Status: "Active"}}
	mockDB.EXPECT().Query(gomock.Eq("SELECT * FROM patients WHERE name != '' AND email != '' AND status != ''"), gomock.Nil()).Return(mockRows(expectedPatients), nil)

	req, _ := http.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	patientsHandler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var actualPatients []patients.Patient
	err := json.NewDecoder(rec.Body).Decode(&actualPatients)
	assert.NoError(t, err)
	assert.Equal(t, expectedPatients, actualPatients)

	// Negative case: Mock error during query
	mockDB.EXPECT().Query(gomock.Eq("SELECT * FROM patients WHERE name != '' AND email != '' AND status != ''"), gomock.Nil()).Return(nil, errors.New("mock error"))

	req, _ = http.NewRequest("GET", "/", nil)
	rec = httptest.NewRecorder()
	patientsHandler(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestGetPatient(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockDB := MockDB{}
	patientHandler := patients.GetPatient(mockDB)

	// Positive case: Mock successful query with valid patient
	expectedPatient := patients.Patient{ID: 1, Name: "John Doe", Email: "john.doe@example.com", Status: "Active"}
	mockDB.EXPECT().QueryRow(gomock.Eq("SELECT * FROM patients WHERE id = $1"), gomock.Eq(1)).Return(mockRow(expectedPatient), nil)

	req, _ := http.NewRequest("GET", "/1", nil)
	rec := httptest.NewRecorder()
	patientHandler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var actualPatient patients.Patient
	err := json.NewDecoder(rec.Body).Decode(&actualPatient)
	assert.NoError(t, err)
	assert.Equal(t, expectedPatient, actualPatient)

	// Negative case: Mock error during query
	mockDB.EXPECT().QueryRow(gomock.Eq("SELECT * FROM patients WHERE id = $1"), gomock.Eq(2)).Return(nil, errors.New("mock error"))

	req, _ = http.NewRequest("GET", "/2", nil)
	rec = httptest.NewRecorder()
	patientHandler(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	// Negative case: Mock no rows found
	mockDB.EXPECT().QueryRow(gomock.Eq("SELECT * FROM patients WHERE id = $1"), gomock.Eq(3)).Return(&sql.Row{}, sql.ErrNoRows)

	req, _ = http.NewRequest("GET", "/3", nil)
	rec = httptest.NewRecorder()
	patientHandler(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// ... Implement similar test cases for createPatient, updatePatient, and deletePatient functions ...
