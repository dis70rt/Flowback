package seed

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

var (
	vipID     = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	genzID    = uuid.MustParse("22222222-2222-2222-2222-222222222222")
	outageID  = uuid.MustParse("33333333-3333-3333-3333-333333333333")
	ignorerID = uuid.MustParse("44444444-4444-4444-4444-444444444444")
	fraudID   = uuid.MustParse("55555555-5555-5555-5555-555555555555")
	nightID   = uuid.MustParse("66666666-6666-6666-6666-666666666666")
	ceoID     = uuid.MustParse("77777777-7777-7777-7777-777777777777")
)

type CustomerSeed struct {
	ID          uuid.UUID
	RazorpayID  string
	Name        string
	Tier        string
	Tenure      string
	Channel     string
	City        string
	State       string
	FailedCount int
	Reliability float64
	Email       string
	Phone       string
}

func createRealRazorpayCustomer(c CustomerSeed, keyID, keySecret string) (string, error) {
	if keyID == "" || keySecret == "" {
		return c.RazorpayID, nil // Fallback to hardcoded mock if no keys provided
	}

	reqBody := map[string]string{
		"name":          c.Name,
		"email":         c.Email,
		"contact":       c.Phone,
		"fail_existing": "0",
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "https://api.razorpay.com/v1/customers", bytes.NewBuffer(body))
	req.SetBasicAuth(keyID, keySecret)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var rzpResp struct {
		ID string `json:"id"`
	}
	json.Unmarshal(respBody, &rzpResp)

	if rzpResp.ID != "" {
		return rzpResp.ID, nil
	}
	return c.RazorpayID, nil
}

func RunDBSeeder(db *sql.DB, rzpKeyID, rzpKeySecret string) {
	log.Println(" Cleaning old customer data...")
	_, err := db.Exec(`TRUNCATE TABLE customers CASCADE;`)
	if err != nil {
		log.Fatal("Failed to truncate:", err)
	}

	log.Println(" Seeding 6 Customer Archetypes into DB (Fetching real RZP IDs if keys exist)...")

	customers := []CustomerSeed{
		{vipID, "cust_vip", "Rajesh (VIP)", "HIGH", "LOYAL", "EMAIL", "Bangalore", "Karnataka", 3, 0.99, "rajesh.vip@example.com", "919000000001"},
		{genzID, "cust_genz", "Aditi (Gen Z)", "LOW", "NEW", "WHATSAPP", "Delhi", "Delhi", 0, 0.80, "aditi.genz@example.com", "919000000002"},
		{outageID, "cust_outage", "Vikram (Outage)", "MEDIUM", "ESTABLISHED", "SMS", "Kathmandu", "Bagmati", 1, 0.95, "vikram.outage@example.com", "919000000003"},
		{ignorerID, "cust_ignorer", "Priya (Ignorer)", "HIGH", "ESTABLISHED", "EMAIL", "Pune", "Maharashtra", 2, 0.90, "priya.ignorer@example.com", "919000000004"},
		{fraudID, "cust_fraud", "Scammer", "LOW", "NEW", "EMAIL", "Unknown", "Unknown", 0, 0.10, "scammer@example.com", "919000000005"},
		{nightID, "cust_night", "Rahul (Night)", "MEDIUM", "NEW", "SMS", "Hyderabad", "Telangana", 1, 0.85, "rahul.night@example.com", "919000000006"},
		{ceoID, "cust_ceo", "Enterprise CEO", "HIGH", "LOYAL", "EMAIL", "Mumbai", "Maharashtra", 5, 0.20, "enterprise.ceo@example.com", "919000000007"},
	}

	for _, c := range customers {
		realRzpID, err := createRealRazorpayCustomer(c, rzpKeyID, rzpKeySecret)
		if err != nil {
			log.Printf("Warning: Failed to fetch real razorpay ID for %s, falling back to mock: %v", c.Name, err)
			realRzpID = c.RazorpayID
		}

		_, err = db.Exec(`
			INSERT INTO customers (id, razorpay_customer_id, name, email, phone, value_tier, tenure, preferred_channel, city, state, failed_payments, reliability_score)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		`, c.ID, realRzpID, c.Name, c.Email, c.Phone, c.Tier, c.Tenure, c.Channel, c.City, c.State, c.FailedCount, c.Reliability)
		
		if err != nil {
			log.Fatalf("Failed to insert %s: %v", c.Name, err)
		}
		fmt.Printf("[+] Seeded: %s (RZP: %s)\n", c.Name, realRzpID)
	}

	log.Println(" Customers seeded successfully! Now run:")
	log.Println("  go run cmd/seeder/main.go --task=vip")
	log.Println("  go run cmd/seeder/main.go --task=outage")
}
