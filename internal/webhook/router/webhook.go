package router

import (
	"github.com/gin-gonic/gin"
	"github.com/williamkoller/payment-system/config"
	"github.com/williamkoller/payment-system/internal/payment/infra"
	"github.com/williamkoller/payment-system/internal/webhook/stripe"
)

func SetupWebhookRouter(e *gin.Engine) {
	cfg, err := config.LoadConfiguration()
	if err != nil {
		panic("cannot load configuration: " + err.Error())
	}

	paymentRepo := infra.NewInMemoryPaymentRepository()
	processor := stripe.NewStripeProcessor(paymentRepo)
	handler := stripe.NewStripeWebhookHandler(cfg.Stripe.StripeWebhook, processor)

	e.POST("/webhook/stripe", handler.Handle)
}
