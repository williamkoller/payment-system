package interfaces

import (
	"time"

	"github.com/williamkoller/payment-system/internal/payment/domain"
)

type PaymentResponse struct {
	ID             string               `json:"id"`
	Amount         int64                `json:"amount"`
	Currency       string               `json:"currency"`
	Status         domain.PaymentStatus `json:"status"`
	Email          string               `json:"email"`
	StripeID       string               `json:"stripe_id"`
	PaymentMethod  string               `json:"payment_method"`
	IdempotencyKey string               `json:"idempotency_key"`
	CreatedAt      time.Time            `json:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at"`
}

func ToPaymentResponse(p *domain.Payment) PaymentResponse {
	return PaymentResponse{
		ID:             p.ID,
		Amount:         p.Amount,
		Currency:       p.Currency,
		Status:         p.Status,
		Email:          p.Email,
		StripeID:       p.StripeID,
		PaymentMethod:  p.PaymentMethod,
		IdempotencyKey: p.IdempotencyKey,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
}
