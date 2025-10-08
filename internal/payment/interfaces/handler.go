package interfaces

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/williamkoller/payment-system/internal/middleware"
	"github.com/williamkoller/payment-system/internal/payment/application"
	"github.com/williamkoller/payment-system/internal/payment/dtos"
)

type PaymentHandler struct {
	Usecase *application.PaymentUseCase
}

func NewPaymentHandler(usecase *application.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{Usecase: usecase}
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	log := middleware.FromContext(c)
	var dto dtos.PaymentDto

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := h.Usecase.CreatePayment(application.PaymentInput{
		Amount:   dto.Amount,
		Currency: dto.Currency,
		Email:    dto.Email,
	})

	if err != nil {
		log.Errorw("Error creating payment", "err", err.Error(), "status", payment.Status)
		c.JSON(http.StatusBadGateway, gin.H{
			"error":   err.Error(),
			"status":  payment.Status,
			"message": "Payment creation failed due to external error (Stripe)",
		})
		return
	}

	log.Infow("Created Payment", "id", payment.ID)
	c.JSON(http.StatusCreated, ToPaymentResponse(payment))
}

func (h *PaymentHandler) GetPaymentByID(c *gin.Context) {
	id := c.Param("id")

	paymentFound, err := h.Usecase.FindPaymentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}

	log := middleware.FromContext(c)
	log.Infow("Found Payment", "id", paymentFound.ID)

	c.JSON(http.StatusOK, ToPaymentResponse(paymentFound))

}
