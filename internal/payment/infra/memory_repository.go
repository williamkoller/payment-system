package infra

import (
	"errors"
	"sync"

	"github.com/williamkoller/payment-system/internal/payment/domain"
)

type InMemoryPaymentRepository struct {
	data      map[string]*domain.Payment
	stripeMap map[string]*domain.Payment
	mu        sync.RWMutex
}

func NewInMemoryPaymentRepository() *InMemoryPaymentRepository {
	return &InMemoryPaymentRepository{
		data: make(map[string]*domain.Payment),
	}
}

func (r *InMemoryPaymentRepository) Save(p *domain.Payment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[p.ID] = p
	return nil
}

func (r *InMemoryPaymentRepository) FindByID(id string) (*domain.Payment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.data[id]
	if !ok {
		return nil, errors.New("payment not found")
	}
	return p, nil
}

func (r *InMemoryPaymentRepository) FindAll() ([]*domain.Payment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ps := make([]*domain.Payment, 0, len(r.data))
	for _, v := range r.data {
		ps = append(ps, v)
	}
	return ps, nil
}

func (r *InMemoryPaymentRepository) Remove(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.data, id)
	return nil
}

func (r *InMemoryPaymentRepository) Update(p *domain.Payment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[p.ID] = p
	return nil
}

func (r *InMemoryPaymentRepository) FindByStripeID(stripeID string) (*domain.Payment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.stripeMap[stripeID]
	if !ok {
		return nil, errors.New("stripe not found")
	}
	return p, nil
}

func (r *InMemoryPaymentRepository) FindByIdempotencyKey(idempotencyKey string) (*domain.Payment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.data {
		if p.GetIdempotencyKey() == idempotencyKey {
			return p, nil
		}
	}

	return nil, nil
}
