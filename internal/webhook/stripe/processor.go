package stripe

import (
	"errors"

	"github.com/stripe/stripe-go"
	"github.com/williamkoller/payment-system/internal/payment/application"
)

type StripeProcessor struct {
	paymentRepo application.PaymentRepository
}

func NewStripeProcessor(paymentRepo application.PaymentRepository) *StripeProcessor {
	return &StripeProcessor{paymentRepo}
}

func (p *StripeProcessor) HandleSucceeded(pi *stripe.PaymentIntent) error {
	if pi == nil {
		return errors.New("nil payment")
	}

	payment, err := p.paymentRepo.FindByStripeID(pi.ID)
	if err != nil {
		return err
	}
	payment.Complete()
	return p.paymentRepo.Update(payment)
}

func (p *StripeProcessor) HandleFailed(pi *stripe.PaymentIntent) error {
	if pi == nil {
		return errors.New("nil PaymentIntent")
	}
	payment, err := p.paymentRepo.FindByStripeID(pi.ID)
	if err != nil {
		return err
	}
	payment.Fail()
	return p.paymentRepo.Update(payment)
}
