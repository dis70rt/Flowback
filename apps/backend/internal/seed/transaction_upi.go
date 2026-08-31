package seed

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/razorpay/razorpay-go"
)

// FailedUPICollectTask creates an order and fails a UPI payment
type FailedUPICollectTask struct {
	AmountPaise int
}

func (t *FailedUPICollectTask) Name() string {
	return fmt.Sprintf("Failed UPI Collect (Amount: %d paise)", t.AmountPaise)
}

func (t *FailedUPICollectTask) Execute(client *razorpay.Client, keyID string) {
	// 1. Create Order
	data := map[string]interface{}{
		"amount":   t.AmountPaise,
		"currency": "INR",
		"receipt":  fmt.Sprintf("sim_upi_%d", time.Now().Unix()),
	}

	order, err := client.Order.Create(data, nil)
	if err != nil {
		log.Fatalf("[FATAL] Failed to create order: %v", err)
	}

	orderID := order["id"].(string)
	log.Printf("[+] Order Created: %s", orderID)

	formData := url.Values{}
	formData.Set("key_id", keyID)
	formData.Set("order_id", orderID)
	formData.Set("amount", fmt.Sprintf("%d", t.AmountPaise))
	formData.Set("currency", "INR")
	formData.Set("method", "upi")
	formData.Set("upi[flow]", "collect")
	formData.Set("upi[vpa]", "failure@razorpay")
	formData.Set("contact", "9999999999")
	formData.Set("email", "test@flowback.ai")

	executeAJAXPayment(formData)
}
