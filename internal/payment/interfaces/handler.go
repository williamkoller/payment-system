package interfaces

import (
	"net/http"
	"strings"

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
		status := ""
		if payment != nil {
			status = string(payment.Status)
		}

		var httpCode int
		var message string

		switch {
		case strings.Contains(err.Error(), "already processed"):
			httpCode = http.StatusOK
			message = "Payment already processed — returning existing transaction"
		case strings.Contains(err.Error(), "already exists"):
			httpCode = http.StatusConflict
			message = "Payment with this idempotency key already exists"
		case strings.Contains(err.Error(), "stripe payment failed"):
			httpCode = http.StatusBadGateway
			message = "Payment creation failed due to Stripe error"
		default:
			httpCode = http.StatusInternalServerError
			message = "Unexpected error while creating payment"
		}

		log.Errorw("Error creating payment", "err", err.Error(), "status", status)
		c.JSON(httpCode, gin.H{
			"error":   err.Error(),
			"status":  status,
			"message": message,
		})
		return
	}

	log.Infow("Created payment", "id", payment.ID, "status", payment.Status)
	c.JSON(http.StatusCreated, ToPaymentResponse(payment))
}

func (h *PaymentHandler) GetPaymentByID(c *gin.Context) {
	var uri dtos.IdentifyPaymentDto

	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment ID"})
		return
	}

	paymentFound, err := h.Usecase.FindPaymentByID(uri)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if paymentFound == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
		return
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
