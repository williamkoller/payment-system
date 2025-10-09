package application

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/stripe/stripe-go"
	"github.com/williamkoller/payment-system/internal/payment/domain"
	"github.com/williamkoller/payment-system/internal/payment/dtos"
	"github.com/williamkoller/payment-system/internal/payment/infra"
	"github.com/williamkoller/payment-system/pkg/ulid"
)

type PaymentRepository interface {
	Save(payment *domain.Payment) error
	FindByID(id string) (*domain.Payment, error)
	FindAll() ([]*domain.Payment, error)
	Remove(id string) error
	Update(payment *domain.Payment) error
	FindByStripeID(stripeID string) (*domain.Payment, error)
}

type PaymentUseCase struct {
	Repository   PaymentRepository
	StripeClient infra.StripeClient
}

type PaymentInput struct {
	Amount        int64
	Currency      string
	Email         string
	PaymentMethod string
}

func NewPaymentUseCase(Repository PaymentRepository, StripeClient infra.StripeClient) *PaymentUseCase {
	return &PaymentUseCase{Repository: Repository, StripeClient: StripeClient}
}

func (u *PaymentUseCase) CreatePayment(input PaymentInput) (*domain.Payment, error) {
	id := ulid.NewULID()
	payment, err := domain.NewPayment(id, input.Amount, strings.ToUpper(input.Currency), input.Email, input.PaymentMethod)

	if err != nil {
		return nil, err
	}

	if err := u.Repository.Save(payment); err != nil {
		return nil, err
	}

	ctx := context.Background()
	intent, err := u.StripeClient.CreatePaymentIntent(ctx, input.Amount, input.Currency, input.Email, input.PaymentMethod)
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

func (u *PaymentUseCase) FindPaymentByID(i dtos.IdentifyPaymentDto) (*domain.Payment, error) {
	paymentFound, err := u.Repository.FindByID(i.PaymentID)
	if err != nil {
		return nil, errors.New("payment not found")
	}
	return paymentFound, nil
}

func (u *PaymentUseCase) Capture(ctx context.Context, i dtos.IdentifyPaymentDto) (*domain.Payment, error) {
	payment, err := u.Repository.FindByID(i.PaymentID)
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

func (u *PaymentUseCase) Cancel(ctx context.Context, i dtos.IdentifyPaymentDto) (*domain.Payment, error) {
	payment, err := u.Repository.FindByID(i.PaymentID)
	if err != nil {
		return nil, fmt.Errorf("payment not found: %w", err)
	}

	if payment.StripeID == "" {
		return nil, errors.New("missing Stripe payment intent ID")
	}

	if err := payment.CanCancel(); err != nil {
		payment.Fail()
		_ = u.Repository.Update(payment)
		return payment, fmt.Errorf("stripe cancel failed: %w", err)
	}

	err = u.StripeClient.Cancel(ctx, payment.StripeID)
	if err != nil {
		var stripeErr *stripe.Error
		if errors.As(err, &stripeErr) {
			if stripeErr.Code == stripe.ErrorCodePaymentIntentUnexpectedState {
				payment.Capture()
				_ = u.Repository.Update(payment)
				return payment, fmt.Errorf("cannot cancel payment: already captured on Stripe")
			}

			if stripeErr.HTTPStatusCode >= 400 && stripeErr.HTTPStatusCode < 500 {
				return payment, fmt.Errorf("stripe client error: %s (%s)", stripeErr.Msg, stripeErr.Code)
			}

			if stripeErr.HTTPStatusCode >= 500 {
				return payment, fmt.Errorf("stripe server error: %s", stripeErr.Msg)
			}
		}

		payment.Fail()
		_ = u.Repository.Update(payment)
		return payment, fmt.Errorf("stripe cancel failed: %w", err)
	}

	payment.Cancel()
	if err := u.Repository.Update(payment); err != nil {
		return payment, err
	}
	return payment, nil
}

func (u *PaymentUseCase) Refund(ctx context.Context, uri dtos.IdentifyPaymentDto, pr dtos.PaymentRefundDto) (*domain.Payment, error) {
	payment, err := u.Repository.FindByID(uri.PaymentID)
	if err != nil {
		return nil, fmt.Errorf("payment not found: %w", err)
	}

	if payment.StripeID == "" {
		return nil, errors.New("missing Stripe payment intent ID")
	}

	if err := payment.CanRefund(); err != nil {
		payment.Fail()
		_ = u.Repository.Update(payment)
		return payment, fmt.Errorf("stripe refund failed: %w", err)
	}

	err = u.StripeClient.Refund(ctx, payment.StripeID, pr.Amount)
	if err != nil {
		payment.Fail()
		_ = u.Repository.Update(payment)
		return payment, fmt.Errorf("stripe refund failed: %w", err)
	}

	payment.Refund()
	if err := u.Repository.Update(payment); err != nil {
		return payment, err
	}

	return payment, nil
}
