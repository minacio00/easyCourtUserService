package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/minacio00/easyCourtUserService/database"
	"github.com/stretchr/testify/assert"
)

func TestHashPasswordHandler(t *testing.T) {
	t.Run("Valid password hashing", func(t *testing.T) {
		// Prepare a request body with valid credentials
		body := map[string]string{
			"email":    "testUser@gmail.com",
			"password": "testUser",
		}
		reqBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/hash-password", bytes.NewReader(reqBody))
		w := httptest.NewRecorder()

		// Call the handler
		hashPassword(w, req)

		// Check the response status code
		assert.Equal(t, http.StatusOK, w.Code)

		// Additional checks if needed for the response body or headers
		// ...
	})

	t.Run("Empty password handling", func(t *testing.T) {
		// Prepare a request body with an empty password
		body := map[string]string{
			"email":    "testUser@gmail.com",
			"password": "",
		}
		reqBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/hash-password", bytes.NewReader(reqBody))
		w := httptest.NewRecorder()

		// Call the handler
		hashPassword(w, req)

		// Check the response status code
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Additional checks if needed for the response body or headers
		// ...
	})
}

func TestSigningHandler(t *testing.T) {
	// Mock your database or set up a test database
	os.Setenv("APP_ENV", "test")
	log.Println("Do stuff BEFORE the tests!")
	database.Connectdb()
	// Prepare a request body with valid credentials
	tenant := &Tenant{
		Email:    "testUser@gmail.com",
		Password: "testUser",
	}
	database.Db.Save(&tenant)

	// não está passando pq não estou hasheando a senha pra salvar
	t.Run("Valid sign-in", func(t *testing.T) {
		body := map[string]string{
			"email":    "testUser@gmail.com",
			"password": "testUser",
		}
		reqBody, _ := json.Marshal(body)
		req := httptest.NewRequest("POST", "/sign-in", bytes.NewReader(reqBody))
		w := httptest.NewRecorder()

		// Mock user data in your test database
		// ...

		// Call the handler
		signing(w, req)

		// Check the response status code
		l := log.Default()
		l.Printf("%v \n", w.Code)
		assert.Equal(t, http.StatusAccepted, w.Code)

		// Additional checks if needed for the response body or headers
		// ...
	})

	t.Run("Empty email or password handling", func(t *testing.T) {
		// Prepare a request body with empty email and/or password
		body := map[string]string{
			"email":    "",
			"password": "test123",
		}
		reqBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/sign-in", bytes.NewReader(reqBody))
		w := httptest.NewRecorder()

		// Call the handler
		signing(w, req)

		// Check the response status code
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Additional checks if needed for the response body or headers
		// ...
	})
}
