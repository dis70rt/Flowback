package seed

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/razorpay/razorpay-go"
)

type FailedSubscriptionTask struct {
	CustomerID  string
	Email       string
	AmountPaise int
}

func (t *FailedSubscriptionTask) Name() string {
	return fmt.Sprintf("Failed Subscription (Customer: %s)", t.CustomerID)
}

func (t *FailedSubscriptionTask) Execute(client *razorpay.Client, keyID string) {
	log.Printf("Starting Task: Failed Subscription (Customer: %s, Amount: %d paise)\n", t.CustomerID, t.AmountPaise)

	// 1. Create a Plan
	planData := map[string]interface{}{
		"period":   "monthly",
		"interval": 1,
		"item": map[string]interface{}{
			"name":        fmt.Sprintf("Test Plan %d", time.Now().Unix()),
			"amount":      t.AmountPaise,
			"currency":    "INR",
			"description": "Flowback Test Subscription Plan",
		},
	}
	plan, err := client.Plan.Create(planData, nil)
	if err != nil {
		log.Fatalf("Failed to create plan: %v", err)
	}
	planID := plan["id"].(string)
	log.Printf("Plan Created: %s\n", planID)

	// 2. Create a Subscription
	subData := map[string]interface{}{
		"plan_id":         planID,
		"total_count":     12,
		"customer_id":     t.CustomerID,
		"customer_notify": 1,
		"notes": map[string]interface{}{
			"test_purpose": "flowback_simulation",
		},
	}
	
	sub, err := client.Subscription.Create(subData, nil)
	if err != nil {
		log.Fatalf("Failed to create subscription: %v", err)
	}

	subID := sub["id"].(string)
	log.Printf("Subscription Created: %s\n", subID)

	log.Printf("Attempting to force failed payment via Server-to-Server AJAX request...")
	
	formData := url.Values{}
	formData.Set("key_id", keyID)
	formData.Set("subscription_id", subID)
	formData.Set("amount", fmt.Sprintf("%d", t.AmountPaise))
	formData.Set("currency", "INR")
	formData.Set("method", "card")
	formData.Set("card[number]", CardSubscriptionDomesticVisa)
	formData.Set("card[expiry_month]", TestCardExpiryMonth)
	formData.Set("card[expiry_year]", TestCardExpiryYear)
	formData.Set("card[cvv]", TestCardCVV)
	formData.Set("contact", "9999999999")
	formData.Set("email", t.Email)

	executeAJAXPayment(formData)
}
