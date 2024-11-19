package sms

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
	"io"
)

// Mocking the Africa's Talking API for testing
func mockAfricaTalkingAPI(t *testing.T) *httptest.Server {
	// Create a mock HTTP server
	handler := http.NewServeMux()
	handler.HandleFunc("/version1/messaging", func(w http.ResponseWriter, r *http.Request) {
		// Assert the expected request method and parameters
		assert.Equal(t, "POST", r.Method)

		// Read the form data from the request
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)

		// Assert the required form parameters
		assert.Contains(t, string(body), "username=sandbox")
		assert.Contains(t, string(body), "to=+254728333926,")
		assert.Contains(t, string(body), "message=Test Message")
		assert.Contains(t, string(body), "from=Sender")

		// Send a mock successful response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"SMSMessageData": {
				"Recipients": [{
					"number": "254728333926",
					"cost": "100",
					"status": "Success",
					"messageId": "msg_id"
				}]
			}
		}`))
	})
	return httptest.NewServer(handler)
}

func TestSendSMS(t *testing.T) {
	// Set up the mock server
	mockServer := mockAfricaTalkingAPI(t)
	defer mockServer.Close()

	// Call the SendSMS function
	SendSMS([]string{"254728333926", }, "Test Message", "Sender")
}
