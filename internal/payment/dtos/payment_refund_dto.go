package dtos

type PaymentRefundDto struct {
	Amount int64 `json:"amount" binding:"required"`
}
