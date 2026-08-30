package seed

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/razorpay/razorpay-go"
)

// FailedCardTask creates an order and fails a Card payment
type FailedCardTask struct {
	CustomerID  string
	Email       string
	AmountPaise int
	CardNumber  string
}

func (t *FailedCardTask) Name() string {
	return fmt.Sprintf("Failed Card Payment (Customer: %s, Card: %s)", t.CustomerID, t.CardNumber)
}

func (t *FailedCardTask) Execute(client *razorpay.Client, keyID string) {
	log.Printf("Starting Task: Failed Card Payment (Customer: %s, Amount: %d paise)\n", t.CustomerID, t.AmountPaise)

	// 1. Create Order
	data := map[string]interface{}{
		"amount":   t.AmountPaise,
		"currency": "INR",
		"receipt":  fmt.Sprintf("sim_card_%d", time.Now().Unix()),
	}
	// Note: We don't strictly need to pass customer_id to the Order for this simulation, 
	// but we MUST pass it to the payment API so Razorpay includes it in the webhook!

	order, err := client.Order.Create(data, nil)
	if err != nil {
		log.Fatalf("Failed to create order: %v", err)
	}

	orderID := order["id"].(string)
	log.Printf("Order Created: %s\n", orderID)

	// 2. Fire S2S Payment API
	formData := url.Values{}
	formData.Set("key_id", keyID)
	formData.Set("order_id", orderID)
	formData.Set("amount", fmt.Sprintf("%d", t.AmountPaise))
	formData.Set("currency", "INR")
	formData.Set("method", "card")
	formData.Set("card[number]", t.CardNumber)
	formData.Set("card[expiry_month]", TestCardExpiryMonth)
	formData.Set("card[expiry_year]", TestCardExpiryYear)
	formData.Set("card[cvv]", TestCardCVV)
	formData.Set("contact", "9999999999")
	formData.Set("email", t.Email)

	executeAJAXPayment(formData)
}

func executeAJAXPayment(formData url.Values) {
	req, err := http.NewRequest("POST", "https://api.razorpay.com/v1/payments/create/ajax", strings.NewReader(formData.Encode()))
	if err != nil {
		log.Fatalf("Failed to build HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpClient := &http.Client{Timeout: 10 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to execute payment API: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// 3. Parse JSON for redirect URL
	var res map[string]interface{}
	if err := json.Unmarshal(body, &res); err == nil {
		if reqData, ok := res["request"].(map[string]interface{}); ok {
			if authURL, ok := reqData["url"].(string); ok {
				log.Printf("Action Required: Complete the OTP step at the link below:\n")
				log.Printf("-> %s\n", authURL)
				return
			}
		}
	}

	// Fallback print
	log.Printf("Razorpay Response: %s\n", string(body))
}
