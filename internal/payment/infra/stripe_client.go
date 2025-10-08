package infra

import (
	"errors"
	"fmt"
	"time"

	"github.com/sony/gobreaker"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"
	"github.com/williamkoller/payment-system/config"
	"github.com/williamkoller/payment-system/pkg/logger"
)

type StripeClient interface {
	CreatePaymentIntent(amount int64, currency, email string) (*stripe.PaymentIntent, error)
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
	configuration, _ := config.LoadConfiguration()
	stripe.Key = configuration.StripeApiKey
	return &stripeClient{
		cb: newDefaultCircuitBreaker(),
	}
}

func (s *stripeClient) CreatePaymentIntent(amount int64, currency, email string) (*stripe.PaymentIntent, error) {
	result, err := s.cb.Execute(func() (interface{}, error) {
		var lastErr error
		backoff := 100 * time.Millisecond
		maxRetries := 3

		for i := 0; i < maxRetries; i++ {
			params := &stripe.PaymentIntentParams{
				Amount:        stripe.Int64(amount),
				Currency:      stripe.String(currency),
				ReceiptEmail:  stripe.String(email),
				CaptureMethod: stripe.String(string(stripe.PaymentIntentCaptureMethodManual)),
				Confirm:       stripe.Bool(true),
				PaymentMethod: stripe.String(configuration.StripeMethod),
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
