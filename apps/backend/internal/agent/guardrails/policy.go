package guardrails

import (
	"fmt"
	"time"

	"github.com/dis70rt/flowback/internal/agent/core"
	"github.com/dis70rt/flowback/internal/repo"
)

// PolicyResult holds the outcome of the guardrail evaluation
type PolicyResult struct {
	IsApproved  bool
	BlockReason string
	// AdjustedDelay allows the policy engine to fix a bad time without failing the whole action
	AdjustedDelay int
}

// EvaluatePolicy acts as a strict deterministic guardrail against LLM hallucinations or bad strategies.
func EvaluatePolicy(
	strategy core.StrategyOutput,
	customer repo.Customer,
	recoveryCase repo.RecoveryCase,
) PolicyResult {

	result := PolicyResult{
		IsApproved:    true,
		AdjustedDelay: strategy.DelayHours,
	}

	// -------------------------------------------------------------------------
	// GUARDRAIL 1: Fraud & Risk Prevention
	// If a user has a terrible reliability score, the AI is NOT allowed to
	// send them a new payment link (which they might exploit) or give a discount.
	// -------------------------------------------------------------------------
	if customer.ReliabilityScore < 0.3 {
		if strategy.Action != "silent_retry" && strategy.Action != "escalate_to_human" {
			result.IsApproved = false
			result.BlockReason = fmt.Sprintf("High Risk Customer (Score: %.2f). Action '%s' blocked. Only silent retries permitted.", customer.ReliabilityScore, strategy.Action)
			return result
		}
	}

	// -------------------------------------------------------------------------
	// GUARDRAIL 2: RBI Communications Compliance (Time of Day)
	// Commercial SMS/WhatsApp messages cannot be sent during quiet hours (9 PM to 8 AM).
	// If the AI suggests sending a message that lands in that window, we adjust the delay.
	// -------------------------------------------------------------------------
	if strategy.Action == "send_sms" || strategy.Action == "send_whatsapp" {
		targetTime := time.Now().Add(time.Duration(result.AdjustedDelay) * time.Hour)
		hour := targetTime.Hour()

		// If execution time is between 21:00 (9 PM) and 08:00 (8 AM)
		if hour >= 21 || hour < 8 {
			// Calculate hours until 9 AM the next day
			var hoursUntilMorning int
			if hour >= 21 {
				hoursUntilMorning = (24 - hour) + 9
			} else {
				hoursUntilMorning = 9 - hour
			}
			result.AdjustedDelay += hoursUntilMorning
			// We don't block it, we just safely delay it.
		}
	}

	// -------------------------------------------------------------------------
	// GUARDRAIL 3: Anti-Spam (Max Contacts)
	// If we have already spammed this user beyond the threshold, block the action.
	// -------------------------------------------------------------------------
	if recoveryCase.ContactCount >= recoveryCase.MaxContacts {
		if strategy.Action == "send_email" || strategy.Action == "send_sms" || strategy.Action == "send_whatsapp" {
			result.IsApproved = false
			result.BlockReason = fmt.Sprintf("Contact limit reached (%d/%d). Further communication blocked.", recoveryCase.ContactCount, recoveryCase.MaxContacts)
			return result
		}
	}

	// -------------------------------------------------------------------------
	// GUARDRAIL 4: LLM Hallucination Check
	// Ensure the action is actually supported by our execution engine.
	// -------------------------------------------------------------------------
	validActions := map[string]bool{
		"silent_retry":      true,
		"send_email":        true,
		"send_sms":          true,
		"send_whatsapp":     true,
		"escalate_to_human": true,
	}
	if !validActions[strategy.Action] {
		result.IsApproved = false
		result.BlockReason = fmt.Sprintf("LLM hallucinated invalid action type: %s", strategy.Action)
		return result
	}

	return result
}
