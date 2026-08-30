package seed

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
)

var (
	vipID     = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	genzID    = uuid.MustParse("22222222-2222-2222-2222-222222222222")
	outageID  = uuid.MustParse("33333333-3333-3333-3333-333333333333")
	ignorerID = uuid.MustParse("44444444-4444-4444-4444-444444444444")
	fraudID   = uuid.MustParse("55555555-5555-5555-5555-555555555555")
	nightID   = uuid.MustParse("66666666-6666-6666-6666-666666666666")
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
}

// RunDBSeeder wipes the customers table and injects our 6 demo archetypes.
func RunDBSeeder(db *sql.DB) {
	log.Println(" Cleaning old customer data...")
	_, err := db.Exec(`TRUNCATE TABLE customers CASCADE;`)
	if err != nil {
		log.Fatal("Failed to truncate:", err)
	}

	log.Println(" Seeding 6 Customer Archetypes into DB...")

	customers := []CustomerSeed{
		{vipID, "cust_vip", "Rajesh (VIP)", "HIGH", "LOYAL", "EMAIL", "Bangalore", "Karnataka", 3, 0.99},
		{genzID, "cust_genz", "Aditi (Gen Z)", "LOW", "NEW", "WHATSAPP", "Delhi", "Delhi", 0, 0.80},
		{outageID, "cust_outage", "Vikram (Outage)", "MEDIUM", "ESTABLISHED", "SMS", "Kathmandu", "Bagmati", 1, 0.95},
		{ignorerID, "cust_ignorer", "Priya (Ignorer)", "HIGH", "ESTABLISHED", "EMAIL", "Pune", "Maharashtra", 2, 0.90},
		{fraudID, "cust_fraud", "Scammer", "LOW", "NEW", "EMAIL", "Unknown", "Unknown", 0, 0.10},
		{nightID, "cust_night", "Rahul (Night)", "MEDIUM", "NEW", "SMS", "Hyderabad", "Telangana", 1, 0.85},
	}

	for _, c := range customers {
		_, err := db.Exec(`
			INSERT INTO customers (id, razorpay_customer_id, name, value_tier, tenure, preferred_channel, city, state, failed_payments, reliability_score)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`, c.ID, c.RazorpayID, c.Name, c.Tier, c.Tenure, c.Channel, c.City, c.State, c.FailedCount, c.Reliability)
		if err != nil {
			log.Fatalf("Failed to insert %s: %v", c.Name, err)
		}
	}

	log.Println(" Customers seeded successfully! Now run:")
	log.Println("  go run cmd/seeder/main.go --task=vip")
	log.Println("  go run cmd/seeder/main.go --task=outage")
}
