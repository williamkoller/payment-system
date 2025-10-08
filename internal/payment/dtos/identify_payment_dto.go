package dtos

type IdentifyPaymentDto struct {
	PaymentID string `uri:"payment_id" binding:"required"`
}
