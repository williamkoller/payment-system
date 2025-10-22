package router

import (
	"github.com/gin-gonic/gin"
	"github.com/williamkoller/payment-system/internal/payment/application"
	"github.com/williamkoller/payment-system/internal/payment/infra"
	"github.com/williamkoller/payment-system/internal/payment/interfaces"
	"github.com/williamkoller/payment-system/internal/payment/repository"
	"gorm.io/gorm"
)

func SetupRouter(e *gin.Engine, db *gorm.DB) {
	repo := repository.NewPaymentRepository(db)
	stripeClient := infra.NewStripeClient()
	usecase := application.NewPaymentUseCase(repo, stripeClient)
	handler := interfaces.NewPaymentHandler(usecase)
	payments := e.Group("/payments")
	{
		payments.POST("/", handler.CreatePayment)
		payments.GET("/:payment_id", handler.GetPaymentByID)
		payments.POST("/:payment_id/capture", handler.CapturePayment)
		payments.POST("/:payment_id/cancel", handler.CancelPayment)
		payments.POST("/:payment_id/refund", handler.RefundPayment)
	}
}
