package repository

import (
	"github.com/williamkoller/payment-system/internal/payment/domain"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	Create(p domain.Payment) (*domain.Payment, error)
	FindByID(id string) (*domain.Payment, error)
	FindAll() ([]*domain.Payment, error)
	Remove(id string) error
	Update(p *domain.Payment) error
}

type PaymentRepositoryImpl struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepositoryImpl {
	return &PaymentRepositoryImpl{db: db}
}

func (r *PaymentRepositoryImpl) Create(p domain.Payment) (*domain.Payment, error) {
	if err := r.db.Create(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
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
