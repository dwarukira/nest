package response

import (
	"time"

	"github.com/solabsafrica/afrikanest/model"
)

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
