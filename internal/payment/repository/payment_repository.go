package repository

import (
	"github.com/williamkoller/payment-system/internal/payment/domain"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	Save(payment domain.Payment) (*domain.Payment, error)
	FindByID(id string) (*domain.Payment, error)
	FindAll() ([]*domain.Payment, error)
	Remove(id string) error
	Update(p *domain.Payment) error
	FindByStripeID(stripeID string) (*domain.Payment, error)
	FindByIdempotencyKey(idempotencyKey string) (*domain.Payment, error)
}

type PaymentRepositoryImpl struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepositoryImpl {
	return &PaymentRepositoryImpl{db: db}
}

func (r *PaymentRepositoryImpl) Save(payment *domain.Payment) (*domain.Payment, error) {
	if err := r.db.Create(&payment).Error; err != nil {
		return nil, err
	}
	return payment, nil
}

func (r *PaymentRepositoryImpl) FindByID(id string) (*domain.Payment, error) {
	var payment domain.Payment
	if err := r.db.First(&payment, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepositoryImpl) FindAll() ([]*domain.Payment, error) {
	var payments []*domain.Payment
	if err := r.db.Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *PaymentRepositoryImpl) Remove(id string) error {
	return r.db.Delete(&domain.Payment{}, "id = ?", id).Error
}

func (r *PaymentRepositoryImpl) Update(p *domain.Payment) error {
	return r.db.Model(&domain.Payment{}).
		Select("StripeID", "Amount", "Currency", "Status", "Email", "PaymentMethod", "IdempotencyKey").
		Where("id = ?", p.ID).
		Updates(p).Error
}

func (r *PaymentRepositoryImpl) FindByStripeID(stripeID string) (*domain.Payment, error) {
	var payment domain.Payment
	if err := r.db.First(&payment, "stripe_id = ?", stripeID).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepositoryImpl) FindByIdempotencyKey(idempotencyKey string) (*domain.Payment, error) {
	var payment domain.Payment
	if err := r.db.First(&payment, "idempotency_key = ?", idempotencyKey).Error; err != nil {
		return nil, err
	}

	return &payment, nil
}
