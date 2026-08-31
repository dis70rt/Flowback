package api

import (
	"github.com/gin-gonic/gin"

	"github.com/dis70rt/flowback/internal/api/handlers"
	"github.com/dis70rt/flowback/internal/events"
	"github.com/dis70rt/flowback/internal/pubsub"
	"github.com/dis70rt/flowback/internal/razorpay"
	"github.com/dis70rt/flowback/internal/repo"
)

type RouterDeps struct {
	Queries        *repo.Queries
	Enqueuer       *events.Enqueuer
	Bus            pubsub.Bus
	RazorpaySecret string
}

func NewRouter(deps RouterDeps) *gin.Engine {
	r := gin.Default()

	// External Webhooks
	rzpHandler := razorpay.NewWebhookHandler(deps.RazorpaySecret, deps.Enqueuer)
	r.POST("/webhooks/razorpay", rzpHandler.Handle)

	// API Group
	apiGroup := r.Group("/api")
	{
		// Real-time SSE Stream
		apiGroup.GET("/stream", handlers.StreamHandler(deps.Bus))

		// Metrics
		apiGroup.GET("/metrics", handlers.MetricsHandler(deps.Queries))

		// Cases
		caseHandler := handlers.NewCaseHandler(deps.Queries)
		apiGroup.GET("/cases", caseHandler.ListCases)
		apiGroup.GET("/cases/:id", caseHandler.GetCase)
		apiGroup.POST("/cases/:id/approve", caseHandler.ApproveDraft)
		apiGroup.PUT("/cases/:id/draft", caseHandler.EditDraft)
		apiGroup.POST("/cases/:id/reject", caseHandler.RejectDraft)

		// Customers
		customerHandler := handlers.NewCustomerHandler(deps.Queries)
		apiGroup.GET("/customers/:id", customerHandler.GetCustomer)
		apiGroup.GET("/customers/:id/payments", customerHandler.GetPayments)
		apiGroup.GET("/customers/:id/communications", customerHandler.GetCommunications)
	}

	return r
}
