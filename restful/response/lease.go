package response

import (
	"time"

	"github.com/solabsafrica/afrikanest/model"
)

// swagger:model CreateLeaseResponse
type CreateLeaseResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

func NewCreateLeaseResponse(lease model.Lease) CreateLeaseResponse {
	return CreateLeaseResponse{
		ID:        lease.ID.String(),
		CreatedAt: lease.CreatedAt,
	}
}

// swagger:model LeaseResponse
type LeaseResponse struct {
	model.Lease
}

// swagger:model LeaseChargeResponse
type LeaseChargeResponse struct {
	model.LeaseCharge
}

// swagger:model LeaseChargePaymentResponse
type LeaseChargePaymentResponse struct {
	model.LeaseChargePayment
}

// swagger:model LeaseBalanceResponse
type LeaseBalanceResponse struct {
	Balance int64 `json:"balance"`
}

func NewLeaseBalanceResponse(balance int64) LeaseBalanceResponse {
	return LeaseBalanceResponse{
		Balance: balance,
	}
}

func NewLeaseChargeResponse(leaseCharge model.LeaseCharge) LeaseChargeResponse {
	return LeaseChargeResponse{
		leaseCharge,
	}
}

func NewLeaseChargePaymentResponse(leaseChargePayment model.LeaseChargePayment) LeaseChargePaymentResponse {
	return LeaseChargePaymentResponse{
		leaseChargePayment,
	}
}
