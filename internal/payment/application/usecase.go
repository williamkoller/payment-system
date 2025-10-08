package application

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/williamkoller/payment-system/internal/payment/domain"
	"github.com/williamkoller/payment-system/internal/payment/infra"
	"github.com/williamkoller/payment-system/pkg/ulid"
)

type PaymentRepository interface {
	Save(payment *domain.Payment) error
	FindByID(id string) (*domain.Payment, error)
	FindAll() ([]*domain.Payment, error)
	Remove(id string) error
	Update(payment *domain.Payment) error
}

type PaymentUseCase struct {
	Repository   PaymentRepository
	StripeClient infra.StripeClient
}

type PaymentInput struct {
	Amount   int64
	Currency string
	Email    string
}

func NewPaymentUseCase(Repository PaymentRepository, StripeClient infra.StripeClient) *PaymentUseCase {
	return &PaymentUseCase{Repository: Repository, StripeClient: StripeClient}
}

func (u *PaymentUseCase) CreatePayment(input PaymentInput) (*domain.Payment, error) {
	id := ulid.NewULID()
	payment, err := domain.NewPayment(id, input.Amount, strings.ToUpper(input.Currency), input.Email)

	if err != nil {
		return nil, err
	}

	if err := u.Repository.Save(payment); err != nil {
		return nil, err
	}

	ctx := context.Background()
	intent, err := u.StripeClient.CreatePaymentIntent(ctx, input.Amount, input.Currency, input.Email)
	if err != nil {
		payment.Fail()
		_ = u.Repository.Update(payment)
		return payment, fmt.Errorf("stripe payment failed: %w", err)
	}

	payment.Complete()
	payment.SetStripeID(intent.ID)

	if err := u.Repository.Update(payment); err != nil {
		return payment, err
	}

	return payment, nil
}

func (u *PaymentUseCase) FindPaymentByID(id string) (*domain.Payment, error) {
	paymentFound, err := u.Repository.FindByID(id)
	if err != nil {
		return nil, errors.New("payment not found")
	}
	return paymentFound, nil
}

func (u *PaymentUseCase) Capture(ctx context.Context, paymentID string) (*domain.Payment, error) {
	payment, err := u.Repository.FindByID(paymentID)
	if err != nil {
		return nil, fmt.Errorf("payment not found: %w", err)
	}

	if payment.StripeID == "" {
		return nil, errors.New("missing Stripe payment intent ID")
	}

	err = u.StripeClient.Capture(ctx, payment.StripeID)
	if err != nil {
		payment.Fail()
		_ = u.Repository.Update(payment)
		return payment, fmt.Errorf("stripe capture failed: %w", err)
	}

	payment.Capture()
	if err := u.Repository.Update(payment); err != nil {
		return payment, err
	}

	return payment, nil
}
