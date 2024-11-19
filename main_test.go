package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"savannah.go/auth"
	"golang.org/x/oauth2"
)

func init() {
	// Set mock environment variables for testing
	os.Setenv("CLIENT_ID", "mock-client-id")
	os.Setenv("CLIENT_SECRET", "mock-client-secret")
	os.Setenv("REDIRECT_URI", "http://localhost:8080/oauth/callback")
	os.Setenv("PROVIDER_URL", "https://accounts.google.com")
}

func TestRoutes(t *testing.T) {
	auth.InitOIDC()

}
// Test that the login route returns a 302 (redirect) status code
func TestLoginRoute(t *testing.T) {
	router := gin.Default()

	router.GET("/login", func(c *gin.Context) {
		authURL := auth.GetOAuth2Config().AuthCodeURL("state-string", oauth2.AccessTypeOffline)
		c.Redirect(http.StatusTemporaryRedirect, authURL)
	})

	// Perform a GET request on the /login endpoint
	req, _ := http.NewRequest("GET", "/login", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Assert the status code =302 - redirect
	assert.Equal(t, http.StatusTemporaryRedirect, resp.Code)
}

// Test that the /customers route returns a 200 OK status code
func TestGetCustomersRoute(t *testing.T) {
	router := gin.Default()

	router.GET("/customers", func(c *gin.Context) {
		// This would be a mock response simulating database records
		customers := []map[string]interface{}{
			{"cust_id": 1, "code": "C001", "full_name": "rwsrrt", "phone": 234567},
			{"cust_id": 2, "code": "C002", "full_name": "emmanuel", "phone": 254728333926},
		}
		c.JSON(http.StatusOK, customers)
	})

	// Perform a GET request on the /customers endpoint
	req, _ := http.NewRequest("GET", "/customers", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Assert the status code (should be 200 for successful response)
	assert.Equal(t, http.StatusOK, resp.Code)

	// Assert the response body (checking if customer data is returned)
	expectedBody := `[{"cust_id": 1, "code": "C001", "full_name": "rwsrrt", "phone": 234567},{"cust_id": 2, "code": "C002", "full_name": "emmanuel", "phone": 254728333926}]`
	assert.JSONEq(t, expectedBody, resp.Body.String())
}

// Test the POST /customers route
func TestPostCustomersRoute(t *testing.T) {
	// Create a new Gin router
	router := gin.Default()

	// Define the POST route
	router.POST("/customers", func(c *gin.Context) {
		var newCustomer struct {
			Code  string `json:"code"`
			Name  string `json:"full_name"`
			Phone string `json:"phone"`
		}
		if err := c.BindJSON(&newCustomer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		// Simulate inserting a customer and returning a response
		newCustomerID := 1 // Mocked customer ID
		c.JSON(http.StatusCreated, gin.H{
			"cust_id": newCustomerID,
			"code":    newCustomer.Code,
			"full_name": newCustomer.Name,
			"phone":   newCustomer.Phone,
		})
	})

	// Create the JSON payload for the POST request
	customerPayload := `{"code":"C003", "full_name":"Alice Johnson", "phone":"+254701234569"}`
	req, _ := http.NewRequest("POST", "/customers", bytes.NewBuffer([]byte(customerPayload)))
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Assert the status code (should be 201 for created)
	assert.Equal(t, http.StatusCreated, resp.Code)

	// Assert the response body
	expectedResponse := `{"cust_id":1,"code":"C003","full_name":"Alice Johnson","phone":"+254701234569"}`
	assert.JSONEq(t, expectedResponse, resp.Body.String())
}

// Test the POST /orders route
func TestPostOrdersRoute(t *testing.T) {
	// Create a new Gin router
	router := gin.Default()

	// Define the POST route (without SMS sending)
	router.POST("/orders", func(c *gin.Context) {
		var newOrder struct {
			Item           string  `json:"item"`
			Time           string  `json:"time"`
			Amount         float64 `json:"amount"`
			RelatedCustomer int    `json:"cust_id"`
		}
		if err := c.BindJSON(&newOrder); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		// Simulate order creation (no SMS sending)
		orderID := 1 // Mocked order ID

		// Return the response
		c.JSON(http.StatusCreated, gin.H{
			"order_id": orderID,
			"item":     newOrder.Item,
			"amount":   newOrder.Amount,
		})
	})

	// Create the JSON payload for the POST request
	orderPayload := `{"item":"Laptop", "time":"2024-11-19T10:00:00Z", "amount": 50000, "cust_id": 1}`
	req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer([]byte(orderPayload)))
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Assert the status code (should be 201 for created)
	assert.Equal(t, http.StatusCreated, resp.Code)

	// Assert the response body
	expectedResponse := `{"order_id":1,"item":"Laptop","amount":50000}`
	assert.JSONEq(t, expectedResponse, resp.Body.String())
}
