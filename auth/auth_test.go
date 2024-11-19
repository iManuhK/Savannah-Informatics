package auth


import (
	"testing"
	"os"
	"net/http"
	"net/http/httptest"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Test InitOIDC with missing environment variables
func TestInitOIDC_MissingEnv(t *testing.T) {
    os.Clearenv() 
    defer func() {
        if r := recover(); r != nil {
            t.Log("Recovered from panic as expected:", r)
        } else {
            t.Errorf("Expected panic, but it did not occur")
        }
    }()

    InitOIDC() 
}


// Test OIDC Middleware with valid and invalid tokens
func TestOIDCAuthMiddleware_InvalidToken(t *testing.T) {
	InitOIDC()

	// Set up Gin context
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(OIDCAuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.String(http.StatusOK, "Access granted")
	})

	// Test with missing Authorization header
	req := httptest.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Authorization header missing")

	// Test with invalid token
	req.Header.Set("Authorization", "Bearer invalid-token")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid token")
}
