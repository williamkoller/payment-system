package stripe

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/webhook"
	"github.com/williamkoller/payment-system/pkg/logger"
)

type StripeWebhookHandler struct {
	secret    string
	processor *StripeProcessor
}

func NewStripeWebhookHandler(secret string, processor *StripeProcessor) *StripeWebhookHandler {
	return &StripeWebhookHandler{secret: secret, processor: processor}
}

func (h *StripeWebhookHandler) Handle(c *gin.Context) {
	const MaxBodyBytes = int64(65536)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)

	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.Default().Errorw("cannot read webhook body", "err", err)
		c.Status(http.StatusServiceUnavailable)
		return
	}

	sigHeader := c.GetHeader("Stripe-Signature")
	event, err := webhook.ConstructEvent(payload, sigHeader, h.secret)
	if err != nil {
		logger.Default().Errorw("invalid webhook signature", "err", err)
		c.Status(http.StatusBadRequest)
		return
	}

	logger.Default().Infow("received stripe webhook event", "type", event.Type)

	switch event.Type {
	case "payment_intent.succeeded":
		var pi stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
			logger.Default().Errorw("failed unmarshal payment_intent.succeeded", "err", err)
		} else {
			if err := h.processor.HandleSucceeded(&pi); err != nil {
				logger.Default().Errorw("processor handle succeeded error", "err", err)
			}
		}
	case "payment_intent.payment_failed":
		var pi stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
			logger.Default().Errorw("failed unmarshal payment_intent.payment_failed", "err", err)
		} else {
			if err := h.processor.HandleFailed(&pi); err != nil {
				logger.Default().Errorw("processor handle failed error", "err", err)
			}
		}
	default:
		logger.Default().Infow("unhandled stripe event", "type", event.Type)
	}

	c.Status(http.StatusOK)
}
