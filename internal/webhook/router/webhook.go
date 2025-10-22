package router

import (
	"github.com/gin-gonic/gin"
	"github.com/williamkoller/payment-system/config"
	"github.com/williamkoller/payment-system/internal/payment/repository"
	"github.com/williamkoller/payment-system/internal/webhook/stripe"
	"gorm.io/gorm"
)

func SetupWebhookRouter(e *gin.Engine, db *gorm.DB) {
	cfg, err := config.LoadConfiguration()
	if err != nil {
		panic("cannot load configuration: " + err.Error())
	}

	repo := repository.NewPaymentRepository(db)
	processor := stripe.NewStripeProcessor(repo)
	handler := stripe.NewStripeWebhookHandler(cfg.Stripe.StripeWebhook, processor)

	e.POST("/webhook/stripe", handler.Handle)
}
