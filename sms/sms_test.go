package sms

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestSendSMS(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
		assert.Equal(t, ApiKey, r.Header.Get("apiKey"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"SMSMessageData": {"recipients": [{"number": "+254728333926", "cost": "KES 0.5", "status": "Success", "messageId": "abc123"}]}}`))
	}))
	defer server.Close()

	// Override API URL
	UrlStr = server.URL

	recipients := []string{"+254728333926"}
	message := "Test SMS"
	senderID := "SAVANNAH INF"

	SendSMS(recipients, message, senderID)

	// Add assertions or mocks for detailed validation if needed
}