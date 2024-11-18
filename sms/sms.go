package sms

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"io"
)

const apiKey   = "atsk_4e44fbe3d8d1140a565c6139fb64edc19f183bc3e64f7ebccdef7cc193bc97c3d35b3434"
const username = "sandbox"


type SMSResponse struct {
	SMSMessageData struct {
		Recipients []struct {
			Number    string `json:"number"`
			Cost      string `json:"cost"`
			Status    string `json:"status"`
			MessageID string `json:"messageId"`
		} `json:"recipients"`
	} `json:"SMSMessageData"`
}

func SendSMS(recipients []string, message, senderID string) {
	urlStr := "https://api.sandbox.africastalking.com/version1/messaging"

	to := strings.Join(recipients, ",")

	// Set up form values
	data := url.Values{}
	data.Set("username", username)
	data.Set("to", to)
	data.Set("message", message)
	data.Set("from", senderID)

	req, err := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatalf("Failed to create HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("apiKey", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send SMS: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("Non-200 status code: %d. Response: %s", resp.StatusCode, string(body))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	fmt.Println("Raw Response:", string(respBody))

}
