package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"strings"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	// You should import the correct database driver (SQLite for mocking, if using SQLite)
	_ "github.com/mattn/go-sqlite3"
)

func TestGetCustomers(t *testing.T) {
	// Mock database connection
	DB, _ := setupMockDB()  // Setup mock DB
	defer DB.Close()

	// Set the DB in your main package (or use dependency injection if applicable)
	main.DB = DB

	router := gin.Default()
	router.GET("/customers", main.GetCustomers)

	req := httptest.NewRequest("GET", "/customers", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "[]") // Assuming empty result set
}

func TestPostCustomers(t *testing.T) {
	// Mock database connection
	DB, _ := setupMockDB()  // Setup mock DB
	defer DB.Close()

	// Set the DB in your main package (or use dependency injection if applicable)
	main.DB = DB

	router := gin.Default()
	router.POST("/customers", main.PostCustomers)

	body := `{"code": "123", "full_name": "John Doe", "phone": "+1234567890"}`
	req := httptest.NewRequest("POST", "/customers", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), `"code":"123"`)
}

// Mock database setup function
func setupMockDB() (*sql.DB, error) {
	// Use an in-memory SQLite database for testing
	// SQLite creates a temporary in-memory database that is discarded after the test
	db, err := sql.Open("sqlite3", ":memory:") 
	if err != nil {
		return nil, err
	}

	// Create tables and mock data if needed
	_, err = db.Exec(`
		CREATE TABLE customers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			code TEXT,
			full_name TEXT,
			phone TEXT
		);
	`)
	if err != nil {
		return nil, err
	}

	// Insert mock data if needed
	_, err = db.Exec(`INSERT INTO customers (code, full_name, phone) VALUES ("123", "John Doe", "+1234567890")`)
	if err != nil {
		return nil, err
	}

	return db, nil
}
