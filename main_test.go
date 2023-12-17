package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)
func TestPostMethod(t *testing.T) {
    // Switch to test mode so you don't get such noisy output
    gin.SetMode(gin.TestMode)

    // Setup your router, just like you did in your main function, and
    // register your routes
    r := gin.Default()
    r.POST("/", PostMethod)

    // Create the mock request you'd like to test. Make sure the second argument
    // here is the same as one of the routes you defined in the router setup
    // block!
    req, err := http.NewRequest(http.MethodPost, "/", nil)
    if err != nil {
        t.Fatalf("Couldn't create request: %v\n", err)
    }

    // Create a response recorder so you can inspect the response
    w := httptest.NewRecorder()

    // Perform the request
    r.ServeHTTP(w, req)
    fmt.Println(w.Body)

    // Check to see if the response was what you expected
    if w.Code == http.StatusOK {
        t.Logf("Expected to get status %d is same ast %d\n", http.StatusOK, w.Code)
    } else {
        t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
    }
}

func TestGetMethod(t *testing.T) {
    // Switch to test mode so you don't get such noisy output
    gin.SetMode(gin.TestMode)

    // Setup your router, just like you did in your main function, and
    // register your routes
    r := gin.Default()
    r.GET("/", GetMethod)

    // Create the mock request you'd like to test. Make sure the second argument
    // here is the same as one of the routes you defined in the router setup
    // block!
    req, err := http.NewRequest(http.MethodGet, "/", nil)
    if err != nil {
        t.Fatalf("Couldn't create request: %v\n", err)
    }

    // Create a response recorder so you can inspect the response
    w := httptest.NewRecorder()

    // Perform the request
    r.ServeHTTP(w, req)
    fmt.Println(w.Body)

    // Check to see if the response was what you expected
    if w.Code == http.StatusOK {
        t.Logf("Expected to get status %d is same ast %d\n", http.StatusOK, w.Code)
    } else {
        t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
    }
}