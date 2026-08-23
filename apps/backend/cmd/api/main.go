package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/dis70rt/flowback/internal/config"
	"github.com/dis70rt/flowback/internal/razorpay"
)

func main() {
	cfg := config.Load()
	r := gin.Default()

	rzpHandler := razorpay.NewWebhookHandler(cfg.RazorpaySecret)

	r.POST("/webhooks/razorpay", rzpHandler.Handle)

	log.Println("STARTING: Flowback Backend listening on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
