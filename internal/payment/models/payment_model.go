package payment_model

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	Amount         int64
	Currency       string
	Status         string
	Email          string
	PaymentMethod  string
	IdempotencyKey string
	StripeID       string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
