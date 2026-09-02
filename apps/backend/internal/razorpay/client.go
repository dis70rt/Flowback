package razorpay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	keyID     string
	keySecret string
}

func NewClient(keyID, secret string) *Client {
	return &Client{keyID: keyID, keySecret: secret}
}

type CreatePaymentLinkRequest struct {
	Amount      int64             `json:"amount"`
	Currency    string            `json:"currency"`
	Description string            `json:"description,omitempty"`
	Notes       map[string]string `json:"notes,omitempty"`
}

type CreatePaymentLinkResponse struct {
	ID       string `json:"id"`
	ShortURL string `json:"short_url"`
}

func (c *Client) CreatePaymentLink(amount int64, caseID string) (string, string, error) {
	reqBody := CreatePaymentLinkRequest{
		Amount:   amount,
		Currency: "INR",
		Notes: map[string]string{
			"recovery_case_id": caseID,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", "", err
	}

	req, err := http.NewRequest("POST", "https://api.razorpay.com/v1/payment_links", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", err
	}

	req.SetBasicAuth(c.keyID, c.keySecret)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("razorpay api error: status=%d body=%s", resp.StatusCode, string(bodyBytes))
	}

	var parsedResp CreatePaymentLinkResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsedResp); err != nil {
		return "", "", err
	}

	return parsedResp.ID, parsedResp.ShortURL, nil
}
