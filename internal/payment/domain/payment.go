package domain

import (
	"errors"
	"time"
)

type PaymentStatus string

const (
	StatusPending   PaymentStatus = "PENDING"
	StatusCompleted PaymentStatus = "COMPLETED"
	StatusFailed    PaymentStatus = "FAILED"
	StatusCanceled  PaymentStatus = "CANCELED"
	StatusCaptured  PaymentStatus = "CAPTURED"
)

type Payment struct {
	ID        string
	StripeID  string
	Amount    int64
	Currency  string
	Status    PaymentStatus
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewPayment(id string, amount int64, currency, email string) (*Payment, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	if currency == "" {
		return nil, errors.New("currency must be empty")
	}

	if email == "" {
		return nil, errors.New("email must be empty")
	}

	now := time.Now()

	return &Payment{
		ID:        id,
		Amount:    amount,
		Currency:  currency,
		Status:    StatusPending,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (p *Payment) Complete() {
	p.Status = StatusCompleted
	p.UpdatedAt = time.Now()
}

func (p *Payment) Cancel() {
	p.Status = StatusCanceled
	p.UpdatedAt = time.Now()
}

func (p *Payment) Fail() {
	p.Status = StatusFailed
	p.UpdatedAt = time.Now()
}

func (p *Payment) Capture() {
	p.Status = StatusCaptured
	p.UpdatedAt = time.Now()
}

func (p *Payment) GetID() string {
	return p.ID
}

func (p *Payment) GetStripeID() string {
	return p.StripeID
}

func (p *Payment) GetAmount() int64 {
	return p.Amount
}

func (p *Payment) GetCurrency() string {
	return p.Currency
}

func (p *Payment) GetStatus() PaymentStatus {
	return p.Status
}

func (p *Payment) GetCreatedAt() time.Time {
	return p.CreatedAt
}

func (p *Payment) GetUpdatedAt() time.Time {
	return p.UpdatedAt
}

func (p *Payment) GetEmail() string {
	return p.Email
}

func (p *Payment) SetStripeID(stripeID string) {
	p.StripeID = stripeID
}
