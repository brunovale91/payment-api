package services

import (
	"github.com/brunovale91/payment-api/store"
	"github.com/brunovale91/payment-api/types"
	"github.com/google/uuid"
)

type PaymentService interface {

	// Generate id, creates payment and returns created payment
	CreatePayment(*types.Payment) (*types.Payment, error)

	// Update payment attributes and return updated payment
	UpdatePayment(string, *types.PaymentAttributes) (*types.Payment, error)

	// Delete payment
	DeletePayment(string) (bool, error)

	// Get payment
	GetPayment(string) (*types.Payment, error)

	// Get slice of payments
	GetPayments() ([]*types.Payment, error)
}

type PaymentServiceImpl struct {
	store store.PaymentStore
}

func NewPaymentService(paymentStore store.PaymentStore) PaymentService {
	return PaymentServiceImpl{
		store: paymentStore,
	}
}

func (p PaymentServiceImpl) CreatePayment(payment *types.Payment) (*types.Payment, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	payment.Id = id.String()
	return p.store.CreatePayment(payment)
}

func (p PaymentServiceImpl) UpdatePayment(id string, attributes *types.PaymentAttributes) (*types.Payment, error) {
	return p.store.UpdatePayment(id, attributes)
}

func (p PaymentServiceImpl) DeletePayment(id string) (bool, error) {
	return p.store.DeletePayment(id)
}

func (p PaymentServiceImpl) GetPayment(id string) (*types.Payment, error) {
	return p.store.GetPayment(id)
}

func (p PaymentServiceImpl) GetPayments() ([]*types.Payment, error) {
	return p.store.GetPayments()
}
