package api

import (
	_ "embed"

	"github.com/gin-gonic/gin"

	"github.com/dis70rt/flowback/internal/api/handlers"
	"github.com/dis70rt/flowback/internal/events"
	"github.com/dis70rt/flowback/internal/pubsub"
	"github.com/dis70rt/flowback/internal/razorpay"
	"github.com/dis70rt/flowback/internal/repo"
)

//go:embed swagger.yaml
var swaggerYAML []byte

type RouterDeps struct {
	Queries        *repo.Queries
	Enqueuer       *events.Enqueuer
	Bus            pubsub.Bus
	RazorpaySecret string
	RazorpayClient *razorpay.Client
}

func NewRouter(deps RouterDeps) *gin.Engine {
	r := gin.Default()

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"service": "flowback-api",
		})
	})

	// External Webhooks
	rzpHandler := razorpay.NewWebhookHandler(deps.RazorpaySecret, deps.Enqueuer, deps.Queries)
	r.POST("/webhooks/razorpay", rzpHandler.Handle)

	// API Group
	apiGroup := r.Group("/api")
	{
		// Serve OpenAPI documentation
		apiGroup.GET("/docs/openapi.yaml", func(c *gin.Context) {
			c.Data(200, "application/yaml", swaggerYAML)
		})
		apiGroup.GET("/docs", func(c *gin.Context) {
			c.Data(200, "text/html; charset=utf-8", []byte(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Flowback API Docs</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui.css" />
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-bundle.js"></script>
    <script>
        window.onload = () => {
            window.ui = SwaggerUIBundle({
                url: '/api/docs/openapi.yaml',
                dom_id: '#swagger-ui',
            });
        };
    </script>
</body>
</html>
			`))
		})

		// Real-time SSE Stream
		apiGroup.GET("/stream", handlers.StreamHandler(deps.Bus))

		// Metrics
		apiGroup.GET("/metrics", handlers.MetricsHandler(deps.Queries))

		// Cases
		caseHandler := handlers.NewCaseHandler(deps.Queries, deps.RazorpayClient)
		apiGroup.GET("/cases", caseHandler.ListCases)
		apiGroup.GET("/cases/:id", caseHandler.GetCase)
		apiGroup.GET("/cases/:id/summary", caseHandler.GetCaseSummary)
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
