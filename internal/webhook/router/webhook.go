package router

import (
	"github.com/gin-gonic/gin"
	"github.com/williamkoller/payment-system/config"
	"github.com/williamkoller/payment-system/internal/payment/infra"
	"github.com/williamkoller/payment-system/internal/webhook/stripe"
)

func SetupWebhookRouter(r *gin.Engine) {
	cfg, err := config.LoadConfiguration()
	if err != nil {
		panic("cannot load configuration: " + err.Error())
	}

	paymentRepo := infra.NewInMemoryPaymentRepository()
	processor := stripe.NewStripeProcessor(paymentRepo)
	handler := stripe.NewStripeWebhookHandler(cfg.StripeWebhook, processor)

	r.POST("/webhook/stripe", handler.Handle)
}
