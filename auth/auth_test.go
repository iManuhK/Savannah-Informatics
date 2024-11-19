package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Mocked token and user data for testing
var mockedValidToken = "mocked-valid-token"

// Mock authentication middleware to bypass actual authentication flow
func MockOIDCAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "Bearer "+mockedValidToken {
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
		}
	}
}

// Mock authorization
func ProtectedRouteHandler(c *gin.Context) {
	c.String(http.StatusOK, "Access granted")
}

// Test case for a valid token
func TestProtectedRoute_WithMockedAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Apply the mocked authentication middleware
	router.Use(MockOIDCAuthMiddleware())
	router.GET("/protected", ProtectedRouteHandler)

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+mockedValidToken) // Set the mocked token in the header

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Access granted")
}

// Test case for invalid token
func TestProtectedRoute_WithInvalidAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.Use(MockOIDCAuthMiddleware())
	router.GET("/protected", ProtectedRouteHandler)

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert that the status code is Unauthorized (401) and the response body contains the error message
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid token")
}
