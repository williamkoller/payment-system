package interfaces

import (
	"fmt"
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
	var dto dtos.AddPaymentDto

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := h.Usecase.CreatePayment(application.PaymentInput{
		Amount:        dto.Amount,
		Currency:      dto.Currency,
		Email:         dto.Email,
		PaymentMethod: dto.PaymentMethod,
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
	var uri dtos.IdentifyPaymentDto
	fmt.Println(uri)
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment ID"})
	}

	paymentFound, err := h.Usecase.FindPaymentByID(uri)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}

	log := middleware.FromContext(c)
	log.Infow("Found Payment", "id", paymentFound.ID)

	c.JSON(http.StatusOK, ToPaymentResponse(paymentFound))

}

func (h *PaymentHandler) CapturePayment(c *gin.Context) {
	var uri dtos.IdentifyPaymentDto
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment ID"})
	}
	log := middleware.FromContext(c)
	ctx := c.Request.Context()
	payment, err := h.Usecase.Capture(ctx, uri)
	if err != nil {
		log.Errorw("Capture failed", "err", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ToPaymentResponse(payment))
}

func (h *PaymentHandler) CancelPayment(c *gin.Context) {
	var uri dtos.IdentifyPaymentDto
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment ID"})
	}
	log := middleware.FromContext(c)
	ctx := c.Request.Context()
	payment, err := h.Usecase.Cancel(ctx, uri)
	if err != nil {
		log.Errorw("Cancel failed", "err", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ToPaymentResponse(payment))
}

func (h *PaymentHandler) RefundPayment(c *gin.Context) {
	log := middleware.FromContext(c)
	ctx := c.Request.Context()

	var uri dtos.IdentifyPaymentDto
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment id"})
		return
	}

	var pr dtos.PaymentRefundDto
	if err := c.ShouldBindJSON(&pr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid amount"})
		return
	}

	log.Infow("Refund Payment", "payment_id", uri.PaymentID, "amount", pr.Amount)

	payment, err := h.Usecase.Refund(ctx, uri, pr)
	if err != nil {
		log.Errorw("Refund failed", "err", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ToPaymentResponse(payment))
}
