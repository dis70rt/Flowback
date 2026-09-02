package handlers

import (
	"net/http"

	"github.com/dis70rt/flowback/internal/repo"
	"github.com/gin-gonic/gin"
)

type MetricsHandlerStruct struct {
	queries *repo.Queries
}

func NewMetricsHandler(q *repo.Queries) *MetricsHandlerStruct {
	return &MetricsHandlerStruct{queries: q}
}

func (h *MetricsHandlerStruct) GetOverview(c *gin.Context) {
	overview, err := h.queries.GetDashboardOverview(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch overview"})
		return
	}

	// Success rate = recovered / (recovered + failed/unrecoverable) * 100
	// This correctly excludes still-active cases from the denominator.
	totalClosed := overview.RecoveredCasesCount + overview.FailedCasesCount
	successRate := 0.0
	if totalClosed > 0 {
		successRate = float64(overview.RecoveredCasesCount) / float64(totalClosed) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"total_amount_at_risk":   overview.TotalAmountAtRisk,
		"total_amount_recovered": overview.TotalAmountRecovered,
		"active_cases":           overview.ActiveCasesCount,
		"recovered_cases":        overview.RecoveredCasesCount,
		"failed_cases":           overview.FailedCasesCount,
		"ai_success_rate":        successRate,
	})
}

func (h *MetricsHandlerStruct) GetTrends(c *gin.Context) {
	trends, err := h.queries.GetRecoveryTrends(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch trends"})
		return
	}
	
	if trends == nil {
		trends = []repo.GetRecoveryTrendsRow{}
	}
	c.JSON(http.StatusOK, trends)
}

func (h *MetricsHandlerStruct) GetChannels(c *gin.Context) {
	channels, err := h.queries.GetChannelDistribution(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch channels"})
		return
	}
	
	if channels == nil {
		channels = []repo.GetChannelDistributionRow{}
	}
	c.JSON(http.StatusOK, channels)
}

func (h *MetricsHandlerStruct) GetPipeline(c *gin.Context) {
	pipeline, err := h.queries.GetPipelineStatus(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch pipeline"})
		return
	}
	
	if pipeline == nil {
		pipeline = []repo.GetPipelineStatusRow{}
	}
	c.JSON(http.StatusOK, pipeline)
}

func (h *MetricsHandlerStruct) GetRecoveredCases(c *gin.Context) {
	cases, err := h.queries.ListRecoveredCases(c.Request.Context(), repo.ListRecoveredCasesParams{
		Limit:  50,
		Offset: 0,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch recovered cases"})
		return
	}

	if cases == nil {
		cases = []repo.ListRecoveredCasesRow{}
	}
	c.JSON(http.StatusOK, cases)
}
