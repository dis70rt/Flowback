package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	
	"github.com/dis70rt/flowback/internal/repo"
	"github.com/dis70rt/flowback/internal/razorpay"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

type CaseHandler struct {
	queries        *repo.Queries
	razorpayClient *razorpay.Client
}

func NewCaseHandler(q *repo.Queries, rzpClient *razorpay.Client) *CaseHandler {
	return &CaseHandler{
		queries:        q,
		razorpayClient: rzpClient,
	}
}

type CaseItemDTO struct {
	ID                  string         `json:"id"`
	SubscriptionID      string         `json:"subscription_id"`
	AmountAtRisk        int64          `json:"amount_at_risk"`
	Status              string         `json:"status"`
	CreatedAt           time.Time      `json:"created_at"`
	LatestActionType    sql.NullString `json:"latest_action_type"`
	LatestActionStatus  sql.NullString `json:"latest_action_status"`
	LatestActionChannel sql.NullString `json:"latest_action_channel"`
}

type ListCasesResponse struct {
	Page  int           `json:"page"`
	Limit int           `json:"limit"`
	Data  []CaseItemDTO `json:"data"`
}

type EditDraftRequest struct {
	DraftPayload json.RawMessage `json:"draft_payload" binding:"required"`
}

func (h *CaseHandler) ListCases(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	filter := c.DefaultQuery("filter", "all")
	
	if page < 1 { page = 1 }
	if limit < 1 || limit > 100 { limit = 20 }
	
	offset := (page - 1) * limit

	var dtos []CaseItemDTO
	
	if filter == "pending" {
		dbCases, err := h.queries.ListPendingCases(c.Request.Context(), repo.ListPendingCasesParams{
			Limit:  int32(limit),
			Offset: int32(offset),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pending cases"})
			return
		}
		for _, row := range dbCases {
			dtos = append(dtos, CaseItemDTO{
				ID:                  row.ID.String(),
				SubscriptionID:      row.SubscriptionID,
				AmountAtRisk:        row.AmountAtRisk,
				Status:              string(row.Status),
				CreatedAt:           row.CreatedAt,
				LatestActionType:    sql.NullString{String: string(row.LatestActionType), Valid: string(row.LatestActionType) != ""},
				LatestActionStatus:  sql.NullString{String: string(row.LatestActionStatus), Valid: string(row.LatestActionStatus) != ""},
				LatestActionChannel: row.LatestActionChannel,
			})
		}
	} else {
		dbCases, err := h.queries.ListRecoveryCases(c.Request.Context(), repo.ListRecoveryCasesParams{
			Limit:  int32(limit),
			Offset: int32(offset),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cases"})
			return
		}
		for _, row := range dbCases {
			dtos = append(dtos, CaseItemDTO{
				ID:                  row.ID.String(),
				SubscriptionID:      row.SubscriptionID,
				AmountAtRisk:        row.AmountAtRisk,
				Status:              string(row.Status),
				CreatedAt:           row.CreatedAt,
				LatestActionType:    sql.NullString{String: string(row.LatestActionType), Valid: string(row.LatestActionType) != ""},
				LatestActionStatus:  sql.NullString{String: string(row.LatestActionStatus), Valid: string(row.LatestActionStatus) != ""},
				LatestActionChannel: row.LatestActionChannel,
			})
		}
	}
	
	c.JSON(http.StatusOK, ListCasesResponse{
		Page:  page,
		Limit: limit,
		Data:  dtos,
	})
}

func (h *CaseHandler) GetCase(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid case uuid"})
		return
	}
	
	caseData, err := h.queries.GetRecoveryCaseByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "case not found"})
		return
	}

	actions, err := h.queries.GetActionsByCase(c.Request.Context(), id)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch case actions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"case":    caseData,
		"actions": actions,
	})
}

func (h *CaseHandler) ApproveDraft(c *gin.Context) {
	actionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid action uuid"})
		return
	}

	data, err := h.queries.GetActionAndCaseForApproval(c.Request.Context(), actionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "action not found"})
		return
	}

	amount := data.AmountAtRisk
	if data.DiscountPercentage.Valid && data.DiscountPercentage.Int32 > 0 {
		amount = amount * int64(100 - data.DiscountPercentage.Int32) / 100
	}

	linkID, linkURL, err := h.razorpayClient.CreatePaymentLink(amount, data.CaseID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate payment link: " + err.Error()})
		return
	}

	mockClerkID := sql.NullString{String: "user_clerk_admin", Valid: true}

	err = h.queries.ApproveAction(c.Request.Context(), repo.ApproveActionParams{
		ID:                actionID,
		ApprovedByClerkID: mockClerkID,
		PaymentLinkID:     sql.NullString{String: linkID, Valid: true},
		PaymentLinkUrl:    sql.NullString{String: linkURL, Valid: true},
	})
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve action"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"status": "approved", "payment_link_url": linkURL})
}

func (h *CaseHandler) EditDraft(c *gin.Context) {
	actionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid action uuid"})
		return
	}

	var req EditDraftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.queries.UpdateActionDraft(c.Request.Context(), repo.UpdateActionDraftParams{
		ID:           actionID,
		DraftPayload: pqtype.NullRawMessage{RawMessage: req.DraftPayload, Valid: len(req.DraftPayload) > 0},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update draft"})
		return
	}

	// Edit acts as Approve
	h.ApproveDraft(c)
}

func (h *CaseHandler) RejectDraft(c *gin.Context) {
	actionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid action uuid"})
		return
	}

	mockClerkID := sql.NullString{String: "user_clerk_admin", Valid: true}

	err = h.queries.RejectAction(c.Request.Context(), repo.RejectActionParams{
		ID:                actionID,
		ApprovedByClerkID: mockClerkID,
	})
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject action"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"status": "rejected"})
}

func (h *CaseHandler) GetCaseSummary(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return
	}

	row, err := h.queries.GetCaseSummary(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "case not found"})
		return
	}

	c.JSON(http.StatusOK, CaseItemDTO{
		ID:                  row.ID.String(),
		SubscriptionID:      row.SubscriptionID,
		AmountAtRisk:        row.AmountAtRisk,
		Status:              string(row.Status),
		CreatedAt:           row.CreatedAt,
		LatestActionType:    sql.NullString{String: string(row.LatestActionType), Valid: string(row.LatestActionType) != ""},
		LatestActionStatus:  sql.NullString{String: string(row.LatestActionStatus), Valid: string(row.LatestActionStatus) != ""},
		LatestActionChannel: row.LatestActionChannel,
	})
}
