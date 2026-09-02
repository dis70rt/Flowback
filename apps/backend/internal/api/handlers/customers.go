package handlers

import (
	"database/sql"
	"net/http"

	"github.com/dis70rt/flowback/internal/repo"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CustomerHandler struct {
	queries *repo.Queries
}

func NewCustomerHandler(q *repo.Queries) *CustomerHandler {
	return &CustomerHandler{queries: q}
}

func (h *CustomerHandler) GetCustomer(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer uuid"})
		return
	}

	customer, err := h.queries.GetCustomerByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) GetPayments(c *gin.Context) {
	// We don't have a payments table synced in this MVP, so we return a placeholder.
	// In production, this would hit Razorpay API or a local 'payments' replica table.
	c.JSON(http.StatusOK, gin.H{
		"message": "Payment history will appear here",
		"data": []interface{}{},
	})
}

func (h *CustomerHandler) GetCommunications(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer uuid"})
		return
	}

	customer, err := h.queries.GetCustomerByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
		return
	}

	if !customer.RazorpayCustomerID.Valid {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}

	history, err := h.queries.GetCommunicationHistory(c.Request.Context(), customer.RazorpayCustomerID)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch communication history"})
		return
	}

	c.JSON(http.StatusOK, history)
}
