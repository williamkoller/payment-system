package infra

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sony/gobreaker"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"
	"github.com/stripe/stripe-go/refund"
	"github.com/williamkoller/payment-system/config"
	"github.com/williamkoller/payment-system/pkg/logger"
)

type StripeClient interface {
	CreatePaymentIntent(ctx context.Context, amount int64, currency, email string, paymentMethod string) (*stripe.PaymentIntent, error)
	Capture(ctx context.Context, piID string) error
	Cancel(ctx context.Context, piID string) error
	Refund(ctx context.Context, stripeID string, amount int64) error
}

var configuration, _ = config.LoadConfiguration()

type stripeClient struct {
	cb *gobreaker.CircuitBreaker
}

func newDefaultCircuitBreaker() *gobreaker.CircuitBreaker {
	settings := gobreaker.Settings{
		Name:        "Stripe",
		MaxRequests: 2,
		Interval:    60 * time.Second,
		Timeout:     10 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 3
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			logger.Default().Infow("circuit breaker state changed", "name", name, "from", from.String(), "to", to.String())
		},
	}
	return gobreaker.NewCircuitBreaker(settings)
}

func NewStripeClient() StripeClient {
	stripe.Key = configuration.StripeApiKey
	return &stripeClient{
		cb: newDefaultCircuitBreaker(),
	}
}

func (c *stripeClient) CreatePaymentIntent(ctx context.Context, amount int64, currency, email string, paymentMethod string) (*stripe.PaymentIntent, error) {
	result, err := c.cb.Execute(func() (interface{}, error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		var lastErr error
		backoff := 100 * time.Millisecond
		maxRetries := 3

		for i := 0; i < maxRetries; i++ {
			params := &stripe.PaymentIntentParams{
				Amount:             stripe.Int64(amount),
				Currency:           stripe.String(currency),
				ReceiptEmail:       stripe.String(email),
				CaptureMethod:      stripe.String(string(stripe.PaymentIntentCaptureMethodManual)),
				Confirm:            stripe.Bool(true),
				PaymentMethod:      stripe.String(configuration.StripeMethod),
				PaymentMethodTypes: []*string{stripe.String(paymentMethod)},
			}

			pi, err := paymentintent.New(params)
			if err == nil {
				return pi, nil
			}

			var stripeErr *stripe.Error
			if errors.As(err, &stripeErr) {
				if stripeErr.HTTPStatusCode == 401 {
					return nil, fmt.Errorf("stripe unauthorized: %s", stripeErr.Msg)
				}
				if stripeErr.HTTPStatusCode >= 400 && stripeErr.HTTPStatusCode < 500 {
					return nil, fmt.Errorf("stripe request error: %s", stripeErr.Msg)
				}
			}

			lastErr = err
			time.Sleep(backoff)
			backoff *= 2
		}

		return nil, lastErr
	})

	if err != nil {
		return nil, err
	}

	pi, ok := result.(*stripe.PaymentIntent)
	if !ok {
		return nil, errors.New("unexpected type from circuit breaker result")
	}

	return pi, nil
}

func (c *stripeClient) Capture(ctx context.Context, piID string) error {
	result, err := c.cb.Execute(func() (interface{}, error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		return paymentintent.Capture(piID, &stripe.PaymentIntentCaptureParams{})
	})

	if err != nil {
		return err
	}

	if _, ok := result.(*stripe.PaymentIntent); !ok {
		return errors.New("unexpected result type from Stripe capture")
	}

	return nil
}

func (c *stripeClient) Cancel(ctx context.Context, piID string) error {
	result, err := c.cb.Execute(func() (interface{}, error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		return paymentintent.Cancel(piID, &stripe.PaymentIntentCancelParams{})
	})

	if err != nil {
		return err
	}

	if _, ok := result.(*stripe.PaymentIntent); !ok {
		return errors.New("unexpected result type from Stripe cancel")
	}

	return nil
}

func (c *stripeClient) Refund(ctx context.Context, stripeID string, amount int64) error {
	logger.Info("stripe payment intent ID", "StripeID", stripeID)

	_, err := c.cb.Execute(func() (interface{}, error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		params := &stripe.RefundParams{
			Amount:        stripe.Int64(amount),
			PaymentIntent: stripe.String(stripeID),
		}

		return refund.New(params)
	})

	return err
}
