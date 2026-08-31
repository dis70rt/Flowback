package handlers

import (
	"net/http"
	"github.com/dis70rt/flowback/internal/repo"
	"github.com/gin-gonic/gin"
)

func MetricsHandler(queries *repo.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		//TODO: Return dummy data for now to satisfy the frontend
		c.JSON(http.StatusOK, gin.H{
			"total_revenue_recovered": 12500.50,
			"active_cases": 14,
			"ai_success_rate": 87.5,
		})
	}
}
