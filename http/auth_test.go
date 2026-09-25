package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mazay/mikromanager/db"
	"go.uber.org/zap"
)

func TestHttpConfigLogin(t *testing.T) {
	// Create a test http config
	logger := zap.NewNop()
	db := &db.DB{}
	
	config := HttpConfig{
		Db:            db,
		EncryptionKey: "test-key",
		Logger:        logger,
	}

	// Create a test request with POST method but no form data
	req, err := http.NewRequest("POST", "/login", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	// Add content type header to simulate form submission
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a ResponseWriter
	rr := httptest.NewRecorder()

	// Call the handler
	config.login(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, status)
	}
}

func TestHttpConfigLogout(t *testing.T) {
	logger := zap.NewNop()
	db := &db.DB{}
	
	config := HttpConfig{
		Db:            db,
		EncryptionKey: "test-key",
		Logger:        logger,
	}

	// Create a test request with GET method
	req, err := http.NewRequest("GET", "/logout", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	// Add a session cookie
	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: "test-session-id",
	})

	// Create a ResponseWriter
	rr := httptest.NewRecorder()

	// Call the handler
	config.logout(rr, req)

	// Check that we redirect (status code 302)
	if status := rr.Code; status != http.StatusFound {
		t.Errorf("Expected status %d, got %d", http.StatusFound, status)
	}
}