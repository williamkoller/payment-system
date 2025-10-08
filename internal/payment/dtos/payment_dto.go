package dtos

type PaymentDto struct {
	Amount   int64  `json:"amount" binding:"required,gt=0"`
	Currency string `json:"currency" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}
