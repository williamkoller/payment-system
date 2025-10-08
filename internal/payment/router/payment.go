package router

import (
	"github.com/gin-gonic/gin"
	"github.com/williamkoller/payment-system/internal/payment/application"
	"github.com/williamkoller/payment-system/internal/payment/infra"
	"github.com/williamkoller/payment-system/internal/payment/interfaces"
)

func SetupRouter(e *gin.Engine) *gin.Engine {
	repo := infra.NewInMemoryPaymentRepository()
	stripeClient := infra.NewStripeClient()
	usecase := application.NewPaymentUseCase(repo, stripeClient)
	handler := interfaces.NewPaymentHandler(usecase)
	e.POST("/payments", handler.CreatePayment)
	e.GET("/payments/:id", handler.GetPaymentByID)
	return e
}
